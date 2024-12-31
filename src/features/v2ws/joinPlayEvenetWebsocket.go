package v2ws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"main/features/v2ws/model/entity"
	_errors "main/features/v2ws/model/errors"
	"main/features/v2ws/model/request"
	"main/features/v2ws/repository"
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func JoinPlayEventWebsocket(msg *entity.WSMessage) {
	ctx := context.Background()
	uID := msg.UserID
	decryptedMessage, err := utils.DecryptAES(msg.Message)
	if err != nil {
		log.Fatalf("AES 복호화 에러: %s", err)
	}
	//string to struct
	req := request.ReqWSJoinPlayEvent{}
	err = json.Unmarshal([]byte(decryptedMessage), &req)
	if err != nil {
		log.Fatalf("JSON 언마샬링 에러: %s", err)
	}

	//비즈니스 로직
	roomInfoMsg := entity.RoomInfo{}
	preloadUsers := []entity.RoomUsers{}
	roomID, newErr := repository.JoinPlayFindOneRoomUsers(ctx, uID)
	if newErr != nil {
		roomInfoMsg.ErrorInfo = newErr
		ErrorHandling(msg, &roomInfoMsg)
		return
	}
	roomDTO, newErr := repository.JoinPlayFindOneRoom(ctx, roomID)
	if newErr != nil {
		roomInfoMsg.ErrorInfo = newErr
		ErrorHandling(msg, &roomInfoMsg)
		return
	}
	err = mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		//유저 정보를 업데이트 한다.
		err := repository.JoinPlayFindOneAndUpdateUser(ctx, tx, uID, roomID)
		if err != nil {
			roomInfoMsg.ErrorInfo = err
			ErrorHandling(msg, &roomInfoMsg)
			return fmt.Errorf("%s", err.Msg)
		}

		preloadUsers, err = repository.JoinPlayFindAllRoomUsers(ctx, tx, roomID)
		if err != nil {
			roomInfoMsg.ErrorInfo = err
			ErrorHandling(msg, &roomInfoMsg)
			return fmt.Errorf("%s", err.Msg)
		}
		return nil
	})
	if err != nil {
		roomInfoMsg.ErrorInfo = &entity.ErrorInfo{
			Code: 500,
			Msg:  err.Error(),
			Type: _errors.ErrInternalServer,
		}
		if roomInfoMsg.ErrorInfo.Msg == "방이 꽉 찼습니다." {
			roomInfoMsg.ErrorInfo.Type = _errors.ErrRoomFull
		} else if roomInfoMsg.ErrorInfo.Msg == "비밀번호가 일치하지 않습니다." {
			roomInfoMsg.ErrorInfo.Type = _errors.ErrWrongPassword
		} else if roomInfoMsg.ErrorInfo.Msg == "게임 중인 방입니다." {
			roomInfoMsg.ErrorInfo.Type = _errors.ErrGameInProgress
		}
	}

	// 메시지 생성
	roomInfoMsg = *CreateRoomInfoMSG(ctx, preloadUsers, 1, roomInfoMsg.ErrorInfo, 0)
	roomInfoMsg.GameInfo.AllReady = false

	if len(preloadUsers) == roomDTO.MinCount {
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
