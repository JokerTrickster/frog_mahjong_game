package find_it

import (
	"context"
	"fmt"
	"main/features/find_it/model/entity"
	_errors "main/features/find_it/model/errors"
	"main/features/find_it/repository"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func CancelMatchEventWebsocket(msg *entity.WSMessage) *entity.ErrorInfo {
	ctx := context.Background()
	uID := msg.UserID
	rID := msg.RoomID

	//비즈니스 로직
	preloadUsers := []entity.PreloadUsers{}
	messageMsg := entity.MessageInfo{}
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
		messageMsg.GameInfo.IsFull = true
	}

	message, err := CreateMessage(&messageMsg)
	if err != nil {
		return CreateErrorMessage(_errors.ErrCodeBadRequest, _errors.ErrMarshalFailed, "메시지 생성 에러")
	}
	msg.Message = message
	sendMessageToClients(rID, msg)
	return nil
}
