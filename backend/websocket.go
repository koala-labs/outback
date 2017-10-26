package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/gorilla/websocket"
	"gitlab.fuzzhq.com/Web-Ops/ufo/pkg/ufo"
)

var waitTime = 2 * time.Second

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func IsDeployed(tasks *ecs.DescribeTasksOutput, s *AppState) {
	for _, task := range tasks.Tasks {
		if *task.TaskDefinitionArn == *s.newT.TaskDefinitionArn && *task.LastStatus == "RUNNING" {
			s.IsDeployed = true
		}
	}
}

func PollForStatus(w http.ResponseWriter, r *http.Request, UFO *ufo.UFO, s *AppState) {
	// Upgrade initial GET request to a websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure we close the connection when the function returns
	defer conn.Close()

	fmt.Println("Client subscribed")
	for {
		serviceDetail, err := UFO.GetService(s.c, *s.s.ServiceName)

		HandleError(err)

		runningTasks, err := UFO.RunningTasks(s.c, serviceDetail)

		HandleError(err)

		if s.IsDeployed {
			conn.Close()
			break
		}

		if len(runningTasks) > 0 {
			tasks, err := UFO.GetTasks(s.c, runningTasks)

			HandleError(err)

			IsDeployed(tasks, s)

			err = conn.WriteJSON(s.IsDeployed)

			if err != nil {
				fmt.Println(err)
				return
			}
		}

		time.Sleep(waitTime)
	}

	fmt.Println("Client unsubscribed")
}
