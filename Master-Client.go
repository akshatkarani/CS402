package main

import (
	"os"
	"fmt"
	"ioutil"
	"time"
	"net"
	"strconv"
)

func main() {
	if os.Args[1] == "-m" {
		// Take IP Port number as input
		CONNECT := os.Args[2]

		// The process local time is set using the input
		time_delay := os.Args[3]
		offSet, _ := time.ParseDuration(os.Args[3])
		t := time.Now().Add(offSet)
		fmt.Println(t)

		// Read IP:Port mapping of all the slaves from file
		content, _ := ioutil.ReadFile(os.Args[4])

		CONNECT_IDS := strings.Split(string(content), "\n")

		// Array containing all the time differences
		var times []float
		// Get time from each slave
		for {
			for _, CONNECT = range CONNECT_IDS {
				// The master acts as a client and tries connecting to
				//  slaves to get its process local time
				s, err := net.ResolveUDPAddr("udp4", CONNECT)
				c, err := net.DialUDP("udp4", nil, s)
				if err != nil {
					fmt.Println(err)
					return
				}

				// Close connection if error is encountered
				defer c.Close()

				// Send message requesting slave to send time
				data := []byte("WHAT IS YOUR TIME?\n")
				_, err = c.Write(data)

				if err != nil {
					fmt.Println(err)
					return
				}

				// Read reply sent by Slave
				buffer := make([]byte, 1024)
				n, _, err := c.ReadFromUDP(buffer)
				if err != nil {
					fmt.Println(err)
					return
				}

				// Print the time delay replied by slave
				fmt.Printf("Reply: %s\n", string(buffer[0:n]))

				// Calculate average over master and slave
				num2, _ := strconv.ParseFloat(string(buffer[0:n])[1:4], 8)
				times = append(times, num2)
				
			}
			// Iterate over all the times of slaves and calculate average
			i := 1
			sumDelay := strconv.ParseFloat(time_delay[1:4], 8)
			for _, time := range  times {
				i = i + 1
				sumDelay = sumDelay + time
			}
			sumDelay = sumDelay/i

			// Convert float to string and set new time delay as avg delay
			string_data := strconv.FormatFloat(sumDelay, 'g', 3, 64)
			time_delay = "+" + string_data +"h"
			offSet, _ := time.ParseDuration(time_delay)
			
			// Write back average time to all the slaves
			for _, CONNECT = range CONNECT_IDS {
				// The master acts as a client and tries connecting to 
				//  slaves to send average time
				s, err := net.ResolveUDPAddr("udp4", CONNECT)
				c, err := net.DialUDP("udp4", nil, s)
				if err != nil {
					fmt.Println(err)
					return
				}

				// Close connection if error is encountered
				defer c.Close()

				data = []byte(string_data)
				_, err = c.Write(data)
			}
		
			time.Sleep(5*time.Second)
		}

	}else {
		// Take IP Port number as input
		CONNECT := os.Args[2]

		// The process local time is set using the input
		time_delay := os.Args[3]
		offSet, _ := time.ParseDuration(os.Args[3])

		// The slave listens to messages from the master 
		s, err := net.ResolveUDPAddr("udp4", CONNECT)
        if err != nil {
			fmt.Println(err)
			return
        }
		
        connection, err := net.ListenUDP("udp4", s)
        if err != nil {
			fmt.Println(err)
			return
		}

		// Close connection if error is encountered
		defer connection.Close()
        buffer := make([]byte, 1024)
		
		// Wait for MASTER to query the SLAVE's time
		for {
			n, addr, err := connection.ReadFromUDP(buffer)
			fmt.Println("-> ", string(buffer[0:n-1]))

			t := time.Now().Add(offSet).Format(time.RFC1123)
			fmt.Printf("The current time is: %s\n", t)

			// The slave replies the time difference between it's process 
			// clock and system clock
			data := []byte(time_delay)
			_, err = connection.WriteToUDP(data, addr)
			if err != nil {
				fmt.Println(err)
				return
			}

			// Read the reply of average time from master
			n, addr, err = connection.ReadFromUDP(buffer)
			fmt.Println("-> ", string(buffer[0:n]))
			time_delay = "+" + string(buffer[0:n]) + "h"

			time.Sleep(5*time.Second)
		}

	}
}