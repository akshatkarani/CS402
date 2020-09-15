package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// opening a UDP connection
	conn, err := net.Dial("udp", "10.250.1.65:8090")
	if err != nil {
		fmt.Printf("Error %v", err)
		return
	}

	// take input from user
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Please give input: ")
	text, _ := reader.ReadString('\n')

	// send message to server
	fmt.Fprintf(conn, text)

	// read response from server
	p := make([]byte, 4096)
	_, err = bufio.NewReader(conn).Read(p)

	// print the response
	if err == nil {
		fmt.Printf("%s", p)
	} else {
		fmt.Printf("Error %v\n", err)
	}
	conn.Close()
}
