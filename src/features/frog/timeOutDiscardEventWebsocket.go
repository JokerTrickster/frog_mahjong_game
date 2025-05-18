package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"main/features/frog/model/entity"
	_errors "main/features/frog/model/errors"
	"main/features/frog/model/request"
	"main/features/frog/repository"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func TimeOutDiscardCardsEventWebsocket(msg *entity.WSMessage) *entity.ErrorInfo {
	//유저 상태를 변경한다. (대기실로 이동)
	ctx := context.Background()
	uID := msg.UserID
	roomID := msg.RoomID

	//string to struct
	req := request.ReqWSTimeOutDiscardCards{}
	err := json.Unmarshal([]byte(msg.Message), &req)
	if err != nil {
		return CreateErrorMessage(_errors.ErrCodeBadRequest, _errors.ErrUnmarshalFailed, "JSON 언마샬링 에러")
	}
	TimeOutDiscardCardsEntity := entity.WSTimeOutDiscardCardsEntity{
		RoomID: roomID,
		UserID: uID,
		CardID: uint(req.CardID),
	}

	// 비즈니스 로직
	roomInfoMsg := entity.RoomInfo{}
	preloadUsers := []entity.RoomUsers{}
	var errInfo *entity.ErrorInfo
	err = mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		// 카드 상태 없데이트
		errInfo := repository.TimeOutDiscardUpdateCardState(ctx, tx, &TimeOutDiscardCardsEntity)
		if err != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}

		preloadUsers, errInfo = repository.PreloadFindGameInfo(ctx, tx, roomID)
		if err != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}
		return nil
	})
	if err != nil {
		return errInfo
	}
	// 메시지 생성
	gameRoomSettings, errInfo := repository.FindOneFrogCurrentRound(ctx, roomID)
	if errInfo != nil {
		return errInfo
	}
	roomInfoMsg = *CreateRoomInfoMSG(ctx, preloadUsers, gameRoomSettings.CurrentRound+1, roomInfoMsg.ErrorInfo)

	// 론 가능 여부를 true로 변경
	roomInfoMsg.GameInfo.IsLoanAllowed = true
	message, err := CreateMessage(&roomInfoMsg)
	if err != nil {
		return CreateErrorMessage(_errors.ErrCodeInternal, err.Error(), _errors.ErrGameTerminated)
	}
	msg.Message = message
	sendMessageToClients(msg.RoomID, msg)
	return nil
}
