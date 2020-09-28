package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var startTime time.Time

func printTime() {
	for {
		fmt.Println("Local time is", time.Since(startTime))
		time.Sleep(5 * time.Second)
	}
}

func printErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}

func askTime(ip string, c chan string) {
	conn, err := net.Dial("udp", ip)
	defer conn.Close()
	printErr(err)
	fmt.Fprintf(conn, "Send your local time")
	p := make([]byte, 1024)
	_, err = bufio.NewReader(conn).Read(p)
	printErr(err)
	p = bytes.Trim(p, "\x00")
	c <- string(p)
}

func synchronizeClocks(slavesTime []string) time.Duration {
	masterTime := time.Since(startTime)
	sumDiff := masterTime - masterTime
	for _, slave := range slavesTime {
		slaveDur, err := time.ParseDuration(slave)
		printErr(err)
		sumDiff += masterTime - slaveDur
	}
	syncTime := masterTime + sumDiff/time.Duration(len(slavesTime)+1)
	return syncTime
}

func sendTime(ips []string, syncTime string) {
	for _, ip := range ips {
		conn, err := net.Dial("udp", ip)
		defer conn.Close()
		printErr(err)
		fmt.Fprintf(conn, syncTime)
	}
}

func startMaster() {
	local := os.Args[3]
	localCl, err := time.ParseDuration("-" + local + "s")
	printErr(err)
	startTime = time.Now().Add(localCl)
	slavesData, err := ioutil.ReadFile(os.Args[4])
	printErr(err)
	slavesID := strings.Split(string(slavesData), "\n")

	go printTime()
	time.Sleep(5 * time.Second)
	for {
		c := make(chan string)
		for _, ip := range slavesID {
			go askTime(ip, c)
		}
		var slavesTime []string
		for range slavesID {
			data := <-c
			slavesTime = append(slavesTime, data)
		}
		syncTime := synchronizeClocks(slavesTime)
		sync := time.Since(startTime) - syncTime
		startTime = startTime.Add(sync)
		sendTime(slavesID, syncTime.String())
		time.Sleep(10 * time.Second)
	}
}

func startSlave() {
	ip := strings.Split(os.Args[2], ":")
	local := os.Args[3]
	localCl, err := time.ParseDuration("-" + local + "s")
	printErr(err)
	startTime = time.Now().Add(localCl)
	port, err := strconv.Atoi(ip[1])
	printErr(err)
	addr := net.UDPAddr{Port: port, IP: net.ParseIP(ip[0])}
	conn, err := net.ListenUDP("udp", &addr)
	printErr(err)
	defer conn.Close()

	go printTime()
	for {
		p := make([]byte, 1024)
		n, remoteaddr, err := conn.ReadFromUDP(p)
		printErr(err)
		_ = p[:n]
		localTime := time.Since(startTime).String()
		_, err = conn.WriteToUDP([]byte(localTime), remoteaddr)

		q := make([]byte, 1024)
		n, _, err = conn.ReadFromUDP(q)
		q = q[:n]
		syncTime, err := time.ParseDuration(string(q))
		printErr(err)
		sync := time.Since(startTime) - syncTime
		startTime = startTime.Add(sync)
	}
}

func main() {
	if os.Args[1] == "-m" {
		startMaster()
	} else {
		startSlave()
	}
}
