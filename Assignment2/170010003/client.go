package main

import (
	"fmt"
	"log"
	"net/rpc"
	"os"
	"time"
)

type Response struct {
	Data    string
	TransNo string
}

type Request struct {
	Data    int
	TransNo string
}

func timestamp() string {
	dt := time.Now()
	const format = "01-02-2006 15:04:05 Monday"
	return dt.Format(format)
}

func main() {
	// Sample code to connect to server using HTTP
	ip := os.Args[1]
	port := os.Args[2]
	client, err := rpc.DialHTTP("tcp", ip+":"+port)
	if err != nil {
		fmt.Println(err.Error())
	}

	/*
		No Duplicate Entry
	*/

	// Request 1
	response := new(Response)
	args := new(Request)
	args.TransNo = timestamp()
	divCall := client.Go("Listener.GetBalance", args, response, nil)
	if divCall == nil {
		log.Fatal(divCall)
	}
	time.Sleep(3 * time.Second)
	fmt.Printf("Response1: Transaction Number - %v : %v\n", response.TransNo, response.Data)

	// Request 2
	response = new(Response)
	args = new(Request)
	args.TransNo = timestamp()
	divCall = client.Go("Listener.GetBalance", args, response, nil)
	if divCall == nil {
		log.Fatal(divCall)
	}
	time.Sleep(3 * time.Second)
	fmt.Printf("Response2: Transaction Number - %v : %v\n", response.TransNo, response.Data)

	/*
		Duplicate Entry
	*/

	// Request 3
	response = new(Response)
	args = new(Request)
	args.TransNo = timestamp()
	divCall = client.Go("Listener.GetBalance", args, response, nil)
	if divCall == nil {
		log.Fatal(divCall)
	}
	time.Sleep(3 * time.Second)
	fmt.Printf("Response3: Transaction Number - %v : %v\n", response.TransNo, response.Data)

	// Request 4
	response = new(Response)
	divCall = client.Go("Listener.GetBalance", args, response, nil)
	if divCall == nil {
		log.Fatal(divCall)
	}
	time.Sleep(3 * time.Second)
	fmt.Printf("Response4: Transaction Number - %v : %v\n", response.TransNo, response.Data)

	/*
		Deposit the Amount
	*/

	// Request 5
	response = new(Response)
	args = new(Request)
	args.TransNo = timestamp()
	args.Data = 3000
	divCall = client.Go("Listener.DepositeAmount", args, response, nil)
	if divCall == nil {
		log.Fatal(divCall)
	}
	time.Sleep(5 * time.Second)
	fmt.Printf("Response5: Transaction Number - %v : %v\n", response.TransNo, response.Data)

	// Request 6
	response = new(Response)
	args = new(Request)
	args.TransNo = timestamp()
	divCall = client.Go("Listener.GetBalance", args, response, nil)
	if divCall == nil {
		log.Fatal(divCall)
	}
	time.Sleep(3 * time.Second)
	fmt.Printf("Response6: Transaction Number - %v : %v\n", response.TransNo, response.Data)

	/*
		Server Crash
	*/

	// Request 7
	response = new(Response)
	args = new(Request)
	args.TransNo = timestamp()
	args.Data = 3000
	divCall = client.Go("Listener.DepositeAmount", args, response, nil)
	if divCall == nil {
		log.Fatal(divCall)
	}
	time.Sleep(5 * time.Second)
	fmt.Printf("Response7: Transaction Number - %v : %v\n", response.TransNo, response.Data)

	fmt.Println("Delay of 10secs")
	time.Sleep(10 * time.Second)
	fmt.Println("Establishing Connection again")

	// Dial up connection again since server is restarted
	ip = os.Args[1]
	port = os.Args[2]
	client, err = rpc.DialHTTP("tcp", ip+":"+port)
	if err != nil {
		fmt.Println(err.Error())
	}

	// Request 8
	response = new(Response)
	divCall = client.Go("Listener.DepositeAmount", args, response, nil)
	if divCall == nil {
		log.Fatal(divCall)
	}
	time.Sleep(3 * time.Second)
	fmt.Printf("Response8: Transaction Number - %v : %v\n", response.TransNo, response.Data)

	// Request 9
	response = new(Response)
	args = new(Request)
	args.TransNo = timestamp()
	divCall = client.Go("Listener.GetBalance", args, response, nil)
	if divCall == nil {
		log.Fatal(divCall)
	}
	time.Sleep(3 * time.Second)
	fmt.Printf("Response9: Transaction Number - %v : %v\n", response.TransNo, response.Data)

}
