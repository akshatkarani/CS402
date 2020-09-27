package main

import (
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

func broadcastTime(conn *net.UDPConn, syncTime int64, wg *sync.WaitGroup) {
	for _, slave := range slaves {
		sTime := syncTime - slave.time
		sendData := []byte(strconv.FormatInt(sTime, 10))
		_, err := conn.WriteToUDP(sendData, slave.addr)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	defer wg.Done()
}

func synchronizeClocks(conn *net.UDPConn, wg *sync.WaitGroup) {
	for {
		time.Sleep(5 * time.Second)
		// fmt.Println("Synchronizing Time")
		if len(slaves) == 0 {
			continue
		}
		masterTime := int64(time.Now().Sub(startTime))
		var sumDiff int64 = 0
		for _, slave := range slaves {
			sumDiff += masterTime - slave.time
		}
		syncTime := masterTime + sumDiff/int64(len(slaves))
		wg.Add(1)
		go broadcastTime(conn, syncTime, wg)
	}
	defer wg.Done()
}

func receiveTime(conn *net.UDPConn, wg *sync.WaitGroup) {
	p := make([]byte, 4096)
	for {
		_, remoteaddr, err := conn.ReadFromUDP(p)
		// fmt.Println("Received time from", remoteaddr)
		if err != nil {
			fmt.Println(err.Error())
		}
		t, _ := strconv.Atoi(string(p))
		slaves[remoteaddr.Port] = data{int64(t), remoteaddr}
	}
	defer wg.Done()
}

func printTime(wg *sync.WaitGroup) {
	for {
		fmt.Println("Local time is", time.Now().Sub(startTime))
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
