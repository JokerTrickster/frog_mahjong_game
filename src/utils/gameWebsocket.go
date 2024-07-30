package utils

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	WSUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
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

/*
유저 ID

	이메일
	이름
	유저 상태 : ready or not ready
	방장인지 여부 :
	가지고 있는 패 정보들 :
	버린 패 정보들 :
	현재 보유하고 있는 코인 :
*/
type RoomUsers struct {
	Users []User
}
type User struct {
	ID             uint   `json:"id"`
	Email          string `json:"email"`
	Name           string `json:"name"`
	PlayerState    string `json:"playerState"`
	IsOwner        bool   `json:"isOwner"`
	Cards          []Card `json:"cards"`
	DiscardedCards []Card `json:"discardedCards"`
	Coin           int    `json:"coin"`
}

/*
카드 ID
이름 : oen, two, three, four, five, six, seven, eight, nine , chung, bal
색깔 : green, red, normal
상태 : 버려진 패 or 소유하고 있는 패 or 가운데 놓여져 있는 패
유저 ID
*/
type Card struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Color  string `json:"color"`
	State  string `json:"state"`
	UserID uint   `json:"userID"`
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
					msg.Message = "방에 참여했습니다."
					err := client.WriteJSON(msg)
					if err != nil {
						log.Printf("error: %v", err)
						client.Close()
						delete(clients, client)
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
