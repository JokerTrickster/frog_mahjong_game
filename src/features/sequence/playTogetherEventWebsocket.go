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

func PlayTogetherEventWebsocket(msg *entity.WSMessage) *entity.ErrorInfo {
	ctx := context.Background()
	uID := msg.UserID

	//비즈니스 로직
	messageInfoMsg := entity.MessageInfo{}
	preloadUsers := []entity.PreloadUsers{}
	var errInfo *entity.ErrorInfo
	roomID, errInfo := repository.PlayTogetherFindOneRoomUsers(ctx, uID)
	if errInfo != nil {
		return errInfo
	}
	err := mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		// 방 정보를 업데이트 한다. (타이머, 인원 수)
		errInfo = repository.PlayTogetherFindOneAndUpdateRoom(ctx, tx, roomID)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}

		//유저 정보를 업데이트 한다.
		errInfo = repository.PlayTogetherFindOneAndUpdateUser(ctx, tx, uID, roomID)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}

		// 게임 유저 정보 생성
		SequenceUserDTO := CreatePlayTogetherUserDTO(uID, roomID)
		errInfo = repository.PlayTogetherInsertUserDTO(ctx, tx, SequenceUserDTO)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}

		// 게임 룸 셋팅 생성
		SequenceGameRoomSettingDTO := CreatePlayTogetherGameRoomSettingDTO(roomID)
		errInfo = repository.PlayTogetherInsertGameRoomSettingDTO(ctx, tx, SequenceGameRoomSettingDTO)
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
	messageInfoMsg = *CreateMessageInfoMSG(ctx, preloadUsers, 1, messageInfoMsg.ErrorInfo, 0)
	if len(preloadUsers) == 1 {
		messageInfoMsg.SequenceGameInfo.IsFull = false
		messageInfoMsg.SequenceGameInfo.AllReady = false
	}

	message, err := CreateMessage(&messageInfoMsg)
	if err != nil {
		return CreateErrorMessage(_errors.ErrCodeBadRequest, _errors.ErrMarshalFailed, "메시지 생성 에러")
	}
	msg.Message = message
	sendMessageToClients(roomID, msg)

	return nil
}
