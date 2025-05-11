package sequence

import (
	"context"
	"encoding/json"
	"fmt"
	"main/features/sequence/model/entity"
	_errors "main/features/sequence/model/errors"
	"main/features/sequence/model/request"
	"main/features/sequence/repository"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func UseCardEventWebsocket(msg *entity.WSMessage) *entity.ErrorInfo {
	//유저 상태를 변경한다. (대기실로 이동)
	ctx := context.Background()
	userID := msg.UserID
	roomID := msg.RoomID
	req := request.ReqWSUseCard{}
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
		// 맵에 표시를 한다.
		errInfo = repository.UseCardUpdateMapState(ctx, tx, int(roomID), int(userID), req.MapID)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}
		// 유저 카드를 사용처리 한다.
		errInfo = repository.UseCardUpdateCardState(ctx, tx, roomID, userID, req.CardID)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}

		// 더미에서 카드를 한장 가져온다.
		errInfo = repository.UseCardGetDummyCard(ctx, tx, roomID, userID)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}
		//현재 턴을 상대방 턴으로 넘긴다.
		errInfo = repository.UseCardUpdateTurn(ctx, tx, roomID)
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
