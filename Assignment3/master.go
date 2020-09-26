package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

var startTime time.Time
var slaves = make(map[net.UDPAddr]int64)

func broadcastTime(conn *net.UDPConn, syncTime int64, wg *sync.WaitGroup) {
	fmt.Println("Length is ", slaves)
	for slaveAddr, slaveTime := range slaves {
		sTime := syncTime - slaveTime
		fmt.Println("Sent time to", slaveAddr)
		_, err := conn.WriteToUDP([]byte(strconv.FormatInt(sTime, 10)), slaveAddr)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	defer wg.Done()
}

func synchronizeClocks(conn *net.UDPConn, wg *sync.WaitGroup) {
	for {
		time.Sleep(5 * time.Second)
		fmt.Println("Synchronizing Time")
		if len(slaves) == 0 {
			continue
		}
		masterTime := int64(time.Now().Sub(startTime))
		var sumDiff int64 = 0
		for _, slaveTime := range slaves {
			sumDiff += masterTime - slaveTime
		}
		syncTime := masterTime + sumDiff/int64(len(slaves))
		wg.Add(1)
		go broadcastTime(conn, syncTime, wg)
	}
	defer wg.Done()
}

func receiveTime(conn *net.UDPConn, wg *sync.WaitGroup) {
	p := make([]byte, 2048)
	for {
		_, remoteaddr, err := conn.ReadFromUDP(p)
		fmt.Println("Received time from", remoteaddr)
		if err != nil {
			fmt.Println(err.Error())
		}
		t, _ := strconv.Atoi(string(p))
		slaves[remoteaddr] = int64(t)
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
	wg.Add(2)
	go receiveTime(conn, &wg)
	go synchronizeClocks(conn, &wg)
	wg.Wait()
}
