package ws

import (
	"context"
	"main/features/ws/model/entity"
	_errors "main/features/ws/model/errors"
	"main/features/ws/model/request"
	"main/features/ws/repository"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func ChatEventWebsocket(msg *entity.WSMessage) *entity.ErrorInfo {
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
	msg.Message = req.Message
	msg.ChatID = ChatID
	sendMessageToClients(roomID, msg)
	return nil
}
