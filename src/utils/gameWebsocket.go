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
		fmt.Println("현재 방 인원수 : ", len(WSClients[msg.RoomID]))
		switch msg.Event {
		case "OPEN": // 방 생성
			// 유저 상태를 변경한다.

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
		case "JOIN": // 방 참여
			//유저 상태를 변경한다. (대기실 -> 방으로 이동)
			if clients, ok := WSClients[msg.RoomID]; ok {
				for client := range clients {
					if msg.UserID == clients[client].UserID {
						client.Close()
						delete(clients, client)
					} else {
						msg.Message = "방을 나갔습니다."
						err := client.WriteJSON(msg)
						if err != nil {
							log.Printf("error: %v", err)
							client.Close()
							delete(clients, client)
						}
					}
				}
			}
		case "CLOSE": // 방 나가기
			//유저 상태를 변경한다. (대기실로 이동)
			if clients, ok := WSClients[msg.RoomID]; ok {
				for client := range clients {
					if msg.UserID == clients[client].UserID {
						client.Close()
						delete(clients, client)
					} else {
						msg.Message = "방을 나갔습니다."
						err := client.WriteJSON(msg)
						if err != nil {
							log.Printf("error: %v", err)
							client.Close()
							delete(clients, client)
						}
					}
				}
			}

		case "READY": // 게임 준비
			// 유저 정보를 ready로 변경
			// 모든 플레이어가 ready 상태라면 방장에게 게임 시작 가능 메시지 전송

			if clients, ok := WSClients[msg.RoomID]; ok {
				for client := range clients {
					if msg.UserID == clients[client].UserID {
						msg.Message = "준비 완료"
					} else {
						msg.Message = "준비를 했습니다."
					}
					err := client.WriteJSON(msg)
					if err != nil {
						log.Printf("error: %v", err)
						client.Close()
						delete(clients, client)
					}
				}
			}
		case "CANCEL": // 게임 준비를 취소
			//유저 정보를 ready로 변경
			// 모든 플레이어가 ready였다가 취소한거라면 방장에게 게임 시작 불가능 메시지 전송

			if clients, ok := WSClients[msg.RoomID]; ok {
				for client := range clients {
					if msg.UserID == clients[client].UserID {
						msg.Message = "준비 취소"
					} else {
						msg.Message = "준비를 취소 했습니다."
					}
					err := client.WriteJSON(msg)
					if err != nil {
						log.Printf("error: %v", err)
						client.Close()
						delete(clients, client)
					}
				}
			}
		case "START": // 게임 시작
			// 모든 플레이어에게 게임 시작을 알리고
			// 게임 준비를 한다.

			if clients, ok := WSClients[msg.RoomID]; ok {
				for client := range clients {

					msg.Message = "게임 시작"
					err := client.WriteJSON(msg)
					if err != nil {
						log.Printf("error: %v", err)
						client.Close()
						delete(clients, client)
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
