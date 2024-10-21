package v2ws

import (
	"context"
	"log"
	"main/features/v2ws/model/entity"
	_errors "main/features/v2ws/model/errors"
	"main/features/v2ws/model/request"
	"main/features/v2ws/repository"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func ChatEventWebsocket(msg *entity.WSMessage) {
	//유저 상태를 변경한다. (대기실로 이동)
	ctx := context.Background()
	uID := msg.UserID
	roomID := msg.RoomID
	//string to struct
	req := request.ReqWSChat{
		UserID:  uID,
		RoomID:  roomID,
		Name:    msg.Name,
		Message: msg.Message,
	}

	// 비즈니스 로직
	ChatInfo := entity.ChatInfo{}
	var ChatID uint
	err := mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {

		// 채팅 DTO 만들기
		chatDTO := CreateChatDTO(req)

		// 채팅 정보 저장
		chatID, err := repository.ChatInsertOneChat(ctx, tx, chatDTO)
		if err != nil {
			return err
		}
		ChatID = chatID

		return nil
	})
	if err != nil {
		ChatInfo.ErrorInfo = &entity.ErrorInfo{
			Code: 500,
			Msg:  err.Error(),
			Type: _errors.ErrInternalServer,
		}
	}

	//유저 상태를 변경한다. (방에 참여)
	if clients, ok := entity.WSClients[msg.RoomID]; ok {
		// 메시지 생성

		//에러 발생시 이벤트 요청한 유저에게만 메시지를 전달한다.
		if ChatInfo.ErrorInfo != nil || err != nil {
			for client := range clients {
				if clients[client].UserID == msg.UserID {
					// 구조체를 JSON 문자열로 변환 (마샬링)
					msg.Message = req.Message
					msg.ChatID = ChatID

					err = client.WriteJSON(msg)
					if err != nil {
						log.Printf("error: %v", err)
						client.Close()
						delete(clients, client)
					}
				}
			}
		} else {
			for client := range clients {
				// 구조체를 JSON 문자열로 변환 (마샬링)
				msg.Message = req.Message
				msg.ChatID = ChatID
				err = client.WriteJSON(msg)
				if err != nil {
					log.Printf("error: %v", err)
					client.Close()
					delete(clients, client)
				}
			}
		}
	}
}
