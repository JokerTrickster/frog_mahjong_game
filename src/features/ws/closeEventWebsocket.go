package ws

import (
	"log"
	"main/features/ws/model/entity"
)

func CloseEventWebsocket(msg *entity.WSMessage) {
	//유저 상태를 변경한다. (대기실로 이동)
	if clients, ok := entity.WSClients[msg.RoomID]; ok {
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
}
