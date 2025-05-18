package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"main/features/frog/model/entity"
	_errors "main/features/frog/model/errors"
	"main/features/frog/model/request"
	"main/features/frog/repository"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func PlayTogetherEventWebsocket(msg *entity.WSMessage) *entity.ErrorInfo {
	ctx := context.Background()
	uID := msg.UserID

	//string to struct
	req := request.ReqWSPlayTogetherEvent{}
	err := json.Unmarshal([]byte(msg.Message), &req)
	if err != nil {
		return CreateErrorMessage(_errors.ErrCodeBadRequest, _errors.ErrUnmarshalFailed, "JSON 언마샬링 에러")
	}

	//비즈니스 로직
	roomInfoMsg := entity.RoomInfo{}
	preloadUsers := []entity.RoomUsers{}
	roomID, errInfo := repository.PlayTogetherFindOneRoomUsers(ctx, uID)
	if errInfo != nil {
		return errInfo
	}
	err = mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		// 방 정보를 업데이트 한다. (타이머, 인원 수)
		errInfo := repository.PlayTogetherFindOneAndUpdateRoom(ctx, tx, roomID)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}

		//유저 정보를 업데이트 한다.
		errInfo = repository.PlayTogetherFindOneAndUpdateUser(ctx, tx, uID, roomID)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}

		preloadUsers, errInfo = repository.PreloadFindGameInfo(ctx, tx, roomID)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}
		return nil
	})
	if err != nil {
		return errInfo
	}

	// 메시지 생성
	roomInfoMsg = *CreateRoomInfoMSG(ctx, preloadUsers, 0, roomInfoMsg.ErrorInfo)
	roomInfoMsg.GameInfo.AllReady = false

	if len(preloadUsers) == 2 {
		roomInfoMsg.GameInfo.IsFull = true
		roomInfoMsg.GameInfo.AllReady = true
	}

	message, err := CreateMessage(&roomInfoMsg)
	if err != nil {
		return CreateErrorMessage(_errors.ErrCodeInternal, err.Error(), _errors.ErrGameTerminated)
	}
	msg.Message = message
	sendMessageToClients(roomID, msg)
	return nil
}
