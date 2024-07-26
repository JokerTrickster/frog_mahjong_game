package utils

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	WSUpgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
		return true
	}}
	WSClients   = make(map[uint]map[*websocket.Conn]WSClient)
	WSBroadcast = make(chan WSMessage)
)

type WSClient struct {
	RoomID uint
	UserID uint
	Conn   *websocket.Conn
}

type WSMessage struct {
	Message string `json:"message"`
	Event   string `json:"event"`
	RoomID  uint   `json:"roomID"`
	UserID  uint   `json:"userID"`
}

func WSHandleMessages() {

	for {
		msg := <-WSBroadcast

		switch msg.Event {
		case "OPEN": // 방 생성
			if clients, ok := WSClients[msg.RoomID]; ok {
				for client := range clients {
					err := client.WriteJSON(msg)
					fmt.Println(msg)
					if err != nil {
						log.Printf("error: %v", err)
						client.Close()
						delete(clients, client)
					}
				}
			}
		case "CLOSE": // 방 나가기
			if clients, ok := WSClients[msg.RoomID]; ok {
				for client := range clients {
					if msg.UserID == clients[client].UserID {
						client.Close()
						delete(clients, client)
					} else {
						err := client.WriteJSON(msg)
						if err != nil {
							log.Printf("error: %v", err)
							client.Close()
							delete(clients, client)
						}
					}
				}
			}
		case "MESSAGE": // 메시지 전송
			if clients, ok := WSClients[msg.RoomID]; ok {
				for client := range clients {
					err := client.WriteJSON(msg)
					if err != nil {
						log.Printf("error: %v", err)
						client.Close()
						delete(clients, client)
					}
				}
			}
		}
	}
}
