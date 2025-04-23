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

func MatchEventWebsocket(msg *entity.WSMessage) *entity.ErrorInfo {
	ctx := context.Background()
	uID := msg.UserID
	// decryptedMessage, err := utils.DecryptAES(msg.Message)
	// if err != nil {
	// 	return CreateErrorMessage(_errors.ErrCodeBadRequest, _errors.ErrCryptoFailed, "AES 복호화 에러")
	// }
	//string to struct

	//비즈니스 로직
	MessageInfoMsg := entity.MessageInfo{}
	preloadUsers := []entity.PreloadUsers{}
	var errInfo *entity.ErrorInfo
	roomID, errInfo := repository.MatchFindOneRoomUsers(ctx, uID)
	if errInfo != nil {
		return errInfo
	}
	err := mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		//유저 정보를 업데이트 한다.
		errInfo = repository.MatchFindOneAndUpdateUser(ctx, tx, uID, roomID)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}
		// 게임 유저 정보 생성
		slimeWarUserDTO := CreateMatchUserDTO(uID, roomID)
		errInfo = repository.MatchInsertUserDTO(ctx, tx, slimeWarUserDTO)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}

		// 게임 룸 셋팅 생성
		slimeWarGameRoomSettingDTO := CreateMatchGameRoomSettingDTO(roomID)
		errInfo = repository.MatchInsertGameRoomSettingDTO(ctx, tx, slimeWarGameRoomSettingDTO)
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
	MessageInfoMsg = *CreateMessageInfoMSG(ctx, preloadUsers, 1, MessageInfoMsg.ErrorInfo, 0)


	message, err := CreateMessage(&MessageInfoMsg)
	if err != nil {
		return CreateErrorMessage(_errors.ErrCodeBadRequest, _errors.ErrMarshalFailed, "메시지 생성 에러")
	}
	msg.Message = message
	sendMessageToClients(roomID, msg)
	return nil
}
