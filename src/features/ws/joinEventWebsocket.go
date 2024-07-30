package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"main/features/ws/model/entity"
	"main/features/ws/model/request"
	"main/features/ws/repository"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func JoinEventWebsocket(msg *entity.WSMessage) {
	ctx := context.Background()
	uID := msg.UserID
	req := request.ReqWSJoin{
		RoomID: uint(msg.RoomID),
	}
	if msg.Message != "" {
		req.Password = msg.Message
	}

	err := mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		// 방 참여 가능한지 체크
		RoomDTO, err := repository.FindOneRoom(ctx, tx, &req)
		if err != nil {
			return err
		}
		if RoomDTO.CurrentCount == RoomDTO.MaxCount {
			return err
		}
		// 방 유저 정보를 생성한다.
		RoomUserDTO, err := CreateRoomUserDTO(uID, int(req.RoomID), "wait")
		if err != nil {
			return err
		}
		err = repository.InsertOneRoomUser(ctx, tx, RoomUserDTO)
		if err != nil {
			return err
		}
		// 방 현재 인원을 증가시킨다.
		err = repository.FindOneAndUpdateRoom(ctx, tx, req.RoomID)
		if err != nil {
			return err
		}

		//유저 정보를 업데이트 한다.
		err = repository.FindOneAndUpdateUser(ctx, tx, uID, req.RoomID)
		if err != nil {
			return err
		}
		return nil
	})
	fmt.Println(err)

	// 메시지 생성
	// 현재 참여하고 있는 유저에 대한 정보를 가져와서 메시지 전달한다.
	preloadUsers, err := repository.FindAllRoomUsers(ctx, uint(msg.RoomID))
	if err != nil {
		log.Println(err)
	}
	roomUsersMsg := entity.TRoomUsers{}
	for _, roomUser := range preloadUsers {
		user := entity.User{
			ID:          uint(roomUser.UserID),
			PlayerState: roomUser.PlayerState,
			Coin:        roomUser.User.Coin,
			Name:        roomUser.User.Name,
			Email:       roomUser.User.Email,
		}
		roomUsersMsg.Users = append(roomUsersMsg.Users, user)
	}
	// 구조체를 JSON 문자열로 변환 (마샬링)
	jsonData, err := json.Marshal(roomUsersMsg)
	if err != nil {
		log.Fatalf("JSON 마샬링 에러: %s", err)
	}

	// JSON 바이트 배열을 문자열로 변환
	jsonString := string(jsonData)
	msg.Message = jsonString

	//유저 상태를 변경한다. (방에 참여)
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
