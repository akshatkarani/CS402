package main

import (
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
	Data      int
	TransData string
}

func main() {

	//Sample code to connect to server using HTTP

	//client, err := rpc.DialHTTP("tcp", serverAddress+":8098")

	//Implement the logic to generate the transaction id using date and time and construct the request.

	//Sample code to make async call to server
	/*
		response := new(Response)
		divCall := client.Go("Listener.GetBalance", args, response, nil)
		if divCall == nil {
			log.Fatal(divCall)
		}
	*/

	time.Sleep(1 * time.Second)

	//Implement the logic to generate duplicate transaction

	//Sample code to make async call to server
	/*
		response2 := new(Response)
		divCall = client.Go("Listener.GetBalance", args, response2, nil)

		if divCall == nil {
			log.Fatal(divCall)
		}*/

	//Wait for the transanctions complete
	time.Sleep(6 * time.Second)

	//Display the response
	/*log.Printf("Response1: Transaction Number - %v : %v", response.TransNo, response.Data)
	log.Printf("Response2: Transaction Number - %v : %v",response2.TransNo, response2.Data)
	*/

	//Implement rest of the transactions

}
