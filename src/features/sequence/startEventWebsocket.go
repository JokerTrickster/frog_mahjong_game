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

func StartEventWebsocket(msg *entity.WSMessage) *entity.ErrorInfo {
	//유저 상태를 변경한다. (대기실로 이동)
	ctx := context.Background()
	uID := msg.UserID
	roomID := msg.RoomID

	// 비즈니스 로직
	//해당 방이 대기상태인지 체크한다.
	preloadUsers := []entity.PreloadUsers{}
	messageMsg := entity.MessageInfo{}
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
		// 방 카드 정보 생성
		SequenceCards := CreateSequenceCards(roomID)
		errInfo = repository.StartCreateSequenceCards(ctx, tx, SequenceCards)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}

		// 슬라임 유저 정보 변경 (colorType, turn 변경)
		errInfo = repository.StartUpdateSequenceUser(ctx, tx, roomID)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}

		// 방 유저 정보를 가져온다.
		roomUsers, errInfo := repository.StartFindRoomUsers(ctx, tx, roomID)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}

		// 유저에게 랜덤으로 5개 카드를 부여한다.
		errInfo = repository.StartCreateSequenceUserCards(ctx, tx, roomUsers)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}

		// 방 맵 정보 생성
		SequenceMaps := CreateSequenceMaps(roomID)
		errInfo = repository.StartCreateSequenceMaps(ctx, tx, SequenceMaps)
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
