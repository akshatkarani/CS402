package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

var startTime time.Time

func startSendingTime(conn net.Conn, wg *sync.WaitGroup) {
	for {
		sendTime := time.Now().Sub(startTime)
		startTime.Add(sendTime)
		fmt.Fprintf(conn, string(sendTime))
		fmt.Println("Recent time successfully sent")
		time.Sleep(5 * time.Second)
	}
	wg.Done()
}

func startReceivingTime(conn net.Conn, wg *sync.WaitGroup) {
	p := make([]byte, 2048)
	for {
		_, err := bufio.NewReader(conn).Read(p)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		sync, err := time.ParseDuration(string(p))
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		startTime = startTime.Add(sync)
		fmt.Println("Time Synchronized")
		//time.Sleep(5 * time.Second)
	}
	defer wg.Done()
}

func printTime(wg *sync.WaitGroup) {
	for {
		fmt.Println("Local time is", startTime)
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
