package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

type data struct {
	time int64
	addr *net.UDPAddr
}

var startTime time.Time
var slaves = make(map[int]data)

func getTime() int64 {
	return int64(time.Since(startTime))
}

func getTimeString() string {
	return strconv.FormatInt(getTime(), 10)
}

func broadcastTime(conn *net.UDPConn, syncTime int64, wg *sync.WaitGroup) {
	for _, slave := range slaves {
		sendData := strconv.FormatInt(syncTime, 10)
		_, err := conn.WriteToUDP([]byte(sendData), slave.addr)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	syncDiff := strconv.FormatInt(getTime()-syncTime, 10) + "ns"
	sync, err := time.ParseDuration(syncDiff)
	if err != nil {
		fmt.Println(err.Error())
	}
	startTime = startTime.Add(sync)
	defer wg.Done()
}

func synchronizeClocks(conn *net.UDPConn, wg *sync.WaitGroup) {
	for {
		time.Sleep(5 * time.Second)
		// fmt.Println("Synchronizing Time")
		if len(slaves) == 0 {
			continue
		}
		fmt.Println(len(slaves))
		masterTime := getTime()
		var sumDiff int64 = 0
		for _, slave := range slaves {
			sumDiff += masterTime - slave.time
		}
		syncTime := masterTime + (sumDiff / int64(len(slaves)+1))
		wg.Add(1)
		go broadcastTime(conn, syncTime, wg)
	}
	defer wg.Done()
}

func receiveTime(conn *net.UDPConn, wg *sync.WaitGroup) {
	p := make([]byte, 4096)
	for {
		_, remoteaddr, err := conn.ReadFromUDP(p)
		p = bytes.Trim(p, "\x00")
		// fmt.Println("Received time from", remoteaddr)
		if err != nil {
			fmt.Println(err.Error())
		}
		t, _ := strconv.ParseInt(string(p), 10, 64)
		slaves[remoteaddr.Port] = data{t, remoteaddr}
	}
	defer wg.Done()
}

func printTime(wg *sync.WaitGroup) {
	for {
		// fmt.Println(getTime())
		fmt.Println("Local time is", time.Since(startTime))
		time.Sleep(5 * time.Second)
	}
	defer wg.Done()
}

func main() {
	startTime = time.Now()
	portS := os.Args[1]
	port, err := strconv.Atoi(portS)
	addr := net.UDPAddr{Port: port, IP: net.ParseIP("127.0.0.1")}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer conn.Close()

	var wg sync.WaitGroup
	wg.Add(3)
	go receiveTime(conn, &wg)
	go synchronizeClocks(conn, &wg)
	go printTime(&wg)
	wg.Wait()
}
