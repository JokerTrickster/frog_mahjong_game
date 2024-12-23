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

func StartEventWebsocket(msg *entity.WSMessage) {
	//유저 상태를 변경한다. (대기실로 이동)
	ctx := context.Background()
	uID := msg.UserID
	roomID := msg.RoomID

	// 비즈니스 로직
	//해당 방이 대기상태인지 체크한다.
	preloadUsers := []entity.RoomUsers{}
	roomInfoMsg := entity.RoomInfo{}
	roomState, newErr := repository.StartCheckRoomState(ctx, roomID)
	if newErr != nil {
		roomInfoMsg.ErrorInfo = newErr
		ErrorHandling(msg, &roomInfoMsg)
		return
	}
	if roomState != "wait" {
		roomInfoMsg.ErrorInfo = &entity.ErrorInfo{
			Code: _errors.ErrCodeBadRequest,
			Msg:  "게임이 시작되었습니다.",
			Type: _errors.ErrAlreadyGame,
		}
		ErrorHandling(msg, &roomInfoMsg)
		return
	}

	err := mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		// 방장이 게임 시작 요청했는지 체크
		err := repository.StartCheckOwner(ctx, tx, uID, roomID)
		if err != nil {
			roomInfoMsg.ErrorInfo = err
			ErrorHandling(msg, &roomInfoMsg)
			return fmt.Errorf("%s", err.Msg)
		}

		// room 데이터 값 변경 (상태 변경, 시작 시간 추가)
		roomUpdateData := StartUpdateRoom(roomID)
		err = repository.StartUpdateRoom(ctx, tx, roomID, roomUpdateData)
		if err != nil {
			roomInfoMsg.ErrorInfo = err
			ErrorHandling(msg, &roomInfoMsg)
			return fmt.Errorf("%s", err.Msg)
		}

		// 유저들 코인 -1 차감한다.
		err = repository.StartDiffCoin(ctx, tx, roomID)
		if err != nil {
			roomInfoMsg.ErrorInfo = err
			ErrorHandling(msg, &roomInfoMsg)
			return fmt.Errorf("%s", err.Msg)
		}

		// 기존 카드가 있다면 모두 제거한다.
		err = repository.StartDeleteCards(ctx, tx, uID)
		if err != nil {
			roomInfoMsg.ErrorInfo = err
			ErrorHandling(msg, &roomInfoMsg)
			return fmt.Errorf("%s", err.Msg)
		}
		birdCards, err := repository.StartBirdCard(ctx, tx)
		if err != nil {
			roomInfoMsg.ErrorInfo = err
			ErrorHandling(msg, &roomInfoMsg)
			return fmt.Errorf("%s", err.Msg)
		}
		// 카드를 생성한다.
		cards := CreateInitCards(roomID, birdCards)
		err = repository.StartCreateCards(ctx, tx, cards)
		if err != nil {
			roomInfoMsg.ErrorInfo = err
			ErrorHandling(msg, &roomInfoMsg)
			return fmt.Errorf("%s", err.Msg)
		}

		return nil
	})

	if err != nil {
		return
	}
	preloadUsers, newErr = repository.StartFindAllRoomUsers(ctx, roomID)
	if newErr != nil {
		fmt.Println(newErr)
		return
	}
	// 메시지 생성
	roomInfoMsg = *CreateRoomInfoMSG(ctx, preloadUsers, 1, roomInfoMsg.ErrorInfo, 0)
	openCards, newErr := repository.StartUpdateCardState(ctx, roomID)
	if newErr != nil {
		fmt.Println(newErr)
	}
	roomInfoMsg.GameInfo.OpenCards = openCards

	message, err := CreateMessage(&roomInfoMsg)
	if err != nil {
		fmt.Println(err)
		return
	}
	msg.Message = message
	msg.SessionID = ""
	sendMessageToClients(roomID, msg)
}
