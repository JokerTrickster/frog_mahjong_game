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

func CloseEventWebsocket(msg *entity.WSMessage) {
	//유저 상태를 변경한다. (대기실로 이동)
	ctx := context.Background()
	uID := msg.UserID
	req := request.ReqWSJoin{
		RoomID: uint(msg.RoomID),
	}
	if msg.Message != "" {
		req.Password = msg.Message
	}

	err := mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		// RoomsID에 해당하는 userID를 삭제한다.
		err := repository.CloseFindOneAndDeleteRoomUser(ctx, tx, uID, req.RoomID)
		if err != nil {
			return err
		}
		// Rooms 현재 인원수를 -1한다.
		roomDTO, err := repository.CloseFindOneAndUpdateRoom(ctx, tx, req.RoomID)
		if err != nil {
			return err
		}
		// user에 rooms_id를 1로 바꾸고 state를 wait으로 변경한다.
		err = repository.CloseFindOneAndUpdateUser(ctx, tx, uID)
		if err != nil {
			return err
		}

		if roomDTO.CurrentCount == 0 {
			// 방 삭제
			err = repository.CloseFindOneAndDeleteRoom(ctx, tx, req.RoomID)
			if err != nil {
				return err
			}

		} else if roomDTO.CurrentCount == 1 {
			// 인원이 1명이면 남아 있는 유저를 방장으로 변경
			//방장이 나가면 다른 유저 중 한명을 방장으로 변경
			//룸에 남아있는 유저 정보를 가져온다.
			roomUserDTO, err := repository.CloseFindOneRoomUser(ctx, tx, req.RoomID)
			if err != nil {
				return err
			}
			userDTO, err := repository.CloseFindOneUser(ctx, tx, uint(roomUserDTO.UserID))
			if err != nil {
				return err
			}

			// 해당 유저를 방장으로 업데이트 한다.

			//방장으로 변경하기 위해 업데이트해야 될 부분들
			// rooms -> owner 변경
			err = repository.CloseChangeRoomOnwer(ctx, tx, req.RoomID, userDTO.ID)
			if err != nil {
				return err
			}

		}
		return nil
	})
	fmt.Println(err)

	// 메시지 생성
	// 현재 참여하고 있는 유저에 대한 정보를 가져와서 메시지 전달한다.
	preloadUsers, err := repository.CloseFindAllRoomUsers(ctx, uint(msg.RoomID))
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
			//방나간 유저 클로즈 처리
			if clients[client].UserID == msg.UserID {
				client.Close()
				delete(clients, client)
			} else {
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
}
