package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const waitTime = 2 * time.Second

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Message ...
type Message struct {
	IsDeployed bool `json:"isDeployed"`
}

func sendDeploymentStatus(w http.ResponseWriter, r *http.Request, cluster string, service string) {
	// Upgrade initial GET request to a websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure we close the connection when the function returns
	defer conn.Close()
	msg := Message{IsDeployed: false}
	fmt.Println("Client subscribed")
	for {
		time.Sleep(waitTime)
		currentService := *describeService(cluster, service)
		if msg.IsDeployed {
			conn.Close()
			break
		}
		if int(*currentService.DesiredCount) > 0 {
			runningTasks := listRunningTasks(cluster, service)
			if len(runningTasks) > 0 {
				tasks := describeTasks(cluster, runningTasks)
				for _, task := range tasks.Tasks {
					if *task.TaskDefinitionArn == TaskDefinition && *task.LastStatus == "RUNNING" {
						msg.IsDeployed = true
					}
				}
				err := conn.WriteJSON(msg.IsDeployed)
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}
	fmt.Println("Client unsubscribed")
}
