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

func MatchEventWebsocket(msg *entity.WSMessage) {
	ctx := context.Background()
	uID := msg.UserID
	decryptedMessage, err := utils.DecryptAES(msg.Message)
	if err != nil {
		errMsg := CreateErrorMessage(_errors.ErrCodeBadRequest, _errors.ErrCryptoFailed, "AES 복호화 에러")
		msg.Message = errMsg
		sendMessageToClient(msg.RoomID, msg)
	}
	//string to struct
	req := request.ReqWSMatchEvent{}
	err = json.Unmarshal([]byte(decryptedMessage), &req)
	if err != nil {
		errMsg := CreateErrorMessage(_errors.ErrCodeBadRequest, _errors.ErrUnmarshalFailed, "JSON 언마샬링 에러")
		msg.Message = errMsg
		sendMessageToClient(msg.RoomID, msg)
	}

	//비즈니스 로직
	roomInfoMsg := entity.RoomInfo{}
	preloadUsers := []entity.RoomUsers{}
	roomID, newErr := repository.MatchFindOneRoomUsers(ctx, uID)
	if newErr != nil {
		roomInfoMsg.ErrorInfo = newErr
		SendErrorMessage(msg, &roomInfoMsg)
		return
	}
	err = mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		//유저 정보를 업데이트 한다.
		err := repository.MatchFindOneAndUpdateUser(ctx, tx, uID, roomID)
		if err != nil {
			roomInfoMsg.ErrorInfo = err
			SendErrorMessage(msg, &roomInfoMsg)
			return fmt.Errorf("%s", err.Msg)
		}
		//해당 방에 미션이 존재하는지 체크한다.
		roomMission, newErr := repository.MatchFindOneRoomMission(ctx, tx, roomID)
		if newErr != nil {
			roomInfoMsg.ErrorInfo = newErr
			SendErrorMessage(msg, &roomInfoMsg)
			return fmt.Errorf("%s", newErr.Msg)
		}
		if len(roomMission) == 0 {
			// 미션을 랜덤으로 3개 생성한다.
			err = repository.MatchCreateMissions(ctx, tx, roomID)
			if err != nil {
				roomInfoMsg.ErrorInfo = err
				SendErrorMessage(msg, &roomInfoMsg)
				return fmt.Errorf("%s", err.Msg)
			}
		}
		preloadUsers, err = repository.MatchFindAllRoomUsers(ctx, tx, roomID)
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
