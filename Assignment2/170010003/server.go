package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/rpc"
	"os"
	"strconv"
	"strings"
	"time"
)

// Listener : int type
type Listener int

type Response struct {
	Data    string
	TransNo string
}

type Request struct {
	Data    int
	TransNo string
}

// TransIds : List of all the processed IDs
var TransIds []string
var index int

func readBalance() string {
	balance, err := ioutil.ReadFile("Balance.txt")
	if err != nil {
		fmt.Println(err.Error())
	}
	return strings.TrimSpace(string(balance))
}

func writeBalance(balance string) {
	balance += "\n"
	err := ioutil.WriteFile("Balance.txt", []byte(balance), 0777)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func checkTransID(transNo string) bool {
	for i := range TransIds {
		if TransIds[i] == transNo {
			return true
		}
	}
	return false
}

func writeFile() {
	file, err := os.OpenFile("Trans_Processed.txt", os.O_WRONLY, os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	for _, line := range TransIds {
		_, err = file.WriteString(line + "\n")
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

// GetBalance : Get the current balance
func (l *Listener) GetBalance(args *Request, response *Response) error {
	// Implement the logic to check the transaction already processed
	// if processed add appropriate message into Response
	if checkTransID(args.TransNo) {
		response.Data = "The transaction is already processed"
		response.TransNo = args.TransNo
		fmt.Println("The transaction is already processed")
		return nil
	}

	// Else Implement the logic to add the transaction id into
	// TransIds array and Trans_Processed.txt file
	TransIds = append(TransIds, args.TransNo)
	writeFile()

	// Implement the logic to read the Balance from Balance.txt
	balance := readBalance()
	fmt.Println("Read balance successfully")

	// Wait for 3 seconds
	// time.Sleep(3 * time.Second)

	// Implement the logic to add Balance and Transaction ID to the response
	response.Data = balance
	response.TransNo = args.TransNo
	return nil
}

// DepositeAmount : Desposit amount and update the balance
func (l *Listener) DepositeAmount(args *Request, response *Response) error {
	// Implement the logic to check the transaction already processed
	// if processed add appropriate message into Response
	if checkTransID(args.TransNo) {
		response.Data = "The transaction is already processed"
		response.TransNo = args.TransNo
		fmt.Println("The transaction is already processed")
		return nil
	}

	// Else Implement the logic to add the transaction id into
	// TransIds array and Trans_Processed.txt file
	TransIds = append(TransIds, args.TransNo)
	writeFile()

	// Implement the logic to read the balance and add the amount to into balance
	// and write calculated balance back to Balance.txt
	balance := readBalance()
	currBalance, err := strconv.Atoi(balance)
	if err != nil {
		fmt.Println(err.Error())
	}
	newBalance := strconv.Itoa(args.Data + currBalance)
	writeBalance(newBalance)
	fmt.Println("Updated balance successfully")

	// Wait for 3 seconds
	time.Sleep(3 * time.Second)

	// Implement the logic add appropriate message and transaction number to response\
	response.Data = "Your Amount is deposited into the account successfully"
	response.TransNo = args.TransNo
	return nil
}

func readFile() {
	//Read the transaction processed from Trans_Processed.txt into TransIds
	file, err := os.Open("Trans_Processed.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		TransIds = append(TransIds, strings.TrimSpace(scanner.Text()))
	}
}

func main() {
	//Read the transaction processed from Trans_Processed.txt into TransIds
	readFile()

	//Sample code to start the server, read the port number from command-line
	listener := new(Listener)
	rpc.Register(listener)
	rpc.HandleHTTP()

	port := os.Args[1]
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}
