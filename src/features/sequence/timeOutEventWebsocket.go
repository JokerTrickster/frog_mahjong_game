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
	"math/rand"

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
		// 현재 소유하고 있는 카드 중 하나를 랜덤으로 사용한다.
		userCards, errInfo := repository.TimeOutFindUserCards(ctx, tx, roomID, uint(req.UserID))
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}
		randomCard := userCards[rand.Intn(len(userCards))]
		errInfo = repository.TimeOutUpdateCardState(ctx, tx, int(roomID), int(randomCard.CardID))
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}
		//  해당 카드의 맵에 사용 여부를 표시한다
		errInfo = repository.TimeOutUpdateMapState(ctx, tx, int(roomID), req.UserID, int(randomCard.CardID))
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}
		// 더미에서 카드 한장을 가져온다.
		errInfo = repository.TimeOutUpdateDummyCardState(ctx, tx, int(roomID), req.UserID)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}
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
