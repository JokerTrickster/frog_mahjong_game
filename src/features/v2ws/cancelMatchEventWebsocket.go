package v2ws

import (
	"context"
	"fmt"
	"main/features/v2ws/model/entity"
	_errors "main/features/v2ws/model/errors"
	"main/features/v2ws/repository"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

/*
	방 매칭할 떄 생성되는 데이터를 원상 복구하거나 삭제하고 연결을 끊어야 된다.
	매칭할 때 생성되는 데이터들
	1. 유저가 0명이면 방을 생성
	2. 유저가 1명이면 방 참여
	3. 방 유저수 증가
	4. room user 정보 생성
	5. 유저 정보 업데이트

*/

func CancelMatchEventWebsocket(msg *entity.WSMessage) *entity.ErrorInfo {
	ctx := context.Background()
	uID := msg.UserID
	rID := msg.RoomID

	//비즈니스 로직
	roomInfoMsg := entity.RoomInfo{}
	preloadUsers := []entity.RoomUsers{}
	var errInfo *entity.ErrorInfo
	err := mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		// 룸 유저 정보를 삭제한다.
		errInfo = repository.CancelMatchDeleteOneRoomUser(ctx, tx, rID, uID)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}

		// 유저 정보를 업데이트 한다.
		errInfo = repository.CancelMatchFindOneAndUpdateUser(ctx, tx, uID)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}

		// 방 정보를 업데이트 한다. (방이 비어있으면 방을 삭제한다.)
		roomDTO, errInfo := repository.CancelMatchFindOneAndUpdateRoom(ctx, tx, rID)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}
		//
		//방장이 나가면 다른 유저 중 한명을 방장으로 변경
		if roomDTO.CurrentCount != 0 && roomDTO.OwnerID == int(uID) {
			//룸 유저 정보를 가져온다.
			roomUserID, errInfo := repository.CancelMatchFindOneRoomUser(ctx, tx, rID)
			if errInfo != nil {
				return fmt.Errorf("%s", errInfo.Msg)
			}
			//해당 유저ID를 방장으로 변경한다.
			errInfo = repository.CancelMatchUpdateRoomOwner(ctx, tx, rID, roomUserID)
			if errInfo != nil {
				return fmt.Errorf("%s", errInfo.Msg)
			}
		}
		preloadUsers, errInfo = repository.CancelMatchFindAllRoomUsers(ctx, tx, rID)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}
		return nil
	})
	if err != nil {
		return errInfo
	}

	// 메시지 생성
	roomInfoMsg = *CreateRoomInfoMSG(ctx, preloadUsers, 1, roomInfoMsg.ErrorInfo, 0)
	roomInfoMsg.GameInfo.AllReady = false

	message, err := CreateMessage(&roomInfoMsg)
	if err != nil {
		fmt.Println(err)
		return CreateErrorMessage(_errors.ErrCodeInternal, err.Error(), _errors.ErrGameTerminated)

	}
	msg.Message = message
	sendMessageToClients(rID, msg)

	// 정상적으로 연결을 끊는다.
	if sessionIDs, ok := entity.RoomSessions[rID]; ok {
		for _, sessionID := range sessionIDs {
			if client, exists := entity.WSClients[sessionID]; exists && client.UserID == uID {
				closeAndRemoveClient(client, sessionID, rID)
			}
		}
	}
	return nil
}
