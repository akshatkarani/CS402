package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var startTime time.Time

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

func startSendingTime(conn net.Conn, wg *sync.WaitGroup) {
	for {
		fmt.Fprintf(conn, getTimeString(startTime))
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
		t, _ := strconv.ParseFloat(string(p), 64)
		syncFl := t - getTime(startTime)
		sync, err := time.ParseDuration(strconv.FormatFloat(syncFl, 'f', -1, 64) + "s")
		if err != nil {
			fmt.Println("Error", err)
		}
		fmt.Println(startTime, sync)
		startTime = startTime.Add(sync)
		fmt.Println(startTime)
		// fmt.Println("Time Synchronized")
		time.Sleep(1 * time.Second)
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
