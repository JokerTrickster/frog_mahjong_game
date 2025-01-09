package v2ws

import (
	"context"
	"encoding/json"
	"fmt"
	"main/features/v2ws/model/entity"
	_errors "main/features/v2ws/model/errors"
	"main/features/v2ws/model/request"
	"main/features/v2ws/repository"
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func ItemChangeEventWebsocket(msg *entity.WSMessage) *entity.ErrorInfo {
	//유저 상태를 변경한다. (대기실로 이동)
	ctx := context.Background()
	uID := msg.UserID
	roomID := msg.RoomID
	// 복호화 후 JSON 언마샬링
	decryptedMessage, err := utils.DecryptAES(msg.Message)
	if err != nil {
		return CreateErrorMessage(_errors.ErrCodeBadRequest, _errors.ErrCryptoFailed, "AES 복호화 에러")
	}
	//string to struct
	req := request.ReqWSItemChange{}
	err = json.Unmarshal([]byte(decryptedMessage), &req)
	if err != nil {
		return CreateErrorMessage(_errors.ErrCodeBadRequest, _errors.ErrUnmarshalFailed, "JSON 언마샬링 에러")
	}

	ItemChangeEntity := entity.WSItemChangeEntity{
		RoomID: roomID,
		UserID: uID,
		ItemID: req.ItemID,
	}
	// 비즈니스 로직
	roomInfoMsg := entity.RoomInfo{}
	preloadUsers := []entity.RoomUsers{}
	var errInfo *entity.ErrorInfo
	// 아이템 사용 가능한지 체크
	errInfo = repository.ItemChangeCheck(ctx, ItemChangeEntity)
	if errInfo != nil {
		return errInfo
	}

	err = mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		// 아이템 사용 (오픈된 카드를 교체한다.)
		errInfo = repository.ItemChange(ctx, tx, ItemChangeEntity)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}
		// 아이템 사용 횟수를 -1 감소한다.
		errInfo = repository.ItemChangeConsumeUserItems(ctx, tx, ItemChangeEntity)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}
		return nil
	})
	if err != nil {
		return errInfo
	}
	// 현재 참여하고 있는 유저에 대한 정보를 가져와서 메시지 전달한다.
	preloadUsers, errInfo = repository.ItemChangeFindAllRoomUsers(ctx, roomID)
	if errInfo != nil {
		return errInfo
	}
	roomInfoMsg = *CreateRoomInfoMSG(ctx, preloadUsers, 1, roomInfoMsg.ErrorInfo, 0)

	message, err := CreateMessage(&roomInfoMsg)
	if err != nil {
		return CreateErrorMessage(_errors.ErrCodeBadRequest, _errors.ErrMarshalFailed, "JSON 마샬링 에러")
	}
	msg.Message = message
	sendMessageToClients(roomID, msg)
	return nil
}
