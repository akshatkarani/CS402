package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

var startTime time.Time

func startSendingTime(conn net.Conn, wg *sync.WaitGroup) {
	for {
		fmt.Fprintf(conn, time.Since(startTime).String())
		// fmt.Println("Recent time successfully sent")
		time.Sleep(5 * time.Second)
	}
	wg.Done()
}

func startReceivingTime(conn net.Conn, wg *sync.WaitGroup) {
	for {
		p := make([]byte, 4096)
		_, err := bufio.NewReader(conn).Read(p)
		if err != nil {
			fmt.Println("Error", err.Error())
		}
		p = bytes.Trim(p, "\x00")
		syncTime, _ := time.ParseDuration(string(p))
		sync := time.Since(startTime) - syncTime
		startTime = startTime.Add(sync)
		// fmt.Println("Time Synchronized")
		time.Sleep(1 * time.Second)
	}
	defer wg.Done()
}

func printTime(wg *sync.WaitGroup) {
	for {
		fmt.Println("Local time is", time.Since(startTime))
		time.Sleep(5 * time.Second)
	}
	defer wg.Done()
}

func main() {
	startTime = time.Now()
	ip, port := os.Args[1], os.Args[2]
	conn, err := net.Dial("udp", ip+":"+port)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer conn.Close()

	var wg sync.WaitGroup
	wg.Add(3)
	go startSendingTime(conn, &wg)
	go startReceivingTime(conn, &wg)
	go printTime(&wg)
	wg.Wait()
}
