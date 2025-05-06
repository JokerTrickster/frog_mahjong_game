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

func JoinPlayEventWebsocket(msg *entity.WSMessage) *entity.ErrorInfo {
	ctx := context.Background()
	uID := msg.UserID

	//string to struct
	req := request.ReqWSJoinPlayEvent{}
	err := json.Unmarshal([]byte(msg.Message), &req)
	if err != nil {
		return CreateErrorMessage(_errors.ErrCodeBadRequest, _errors.ErrUnmarshalFailed, "JSON 언마샬링 에러")
	}

	//비즈니스 로직
	messageInfoMsg := entity.MessageInfo{}
	preloadUsers := []entity.PreloadUsers{}
	var errInfo *entity.ErrorInfo
	roomID, errInfo := repository.JoinPlayFindOneRoomUsers(ctx, uID)
	if errInfo != nil {
		return errInfo
	}

	err = mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		//유저 정보를 업데이트 한다.
		errInfo = repository.JoinPlayFindOneAndUpdateUser(ctx, tx, uID, roomID)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}
		// 게임 유저 정보 생성
		SequenceUserDTO := CreateJoinPlayUserDTO(uID, roomID)
		errInfo = repository.JoinPlayInsertUserDTO(ctx, tx, SequenceUserDTO)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}

		// 게임 룸 셋팅 생성
		SequenceGameRoomSettingDTO := CreateJoinPlayGameRoomSettingDTO(roomID)
		errInfo = repository.JoinPlayInsertGameRoomSettingDTO(ctx, tx, SequenceGameRoomSettingDTO)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}

		preloadUsers, errInfo = repository.PreloadUsers(ctx, tx, roomID)
		if err != nil {
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
