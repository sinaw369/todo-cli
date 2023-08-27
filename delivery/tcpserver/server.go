package main

import (
	"cli/delivery/deliveryparam"
	"cli/repository/memorystore"
	"cli/service/task"
	"encoding/json"
	"fmt"
	"log"
	"net"
)

func main() {
	const (
		network = "tcp"
		address = "127.0.0.1:1379"
	)

	//create new listener
	listener, err := net.Listen(network, address)
	if err != nil {
		log.Fatalln("can't listen on given address:", address, err)
	}
	defer listener.Close()

	fmt.Println("server listening on:", listener.Addr())
	taskMemoryRepo := memorystore.NewTaskStore()
	taskService := task.NewService(taskMemoryRepo)

	for {
		//listen for new connections
		connection, err := listener.Accept()
		if err != nil {
			log.Println("can't listen to new  connection:", err)

			continue
		}
		//process request
		var rawRequest = make([]byte, 1024)
		numberOfReadBytes, rErr := connection.Read(rawRequest)
		if rErr != nil {
			fmt.Println("can't read data from connection ", err)

			continue
		}
		fmt.Printf("client Address: %s ,numOfReadBytes: %d,data: %s\n", connection.RemoteAddr(), numberOfReadBytes, string(rawRequest))

		req := &deliveryparam.Request{}
		if err := json.Unmarshal(rawRequest[:numberOfReadBytes], req); err != nil {
			log.Println("bad request", err)
			continue
		}
		switch req.Command {
		case "create-task":
			response, err := taskService.Create(task.CreateRequest{
				Title:               req.CreateTaskRequest.Title,
				DueDate:             req.CreateTaskRequest.DueDate,
				CategoryID:          req.CreateTaskRequest.CategoryID,
				AuthenticatedUserID: 0,
			})
			if err != nil {
				_, wErr := connection.Write([]byte(err.Error()))
				if wErr != nil {
					log.Println("can't write data to connection ", rErr)

					continue
				}
			}
			data, mErr := json.Marshal(&response)
			if mErr != nil {
				_, wErr := connection.Write([]byte(mErr.Error()))
				if wErr != nil {
					log.Println("can't marshal response  ", rErr)

					continue
				}
				continue
			}
			_, wErr := connection.Write(data)
			if wErr != nil {
				log.Println("can't write data to connection ", rErr)

				continue
			}

		}

		connection.Close()
	}

}
