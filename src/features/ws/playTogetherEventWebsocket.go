package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"main/features/ws/model/entity"
	_errors "main/features/ws/model/errors"
	"main/features/ws/model/request"
	"main/features/ws/repository"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func PlayTogetherEventWebsocket(msg *entity.WSMessage) {
	ctx := context.Background()
	uID := msg.UserID

	//string to struct
	req := request.ReqWSPlayTogetherEvent{}
	err := json.Unmarshal([]byte(msg.Message), &req)
	if err != nil {
		log.Fatalf("JSON 언마샬링 에러: %s", err)
	}

	//비즈니스 로직
	roomInfoMsg := entity.RoomInfo{}
	preloadUsers := []entity.RoomUsers{}
	roomID, err := repository.PlayTogetherFindOneRoomUsers(ctx, uID)
	if err != nil {
		log.Fatalf("방 유저 정보 조회 에러: %s", err)
	}
	err = mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		// 방 정보를 업데이트 한다. (타이머, 인원 수)
		err := repository.PlayTogetherFindOneAndUpdateRoom(ctx, tx, roomID, uint(req.Count), uint(req.Timer))
		if err != nil {
			return err
		}

		//유저 정보를 업데이트 한다.
		err = repository.PlayTogetherFindOneAndUpdateUser(ctx, tx, uID, roomID)
		if err != nil {
			return err
		}

		preloadUsers, err = repository.PlayTogetherFindAllRoomUsers(ctx, tx, roomID)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		roomInfoMsg.ErrorInfo = &entity.ErrorInfo{
			Code: 500,
			Msg:  err.Error(),
			Type: _errors.ErrInternalServer,
		}
		if roomInfoMsg.ErrorInfo.Msg == "방이 꽉 찼습니다." {
			roomInfoMsg.ErrorInfo.Type = _errors.ErrRoomFull
		} else if roomInfoMsg.ErrorInfo.Msg == "비밀번호가 일치하지 않습니다." {
			roomInfoMsg.ErrorInfo.Type = _errors.ErrWrongPassword
		} else if roomInfoMsg.ErrorInfo.Msg == "게임 중인 방입니다." {
			roomInfoMsg.ErrorInfo.Type = _errors.ErrGameInProgress
		}
	}

	// 메시지 생성
	roomInfoMsg = *CreateRoomInfoMSG(ctx, preloadUsers, 1, roomInfoMsg.ErrorInfo)
	roomInfoMsg.GameInfo.AllReady = false

	if len(preloadUsers) == req.Count {
		roomInfoMsg.GameInfo.IsFull = true
		roomInfoMsg.GameInfo.AllReady = true
	}
	// 구조체를 JSON 문자열로 변환 (마샬링)
	message, err := CreateMessage(&roomInfoMsg)
	if err != nil {
		fmt.Println(err)
	}
	msg.Message = message
	msg.RoomID = roomID

	//방 유저들에게 메시지 전달
	if clients, ok := entity.WSClients[msg.RoomID]; ok {
		//에러 발생시 이벤트 요청한 유저에게만 메시지를 전달한다.
		if roomInfoMsg.ErrorInfo != nil || err != nil {
			for client := range clients {
				if clients[client].UserID == msg.UserID {
					_ = client.WriteJSON(msg)
					clientData := clients[client]
					clientData.Close()
					clients[client] = clientData
					delete(clients, client)
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