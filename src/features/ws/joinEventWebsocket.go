package ws

import (
	"context"
	"fmt"
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
	//비즈니스 로직
	roomInfoMsg := entity.RoomInfo{}
	preloadUsers := []entity.RoomUsers{}
	err := mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		// 방 참여 가능한지 체크
		RoomDTO, err := repository.JoinFindOneRoom(ctx, tx, &req)
		if err != nil {
			return err
		}
		if RoomDTO.Password != req.Password {
			return fmt.Errorf("비밀번호가 일치하지 않습니다.")
		}

		if RoomDTO.CurrentCount == RoomDTO.MaxCount {
			return fmt.Errorf("방이 꽉 찼습니다.")
		}
		// TODO 기존에 방 유저 정보가 있는지 가져온다.
		// 유저 정보가 있으면 삭제하고 방 인원수를 감소시킨다.
		err = repository.JoinFindOneAndDeleteRoomUser(ctx, tx, uID, req.RoomID)
		if err != nil {
			return err
		}

		// 방 유저 정보를 생성한다.
		RoomUserDTO, err := CreateRoomUserDTO(uID, int(req.RoomID), "wait")
		if err != nil {
			return err
		}
		err = repository.JoinInsertOneRoomUser(ctx, tx, RoomUserDTO)
		if err != nil {
			return err
		}
		// 방 현재 인원을 증가시킨다.
		err = repository.JoinFindOneAndUpdateRoom(ctx, tx, req.RoomID)
		if err != nil {
			return err
		}

		//유저 정보를 업데이트 한다.
		err = repository.JoinFindOneAndUpdateUser(ctx, tx, uID, req.RoomID)
		if err != nil {
			return err
		}
		preloadUsers, err = repository.JoinFindAllRoomUsers(ctx, tx, req.RoomID)
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
	roomInfoMsg = *CreateRoomInfoMSG(ctx, preloadUsers, 1)
	roomInfoMsg.GameInfo.AllReady = false
	// 구조체를 JSON 문자열로 변환 (마샬링)
	message, err := CreateMessage(&roomInfoMsg)
	if err != nil {
		fmt.Println(err)
	}
	msg.Message = message

	//방 유저들에게 메시지 전달
	if clients, ok := entity.WSClients[msg.RoomID]; ok {
		//에러 발생시 이벤트 요청한 유저에게만 메시지를 전달한다.
		if roomInfoMsg.ErrorInfo != nil || err != nil {
			for client := range clients {
				if clients[client].UserID == msg.UserID {
					err := client.WriteJSON(msg)
					if err != nil {
						fmt.Printf("error: %v", err)
						client.Close()
						delete(clients, client)
					}
				}
			}
		} else {
			for client := range clients {
				err := client.WriteJSON(msg)
				if err != nil {
					fmt.Printf("error: %v", err)
					client.Close()
					delete(clients, client)
				}
			}
		}
	}
}
