package find_it

import (
	"context"
	"encoding/json"
	"fmt"
	"main/features/find_it/model/entity"
	_errors "main/features/find_it/model/errors"
	"main/features/find_it/model/request"
	"main/features/find_it/repository"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NextRoundEventWebsocket(msg *entity.WSMessage) *entity.ErrorInfo {
	//유저 상태를 변경한다. (대기실로 이동)
	ctx := context.Background()
	uID := msg.UserID
	roomID := msg.RoomID
	req := request.ReqWSNextRound{}
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
		// 방장이 라운드 요청했는지 체크
		errInfo := repository.StartCheckOwner(ctx, tx, uID, roomID)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}

		// 라운드 수 증가
		errInfo = repository.NextRoundRoundIncrease(ctx, tx, roomID, req.Round+1)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}
		//

		//TODO 30라운드 이미지를 선택해서 각 라운드마다 이미지를 만든다.
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

	if len(preloadUsers) == 2 {
		messageMsg.GameInfo.IsFull = true
	}

	message, err := CreateMessage(&messageMsg)
	if err != nil {
		return CreateErrorMessage(_errors.ErrCodeBadRequest, _errors.ErrMarshalFailed, "메시지 생성 에러")
	}
	if messageMsg.GameInfo.Round == messageMsg.GameInfo.RoundCount+1 {
		msg.Event = "GAME_CLEAR"
	} else {
		msg.Event = "ROUND_START"
	}
	msg.Message = message
	sendMessageToClients(roomID, msg)
	return nil
}
