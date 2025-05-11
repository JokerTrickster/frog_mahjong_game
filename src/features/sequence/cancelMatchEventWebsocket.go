package sequence

import (
	"context"
	"fmt"
	"main/features/sequence/model/entity"
	_errors "main/features/sequence/model/errors"
	"main/features/sequence/repository"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func CancelMatchEventWebsocket(msg *entity.WSMessage) *entity.ErrorInfo {
	if msg == nil {
		return CreateErrorMessage(_errors.ErrCodeBadRequest, _errors.ErrInvalidRequest, "잘못된 메시지 형식입니다.")
	}

	ctx := context.Background()
	uID := msg.UserID
	rID := msg.RoomID

	// Check if room exists before proceeding
	if _, exists := entity.RoomSessions[rID]; !exists {
		return CreateErrorMessage(_errors.ErrCodeBadRequest, _errors.ErrRoomNotFound, "방이 이미 삭제되었습니다.")
	}

	//비즈니스 로직
	preloadUsers := []entity.PreloadUsers{}
	messageMsg := entity.MessageInfo{
		SequenceGameInfo: &entity.SequenceGameInfo{
			RoomID: rID,
		},
		Users: make([]*entity.User, 0),
	}
	var errInfo *entity.ErrorInfo
	err := mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		// 방 유저 정보를 삭제한다.
		errInfo = repository.CancelMatchDeleteRoomUser(ctx, tx, uID)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}
		// 방 세팅 정보 삭제한다.
		errInfo = repository.CancelMatchDeleteRoomSetting(ctx, tx, rID)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}

		// 방 정보를 삭제한다.
		errInfo = repository.CancelMatchDeleteRoom(ctx, tx, rID)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}

		return nil
	})
	if err != nil {
		return errInfo
	}

	// 메시지 생성
	messageMsg = *CreateMessageInfoMSG(ctx, preloadUsers, 1, messageMsg.ErrorInfo, 0)

	if len(preloadUsers) == 2 {
		if messageMsg.SequenceGameInfo != nil {
			messageMsg.SequenceGameInfo.IsFull = true
		}
	}

	message, err := CreateMessage(&messageMsg)
	if err != nil {
		return CreateErrorMessage(_errors.ErrCodeBadRequest, _errors.ErrMarshalFailed, "메시지 생성 에러")
	}
	msg.Message = message

	// Clean up websocket connections before sending message
	if sessionIDs, ok := entity.RoomSessions[rID]; ok {
		for _, sessionID := range sessionIDs {
			if client, exists := entity.WSClients[sessionID]; exists {
				client.Closed = true
				delete(entity.WSClients, sessionID)
			}
		}
		delete(entity.RoomSessions, rID)
	}

	// Send message only if there are still active connections
	if sessionIDs, ok := entity.RoomSessions[rID]; ok && len(sessionIDs) > 0 {
		sendMessageToClients(rID, msg)
	}

	return nil
}
