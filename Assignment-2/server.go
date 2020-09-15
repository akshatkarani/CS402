package main

import (
	"time"
)

type Listener int

type Response struct {
	Data    string
	TransNo string
}

type Request struct {
	Data      int
	TransData string
}

var TransIds [100]string
var index int

func (l *Listener) GetBalance(args *Request, response *Response) error {

	//Implement the logic to check the transaction already processed if processed add appropriate message into Response

	//Else Implement the logic to add the transaction id into TransIds array and Trans_Processed.txt file

	//Implement the logic to read the Balance from Balance.txt

	//Wait for 3 seconds
	time.Sleep(3 * time.Second)
	//Implement the logic to add Balance and Transaction ID to the response

	return nil
}

func (l *Listener) DepositeAmount(args *Request, response *Response) error {

	//Implement the logic to check the transaction already processed if processed add appropriate message into Response

	//Else Implement the logic to add the transaction id into TransIds array and Trans_Processed.txt file

	//Implement the logic to read the balance and add the amount to into balance and write calculated balance back to Balance.txt

	time.Sleep(3 * time.Second)

	//Implement the logic add appropriate message and transaction number to response

	return nil
}

func main() {

	//Read the transaction processed from Trans_Processed.txt into TransIds

	//Sample code to start the server, read the port number from command-line
	/*
		listener := new(Listener)
		rpc.Register(listener)
		rpc.HandleHTTP()

		err = http.ListenAndServe(":8098", nil)
		if err != nil {
			fmt.Println(err.Error())
		}*/

}
