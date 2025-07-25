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

func HeroEventWebsocket(msg *entity.WSMessage) *entity.ErrorInfo {
	//유저 상태를 변경한다. (대기실로 이동)
	ctx := context.Background()
	uID := msg.UserID
	roomID := msg.RoomID
	req := request.ReqWSHero{}
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
		// // 카드 정보를 가져온다.
		// cardInfo, errInfo := repository.HeroFindOneCardInfo(ctx, tx, roomID, req.CardID)
		// if errInfo != nil {
		// 	return fmt.Errorf("%s", errInfo.Msg)
		// }

		// 왕을 이동시킨다. 라운드 수를 증가시킨다.
		errInfo = repository.HeroUpdateKing(ctx, tx, roomID, req.KingIndex)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}

		// 왕 이동 자리에 유저 슬라임을 놓는다.
		errInfo = repository.HeroUpdateUserSlime(ctx, tx, roomID, uID, req.KingIndex)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}

		// 유저가 사용한 카드 상태값을 discard 로 변경한다.
		errInfo = repository.HeroUpdateCardState(ctx, tx, roomID, uID, req.CardID)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}

		// 유저 히어로 카드 카운트 감소한다.
		errInfo = repository.HeroUpdateUserHeroCardDecrease(ctx, tx, roomID, uID)
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
	if messageMsg.Users[0].ID == uID {
		messageMsg.Users[0].LastCardID = req.CardID
	} else {
		messageMsg.Users[1].LastCardID = req.CardID
	}
	message, err := CreateMessage(&messageMsg)
	if err != nil {
		return CreateErrorMessage(_errors.ErrCodeBadRequest, _errors.ErrMarshalFailed, "메시지 생성 에러")
	}
	msg.Message = message
	sendMessageToClients(roomID, msg)
	return nil
}
