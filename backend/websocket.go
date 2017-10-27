package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

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

func IsDeployed(UFO *ufo.UFO, s *AppState) bool {
	runningTasks, err := UFO.RunningTasks(s.c, s.s)

	HandleError(err)

	tasks, err := UFO.GetTasks(s.c, runningTasks)

	HandleError(err)

	for _, task := range tasks.Tasks {
		if *task.TaskDefinitionArn == *s.newT.TaskDefinitionArn && *task.LastStatus == "RUNNING" {
			return true
		}
	}

	return false
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

	HandleError(err)

	for !IsDeployed(UFO, s) {
		err = conn.WriteJSON(false)
		HandleWriteError(err)

		time.Sleep(waitTime)
	}

	err = conn.WriteJSON(true)
	HandleWriteError(err)

	fmt.Println("Client unsubscribed")
}

func HandleWriteError(err error) {
	if err != nil {
		log.Println("write:", err)
		return
	}
}
