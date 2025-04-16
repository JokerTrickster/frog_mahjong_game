package slime_war

import (
	"context"
	"fmt"
	"main/features/slime_war/model/entity"
	_errors "main/features/slime_war/model/errors"
	"main/features/slime_war/repository"
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

		//TODO 20라운드 이미지를 선택해서 각 라운드마다 이미지를 만든다.
		imagesDTO, errInfo := repository.FindImages(ctx, tx)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}

		// 각 라운드마다 이미지를 지정한다.
		roundImagesDTO := CreateRoundImages(roomID, imagesDTO)
		errInfo = repository.CreateRoundImages(ctx, tx, roundImagesDTO)
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

	if len(preloadUsers) == 2 {
		messageMsg.GameInfo.IsFull = true
	}

	message, err := CreateMessage(&messageMsg)
	if err != nil {
		return CreateErrorMessage(_errors.ErrCodeBadRequest, _errors.ErrMarshalFailed, "메시지 생성 에러")
	}
	msg.Message = message
	sendMessageToClients(roomID, msg)
	return nil
}
