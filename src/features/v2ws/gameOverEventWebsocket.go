package v2ws

import (
	"context"
	"fmt"
	"main/features/v2ws/model/entity"
	_errors "main/features/v2ws/model/errors"
	"main/features/v2ws/repository"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func GameOverEventWebsocket(msg *entity.WSMessage) {
	//유저 상태를 변경한다. (대기실로 이동)
	ctx := context.Background()
	uID := msg.UserID
	roomID := msg.RoomID

	GameOverEntity := entity.WSGameOverEntity{
		RoomID: roomID,
		UserID: uID,
	}
	// 비즈니스 로직
	roomInfoMsg := entity.RoomInfo{}
	preloadUsers := []entity.RoomUsers{}
	var err error
	err = mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {

		// 유저 상태 변경
		err = repository.GameOverUpdateRoomUsers(ctx, tx, &GameOverEntity)
		if err != nil {
			return err
		}
		preloadUsers, err = repository.GameOverFindAllRoomUsers(ctx, tx, roomID)
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
	}

	// 메시지 생성
	roomInfoMsg = *CreateRoomInfoMSG(ctx, preloadUsers, 1, roomInfoMsg.ErrorInfo,0)
	roomInfoMsg.GameInfo.AllReady = false

	// 구조체를 JSON 문자열로 변환 (마샬링)
	message, err := CreateMessage(&roomInfoMsg)
	if err != nil {
		fmt.Println(err)
	}
	msg.Message = message

	//유저 상태를 변경한다. (방에 참여)
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
