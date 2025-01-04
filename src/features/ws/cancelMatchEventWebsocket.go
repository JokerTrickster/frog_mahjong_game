package ws

import (
	"context"
	"fmt"
	"main/features/ws/model/entity"
	"main/features/ws/repository"
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

func CancelMatchEventWebsocket(msg *entity.WSMessage) {
	ctx := context.Background()
	uID := msg.UserID
	roomID := msg.RoomID

	//비즈니스 로직
	roomInfoMsg := entity.RoomInfo{}
	preloadUsers := []entity.RoomUsers{}
	err := mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		// 룸 유저 정보를 삭제한다.
		newErr := repository.CancelMatchDeleteOneRoomUser(ctx, tx, roomID, uID)
		if newErr != nil {
			SendErrorMessage(msg, newErr)
			return fmt.Errorf("%s", newErr.Msg)
		}

		// 유저 정보를 업데이트 한다.
		newErr = repository.CancelMatchFindOneAndUpdateUser(ctx, tx, uID)
		if newErr != nil {
			SendErrorMessage(msg, newErr)
			return fmt.Errorf("%s", newErr.Msg)
		}

		// 방 정보를 업데이트 한다. (방이 비어있으면 방을 삭제한다.)
		roomDTO, newErr := repository.CancelMatchFindOneAndUpdateRoom(ctx, tx, roomID)
		if newErr != nil {
			SendErrorMessage(msg, newErr)
			return fmt.Errorf("%s", newErr.Msg)
		}
		//
		//방장이 나가면 다른 유저 중 한명을 방장으로 변경
		if roomDTO.CurrentCount != 0 && roomDTO.OwnerID == int(uID) {
			//룸 유저 정보를 가져온다.
			roomUserID, newErr := repository.CancelMatchFindOneRoomUser(ctx, tx, roomID)
			if newErr != nil {
				SendErrorMessage(msg, newErr)
				return fmt.Errorf("%s", newErr.Msg)
			}
			//해당 유저ID를 방장으로 변경한다.
			newErr = repository.CancelMatchUpdateRoomOwner(ctx, tx, roomID, roomUserID)
			if newErr != nil {
				SendErrorMessage(msg, newErr)
				return fmt.Errorf("%s", newErr.Msg)
			}
		}
		preloadUsers, newErr = repository.PreloadFindGameInfo(ctx, tx, roomID)
		if newErr != nil {
			SendErrorMessage(msg, newErr)
			return fmt.Errorf("%s", newErr.Msg)
		}
		return nil
	})
	if err != nil {
		return
	}

	// 메시지 생성
	roomInfoMsg = *CreateRoomInfoMSG(ctx, preloadUsers, 1, roomInfoMsg.ErrorInfo)
	roomInfoMsg.GameInfo.AllReady = false

	// 구조체를 JSON 문자열로 변환 (마샬링)
	message, err := CreateMessage(&roomInfoMsg)
	if err != nil {
		fmt.Println(err)
	}
	msg.Message = message
	msg.RoomID = roomID
	sendMessageToClients(roomID, msg)

}
