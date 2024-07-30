package ws

import (
	"fmt"
	"log"
	"main/features/ws/model/entity"
)

func WSHandleMessages() {

	for {
		msg := <-entity.WSBroadcast
		fmt.Println("현재 방 인원수 : ", len(entity.WSClients[msg.RoomID]))
		switch msg.Event {
		case "JOIN": // 방 참여
			JoinEventWebsocket(&msg)
		case "CLOSE": // 방 나가기
			CloseEventWebsocket(&msg)
		case "READY": // 게임 준비
			ReadyEventWebsocket(&msg)
		case "CANCEL": // 게임 준비를 취소
			//유저 정보를 ready로 변경
			// 모든 플레이어가 ready였다가 취소한거라면 방장에게 게임 시작 불가능 메시지 전송

			if clients, ok := entity.WSClients[msg.RoomID]; ok {
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

			if clients, ok := entity.WSClients[msg.RoomID]; ok {
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
			if clients, ok := entity.WSClients[msg.RoomID]; ok {
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
