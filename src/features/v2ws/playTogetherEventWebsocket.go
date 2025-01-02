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

func PlayTogetherEventWebsocket(msg *entity.WSMessage) {
	ctx := context.Background()
	uID := msg.UserID
	decryptedMessage, err := utils.DecryptAES(msg.Message)
	if err != nil {
		errMsg := CreateErrorMessage(_errors.ErrCodeBadRequest, _errors.ErrCryptoFailed, "AES 복호화 에러")
		msg.Message = errMsg
		sendMessageToClient(msg.RoomID, msg)
	}
	//string to struct
	req := request.ReqWSPlayTogetherEvent{}
	err = json.Unmarshal([]byte(decryptedMessage), &req)
	if err != nil {
		errMsg := CreateErrorMessage(_errors.ErrCodeBadRequest, _errors.ErrUnmarshalFailed, "JSON 언마샬링 에러")
		msg.Message = errMsg
		sendMessageToClient(msg.RoomID, msg)
	}

	//비즈니스 로직
	roomInfoMsg := entity.RoomInfo{}
	preloadUsers := []entity.RoomUsers{}
	roomID, newErr := repository.PlayTogetherFindOneRoomUsers(ctx, uID)
	if newErr != nil {
		roomInfoMsg.ErrorInfo = newErr
		SendErrorMessage(msg, &roomInfoMsg)
		return
	}
	err = mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		// 방 정보를 업데이트 한다. (타이머, 인원 수)
		err := repository.PlayTogetherFindOneAndUpdateRoom(ctx, tx, roomID, uint(req.Count), uint(req.Timer))
		if err != nil {
			roomInfoMsg.ErrorInfo = err
			SendErrorMessage(msg, &roomInfoMsg)
			return fmt.Errorf("%s", err.Msg)
		}

		//유저 정보를 업데이트 한다.
		err = repository.PlayTogetherFindOneAndUpdateUser(ctx, tx, uID, roomID)
		if err != nil {
			roomInfoMsg.ErrorInfo = err
			SendErrorMessage(msg, &roomInfoMsg)
			return fmt.Errorf("%s", err.Msg)
		}
		// 미션 3개를 생성한다.
		err = repository.PlayTogetherCreateMissions(ctx, tx, roomID)
		if err != nil {
			roomInfoMsg.ErrorInfo = err
			SendErrorMessage(msg, &roomInfoMsg)
			return fmt.Errorf("%s", err.Msg)
		}
		preloadUsers, err = repository.PlayTogetherFindAllRoomUsers(ctx, tx, roomID)
		if err != nil {
			roomInfoMsg.ErrorInfo = err
			SendErrorMessage(msg, &roomInfoMsg)
			return fmt.Errorf("%s", err.Msg)
		}
		return nil
	})
	if err != nil {
		return
	}

	// 메시지 생성
	roomInfoMsg = *CreateRoomInfoMSG(ctx, preloadUsers, 1, roomInfoMsg.ErrorInfo, 0)
	roomInfoMsg.GameInfo.AllReady = false

	if len(preloadUsers) == req.Count {
		roomInfoMsg.GameInfo.IsFull = true
		roomInfoMsg.GameInfo.AllReady = true
	}

	message, err := CreateMessage(&roomInfoMsg)
	if err != nil {
		fmt.Println(err)
		return
	}
	msg.Message = message
	sendMessageToClients(roomID, msg)
}
