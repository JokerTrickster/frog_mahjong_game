package slime_war

import (
	"context"
	"encoding/json"
	"fmt"
	"main/features/slime_war/model/entity"
	_errors "main/features/slime_war/model/errors"
	"main/features/slime_war/model/request"
	"main/features/slime_war/repository"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func JoinPlayEventWebsocket(msg *entity.WSMessage) *entity.ErrorInfo {
	ctx := context.Background()
	uID := msg.UserID

	//string to struct
	req := request.ReqWSJoinPlayEvent{}
	err := json.Unmarshal([]byte(msg.Message), &req)
	if err != nil {
		return CreateErrorMessage(_errors.ErrCodeBadRequest, _errors.ErrUnmarshalFailed, "JSON 언마샬링 에러")
	}

	//비즈니스 로직
	messageInfoMsg := entity.MessageInfo{}
	preloadUsers := []entity.PreloadUsers{}
	var errInfo *entity.ErrorInfo
	roomID, errInfo := repository.JoinPlayFindOneRoomUsers(ctx, uID)
	if errInfo != nil {
		return errInfo
	}

	err = mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		//유저 정보를 업데이트 한다.
		errInfo = repository.JoinPlayFindOneAndUpdateUser(ctx, tx, uID, roomID)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}

		preloadUsers, errInfo = repository.PreloadUsers(ctx, tx, roomID)
		if err != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}
		return nil
	})
	if err != nil {
		return errInfo
	}
	// 메시지 생성
	messageInfoMsg = *CreateMessageInfoMSG(ctx, preloadUsers, 1, messageInfoMsg.ErrorInfo, 0)
	if len(preloadUsers) == 1 {
		messageInfoMsg.SlimeWarGameInfo.IsFull = false
		messageInfoMsg.SlimeWarGameInfo.AllReady = false
	}

	message, err := CreateMessage(&messageInfoMsg)
	if err != nil {
		return CreateErrorMessage(_errors.ErrCodeBadRequest, _errors.ErrMarshalFailed, "메시지 생성 에러")
	}
	msg.Message = message
	sendMessageToClients(roomID, msg)

	return nil
}
