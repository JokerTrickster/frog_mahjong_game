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

func MatchEventWebsocket(msg *entity.WSMessage) *entity.ErrorInfo {
	ctx := context.Background()
	uID := msg.UserID

	//string to struct
	req := request.ReqWSMatchEvent{}
	err := json.Unmarshal([]byte(msg.Message), &req)
	if err != nil {
		return CreateErrorMessage(_errors.ErrCodeBadRequest, _errors.ErrUnmarshalFailed, "JSON 언마샬링 에러")
	}

	//비즈니스 로직
	roomInfoMsg := entity.RoomInfo{}
	preloadUsers := []entity.RoomUsers{}
	roomID, errInfo := repository.MatchFindOneRoomUsers(ctx, uID)
	if errInfo != nil {
		return CreateErrorMessage(_errors.ErrCodeBadRequest, _errors.ErrInternalServer, err.Error())
	}
	err = mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		errInfo := &entity.ErrorInfo{}
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
	roomInfoMsg = *CreateRoomInfoMSG(ctx, preloadUsers, 1, roomInfoMsg.ErrorInfo)
	roomInfoMsg.GameInfo.AllReady = false

	if len(preloadUsers) == 2 {
		roomInfoMsg.GameInfo.IsFull = true
		roomInfoMsg.GameInfo.AllReady = true
	}
	// 구조체를 JSON 문자열로 변환 (마샬링)
	message, err := CreateMessage(&roomInfoMsg)
	if err != nil {
		return CreateErrorMessage(_errors.ErrCodeInternal, err.Error(), _errors.ErrGameTerminated)
	}
	msg.Message = message
	msg.RoomID = roomID

	sendMessageToClients(roomID, msg)
	return nil
}
