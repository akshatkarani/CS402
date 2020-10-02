package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/DistributedClocks/GoVector/govec"
)

var startTime time.Time
var logger *govec.GoLog

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

func askTime(ips string, conn *net.UDPConn, c chan string) {
	ip := strings.Split(ips, ":")
	port, err := strconv.Atoi(ip[1])
	printErr(err)
	addr := net.UDPAddr{Port: port, IP: net.ParseIP(ip[0])}
	data := logger.PrepareSend("Time request sent", "Send your local time", govec.GetDefaultLogOptions())
	conn.WriteToUDP(data, &addr)
	p := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(p)
	var recv string
	logger.UnpackReceive("Time request received", p, recv, govec.GetDefaultLogOptions())
	printErr(err)
	// p = p[:n]
	// p = bytes.Trim(p, "\x00")
	c <- string(recv)
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

func sendTimeToOne(ips string, syncTime string, conn *net.UDPConn) {
	ip := strings.Split(ips, ":")
	port, err := strconv.Atoi(ip[1])
	printErr(err)
	addr := net.UDPAddr{Port: port, IP: net.ParseIP(ip[0])}
	data := logger.PrepareSend("Sending Synchronized time", syncTime, govec.GetDefaultLogOptions())
	conn.WriteToUDP(data, &addr)
}

func sendTime(ips []string, syncTime string, conn *net.UDPConn) {
	for _, ip := range ips {
		go sendTimeToOne(ip, syncTime, conn)
	}
}

func startMaster() {
	ip := strings.Split(os.Args[2], ":")
	port, err := strconv.Atoi(ip[1])
	printErr(err)
	addr := net.UDPAddr{Port: port, IP: net.ParseIP(ip[0])}
	conn, err := net.ListenUDP("udp", &addr)

	local := os.Args[3]
	localCl, err := time.ParseDuration("-" + local + "s")
	printErr(err)
	startTime = time.Now().Add(localCl)
	slavesData, err := ioutil.ReadFile(os.Args[4])
	printErr(err)
	slavesID := strings.Split(strings.TrimSpace(string(slavesData)), "\n")

	logger = govec.InitGoVector("Master", os.Args[5], govec.GetDefaultConfig())
	go printTime()
	time.Sleep(5 * time.Second)
	for {
		c := make(chan string)
		for _, ip := range slavesID {
			go askTime(ip, conn, c)
		}
		var slavesTime []string
		for range slavesID {
			data := <-c
			slavesTime = append(slavesTime, data)
		}
		syncTime := synchronizeClocks(slavesTime)
		sync := time.Since(startTime) - syncTime
		startTime = startTime.Add(sync)
		logger.LogLocalEvent("Local Time Synchronized", govec.GetDefaultLogOptions())
		sendTime(slavesID, syncTime.String(), conn)
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

	logger = govec.InitGoVector("Slave: "+ip[1], os.Args[4], govec.GetDefaultConfig())
	go printTime()
	for {
		p := make([]byte, 1024)
		n, remoteaddr, err := conn.ReadFromUDP(p)
		printErr(err)
		logger.UnpackReceive("Time request received", p, p, govec.GetDefaultLogOptions())
		_ = p[:n]
		localTime := time.Since(startTime).String()
		_ = logger.PrepareSend("Local time sent", localTime, govec.GetDefaultLogOptions())
		_, err = conn.WriteToUDP([]byte(localTime), remoteaddr)
		printErr(err)

		q := make([]byte, 1024)
		n, _, err = conn.ReadFromUDP(q)
		q = q[:n]
		syncTime, err := time.ParseDuration(string(q))
		printErr(err)
		logger.LogLocalEvent("Synchronized time received", govec.GetDefaultLogOptions())
		sync := time.Since(startTime) - syncTime
		startTime = startTime.Add(sync)
		logger.LogLocalEvent("Local time synchronized", govec.GetDefaultLogOptions())
	}
}

func main() {
	if os.Args[1] == "-m" {
		if len(os.Args) != 6 {
			fmt.Println("Usage: go run MasterClient.go -m ip:port time slavesfile logfile")
			return
		}
		startMaster()
	} else if os.Args[1] == "-s" {
		if len(os.Args) != 5 {
			fmt.Println("Usage: go run MasterClient.go -s ip:port time logfile")
			return
		}
		startSlave()
	} else {
		fmt.Println("Wrong usage: Either run as -m or -s")
	}
}
