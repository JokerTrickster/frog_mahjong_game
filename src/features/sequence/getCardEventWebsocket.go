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

func GetCardEventWebsocket(msg *entity.WSMessage) *entity.ErrorInfo {
	//유저 상태를 변경한다. (대기실로 이동)
	ctx := context.Background()
	uID := msg.UserID
	roomID := msg.RoomID

	// 비즈니스 로직
	//해당 방이 대기상태인지 체크한다.
	preloadUsers := []entity.PreloadUsers{}
	messageMsg := entity.MessageInfo{}
	var errInfo *entity.ErrorInfo

	err := mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		// 남은 카드수가 0이면 버린 카드 상태값을 모두 초기화 한다.
		dummyCardCount, errInfo := repository.GetCardCountDummyCard(ctx, tx, roomID)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}
		if dummyCardCount == 0 {
			errInfo = repository.GetCardUpdateDummyCard(ctx, tx, roomID)
			if errInfo != nil {
				return fmt.Errorf("%s", errInfo.Msg)
			}
		}

		// 남은 카드 하나를 가져온다.
		cardInfo, errInfo := repository.GetCardFindOneCardInfo(ctx, tx, roomID, uID)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}
		// 유저에게 카드 정보를 업데이트 한다.
		errInfo = repository.GetCardUpdateCardState(ctx, tx, roomID, uID, int(cardInfo.CardID))
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}

		// 라운드 수를 증가시킨다. 남은 카드 수를 감소시킨다.
		errInfo = repository.GetCardUpdateRoomSetting(ctx, tx, roomID)
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
