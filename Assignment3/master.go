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
	time time.Duration
	addr *net.UDPAddr
}

var startTime time.Time
var slaves = make(map[int]data)

func broadcastTime(conn *net.UDPConn, syncTime time.Duration, wg *sync.WaitGroup) {
	for _, slave := range slaves {
		_, err := conn.WriteToUDP([]byte(syncTime.String()), slave.addr)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	sync := time.Since(startTime) - syncTime
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
		masterTime := time.Since(startTime)
		sumDiff := masterTime
		for _, slave := range slaves {
			sumDiff += masterTime - slave.time
		}
		syncTime := masterTime + sumDiff/time.Duration(len(slaves)+1)
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
		t, _ := time.ParseDuration(string(p))
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
