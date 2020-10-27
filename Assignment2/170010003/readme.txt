client.go: This file contains client program
server.go: This file contains server program
Balance.txt: This file stores the current Balance
Trans_Processed.txt: This file stores all the processed ids.
output.png: This is screenshot for output of all 4 questions

First start the server and then start the client using
`go run server.go <port>`
`go run client.go <ip> <port>`

To test the server crash functionality, when server prints "Updated balance successfully" immediately manually crash the server.
Then restart it immediately.