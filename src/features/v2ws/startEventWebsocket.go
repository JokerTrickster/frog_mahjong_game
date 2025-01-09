package v2ws

import (
	"context"
	"fmt"
	"main/features/v2ws/model/entity"
	_errors "main/features/v2ws/model/errors"
	"main/features/v2ws/repository"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func StartEventWebsocket(msg *entity.WSMessage) *entity.ErrorInfo {
	//유저 상태를 변경한다. (대기실로 이동)
	ctx := context.Background()
	uID := msg.UserID
	roomID := msg.RoomID

	// 비즈니스 로직
	//해당 방이 대기상태인지 체크한다.
	preloadUsers := []entity.RoomUsers{}
	roomInfoMsg := entity.RoomInfo{}
	var errInfo *entity.ErrorInfo
	roomState, newErr := repository.StartCheckRoomState(ctx, roomID)
	if newErr != nil {
		return newErr
	}
	if roomState != "wait" {
		return CreateErrorMessage(_errors.ErrCodeBadRequest, "게임이 시작되었습니다.", _errors.ErrAlreadyGame)
	}

	err := mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		// 방장이 게임 시작 요청했는지 체크
		errInfo := repository.StartCheckOwner(ctx, tx, uID, roomID)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}

		// room 데이터 값 변경 (상태 변경, 시작 시간 추가)
		roomUpdateData := StartUpdateRoom(roomID)
		errInfo = repository.StartUpdateRoom(ctx, tx, roomID, roomUpdateData)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}

		// 유저들 코인 -1 차감한다.
		errInfo = repository.StartDiffCoin(ctx, tx, roomID)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}

		// 기존 카드가 있다면 모두 제거한다.
		errInfo = repository.StartDeleteCards(ctx, tx, uID)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}
		birdCards, err := repository.StartBirdCard(ctx, tx)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}
		// 카드를 생성한다.
		cards := CreateInitCards(roomID, birdCards)
		err = repository.StartCreateCards(ctx, tx, cards)
		if err != nil {
			return fmt.Errorf("%s", err.Msg)
		}

		return nil
	})

	if err != nil {
		return errInfo
	}
	preloadUsers, errInfo = repository.StartFindAllRoomUsers(ctx, roomID)
	if errInfo != nil {
		fmt.Println(newErr)
		return errInfo
	}
	// 메시지 생성
	roomInfoMsg = *CreateRoomInfoMSG(ctx, preloadUsers, 1, roomInfoMsg.ErrorInfo, 0)
	openCards, errInfo := repository.StartUpdateCardState(ctx, roomID)
	if errInfo != nil {
		fmt.Println(newErr)
		return errInfo
	}
	roomInfoMsg.GameInfo.OpenCards = openCards

	message, err := CreateMessage(&roomInfoMsg)
	if err != nil {
		return CreateErrorMessage(_errors.ErrCodeInternal, err.Error(), _errors.ErrGameTerminated)
	}
	msg.Message = message
	msg.SessionID = ""
	sendMessageToClients(roomID, msg)
	return nil
}
