package main

import (
	"cli/delivery/deliveryparam"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	fmt.Println("command", os.Args[0])
	message := "default message"
	if len(os.Args) > 1 {
		message = os.Args[1]
	}
	fmt.Println("message", message)
	connection, err := net.Dial("tcp", "127.0.0.1:1379")
	if err != nil {
		log.Fatalln("can't dial the given address", err)
	}
	defer connection.Close()

	fmt.Println("local address", connection.LocalAddr())
	req := deliveryparam.Request{Command: message}
	if req.Command == "create-task" {
		req.CreateTaskRequest = deliveryparam.CreateTaskRequest{
			Title:      "tes5t",
			DueDate:    "test",
			CategoryID: 1,
		}
	}
	serializedData, merr := json.Marshal(&req)
	if merr != nil {
		log.Fatalln("failed to marshal", merr)
	}

	numOfWrittenBytes, werr := connection.Write([]byte(serializedData))
	if werr != nil {
		log.Fatalln("can't write data to connection", werr)
	}
	fmt.Println("num of written bytes", numOfWrittenBytes)

	var data = make([]byte, 1024)
	_, rErr := connection.Read(data)
	if rErr != nil {
		log.Fatalln("can't read data from connection", rErr)
	}
	fmt.Println("server response", string(data))

}
