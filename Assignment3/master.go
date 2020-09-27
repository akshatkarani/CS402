package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type data struct {
	time float64
	addr *net.UDPAddr
}

var startTime time.Time
var slaves = make(map[int]data)

func getTimeString(t time.Time) string {
	str := t.String()
	i := strings.Index(str, " m=")
	return str[i+4:]
}

func getTime(t time.Time) float64 {
	ret, err := strconv.ParseFloat(getTimeString(t), 64)
	if err != nil {
		fmt.Println(err.Error())
	}
	return ret
}

func broadcastTime(conn *net.UDPConn, syncTime float64, wg *sync.WaitGroup) {
	for _, slave := range slaves {
		sendData := strconv.FormatFloat(syncTime, 'f', -1, 64)
		_, err := conn.WriteToUDP([]byte(sendData), slave.addr)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	syncFl := syncTime - getTime(startTime)
	sync, err := time.ParseDuration(strconv.FormatFloat(syncFl, 'f', -1, 64) + "s")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(startTime, sync)
	startTime = startTime.Add(sync)
	fmt.Println(startTime)
	defer wg.Done()
}

func synchronizeClocks(conn *net.UDPConn, wg *sync.WaitGroup) {
	for {
		time.Sleep(5 * time.Second)
		// fmt.Println("Synchronizing Time")
		if len(slaves) == 0 {
			continue
		}
		masterTime := getTime(startTime)
		var sumDiff float64 = 0
		for _, slave := range slaves {
			sumDiff += masterTime - slave.time
		}
		syncTime := masterTime + sumDiff/float64(len(slaves))
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
		t, _ := strconv.ParseFloat(string(p), 64)
		slaves[remoteaddr.Port] = data{t, remoteaddr}
	}
	defer wg.Done()
}

func printTime(wg *sync.WaitGroup) {
	for {
		fmt.Println("Local time is", getTime(time.Now())-getTime(startTime))
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
