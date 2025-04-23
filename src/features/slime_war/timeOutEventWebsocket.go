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

func TimeOutEventWebsocket(msg *entity.WSMessage) *entity.ErrorInfo {
	//유저 상태를 변경한다. (대기실로 이동)
	ctx := context.Background()
	roomID := msg.RoomID
	req := request.ReqWSTimeOut{}
	err := json.Unmarshal([]byte(msg.Message), &req)
	if err != nil {
		return CreateErrorMessage(_errors.ErrCodeBadRequest, _errors.ErrUnmarshalFailed, "JSON 언마샬링 에러")
	}
	// 비즈니스 로직
	//해당 방이 대기상태인지 체크한다.
	preloadUsers := []entity.PreloadUsers{}
	messageMsg := entity.MessageInfo{}
	var errInfo *entity.ErrorInfo

	err = mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		//현재 턴을 상대방 턴으로 넘긴다.
		errInfo = repository.TimeOutUpdateTurn(ctx, tx, roomID)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}

		preloadUsers, errInfo = repository.PreloadUsers(ctx, tx, roomID)
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

	message, err := CreateMessage(&messageMsg)
	if err != nil {
		return CreateErrorMessage(_errors.ErrCodeBadRequest, _errors.ErrMarshalFailed, "메시지 생성 에러")
	}

	msg.Message = message
	sendMessageToClients(roomID, msg)
	return nil
}
