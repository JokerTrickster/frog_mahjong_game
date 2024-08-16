package ws

import (
	"context"
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
	//비즈니스 로직
	roomInfoMsg := entity.RoomInfo{}
	preloadUsers := []entity.RoomUsers{}
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
		preloadUsers, err = repository.CloseFindAllRoomUsers(ctx, tx, req.RoomID)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		roomInfoMsg.ErrorInfo = &entity.ErrorInfo{
			Code: 500,
			Msg:  err.Error(),
		}
	}

	// 메시지 생성
	roomInfoMsg = *CreateRoomInfoMSG(ctx, preloadUsers, 1, roomInfoMsg.ErrorInfo)
	// 구조체를 JSON 문자열로 변환 (마샬링)
	message, err := CreateMessage(&roomInfoMsg)
	if err != nil {
		fmt.Println(err)
	}
	msg.Message = message

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
