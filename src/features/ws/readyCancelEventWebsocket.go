package ws

import (
	"context"
	"encoding/json"
	"log"
	"main/features/ws/model/entity"
	"main/features/ws/repository"
)

func ReadyCancelEventWebsocket(msg *entity.WSMessage) {
	//유저 상태를 변경한다. (대기실로 이동)
	ctx := context.Background()
	uID := msg.UserID
	roomID := msg.RoomID

	err := repository.ReadyCancelFindOneAndUpdateRoomUser(ctx, uID, roomID)
	if err != nil {
		log.Println(err)
	}

	// 메시지 생성
	// 현재 참여하고 있는 유저에 대한 정보를 가져와서 메시지 전달한다.
	preloadUsers, err := repository.ReadyCancelFindAllRoomUsers(ctx, uint(msg.RoomID))
	if err != nil {
		log.Println(err)
	}
	roomInfoMsg := entity.RoomInfo{}
	for _, roomUser := range preloadUsers {
		user := entity.User{
			ID:          uint(roomUser.UserID),
			PlayerState: roomUser.PlayerState,
			Coin:        roomUser.User.Coin,
			Name:        roomUser.User.Name,
			Email:       roomUser.User.Email,
		}
		if roomUser.Room.OwnerID == roomUser.UserID {
			user.IsOwner = true
		}
		roomInfoMsg.Users = append(roomInfoMsg.Users, &user)
	}
	roomInfoMsg.GameInfo = &entity.GameInfo{
		PlayTurn: 1,
		AllReady: false,
	}

	// 구조체를 JSON 문자열로 변환 (마샬링)
	jsonData, err := json.Marshal(roomInfoMsg)
	if err != nil {
		log.Fatalf("JSON 마샬링 에러: %s", err)
	}

	// JSON 바이트 배열을 문자열로 변환
	jsonString := string(jsonData)
	msg.Message = jsonString

	//유저 상태를 변경한다. (방에 참여)
	if clients, ok := entity.WSClients[msg.RoomID]; ok {
		for client := range clients {
			//나머지 유저에게 메시지 전달
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}