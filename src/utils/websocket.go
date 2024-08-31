package utils

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	Upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
		return true
	}}
	Clients   = make(map[*websocket.Conn]Client)
	Broadcast = make(chan Message)
)

type Client struct {
	Name   string
	UserID int
	Conn   *websocket.Conn
}

type Message struct {
	Username string `json:"username"`
	Message  string `json:"message"`
	Event    string `json:"event"`
}

func HandleMessages() {

	for {
		msg := <-Broadcast
		switch msg.Event {
		case "OPEN":
			for client := range Clients {
				err := client.WriteJSON(msg)
				if err != nil {
					log.Printf("error: %v", err)
					client.Close()
					delete(Clients, client)
				}
			}
		case "MESSAGE":
			for client := range Clients {
				err := client.WriteJSON(msg)
				if err != nil {
					log.Printf("error: %v", err)
					client.Close()
					delete(Clients, client)
				}
			}
		}
	}
}
