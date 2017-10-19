package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"gitlab.fuzzhq.com/Web-Ops/ufo/pkg/ufo"
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

func sendDeploymentStatus(w http.ResponseWriter, r *http.Request, s *AppState, UFO *ufo.UFO) {
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
		currentService, err := UFO.GetService(s.c, *s.s.ServiceName)

		HandleError(err)

		if msg.IsDeployed {
			conn.Close()
			break
		}
		if int(*currentService.DesiredCount) > 0 {
			runningTasks, err := UFO.RunningTasks(s.c, currentService)

			HandleError(err)

			if len(runningTasks) > 0 {
				tasks, err := UFO.GetTasks(s.c, runningTasks)

				for _, task := range tasks.Tasks {
					if *task.TaskDefinitionArn == *s.newT.TaskDefinitionArn && *task.LastStatus == "RUNNING" {
						msg.IsDeployed = true
					}
				}

				err = conn.WriteJSON(msg.IsDeployed)

				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}
	fmt.Println("Client unsubscribed")
}
