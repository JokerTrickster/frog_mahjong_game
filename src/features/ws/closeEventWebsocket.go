package ws

import (
	"context"
	"fmt"
	"main/features/ws/model/entity"
	"main/features/ws/model/request"
	"main/features/ws/repository"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func CloseEventWebsocket(msg *entity.WSMessage) *entity.ErrorInfo {
	//유저 상태를 변경한다. (대기실로 이동)
	ctx := context.Background()
	uID := msg.UserID
	req := request.ReqWSClose{
		RoomID: uint(msg.RoomID),
	}

	//비즈니스 로직
	roomInfoMsg := entity.RoomInfo{}
	preloadUsers := []entity.RoomUsers{}
	var errInfo *entity.ErrorInfo
	err := mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		// RoomsID에 해당하는 userID를 삭제한다.
		errInfo = repository.CloseFindOneAndDeleteRoomUser(ctx, tx, uID, req.RoomID)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}
		// Rooms 현재 인원수를 -1한다.
		roomDTO, errInfo := repository.CloseFindOneAndUpdateRoom(ctx, tx, req.RoomID)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}
		// user에 rooms_id를 1로 바꾸고 state를 wait으로 변경한다.
		errInfo = repository.CloseFindOneAndUpdateUser(ctx, tx, uID)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}
		if roomDTO.CurrentCount == 0 {
			// 방 삭제
			errInfo = repository.CloseFindOneAndDeleteRoom(ctx, tx, req.RoomID)
			if errInfo != nil {
				return fmt.Errorf("%s", errInfo.Msg)
			}

		} else if roomDTO.CurrentCount == 1 {
			// 인원이 1명이면 남아 있는 유저를 방장으로 변경
			//방장이 나가면 다른 유저 중 한명을 방장으로 변경
			//룸에 남아있는 유저 정보를 가져온다.
			roomUserDTO, errInfo := repository.CloseFindOneRoomUser(ctx, tx, req.RoomID)
			if errInfo != nil {
				return fmt.Errorf("%s", errInfo.Msg)
			}
			userDTO, errInfo := repository.CloseFindOneUser(ctx, tx, uint(roomUserDTO.UserID))
			if errInfo != nil {
				return fmt.Errorf("%s", errInfo.Msg)
			}

			// 해당 유저를 방장으로 업데이트 한다.
			//방장으로 변경하기 위해 업데이트해야 될 부분들
			// rooms -> owner 변경
			errInfo = repository.CloseChangeRoomOnwer(ctx, tx, req.RoomID, userDTO.ID)
			if errInfo != nil {
				return fmt.Errorf("%s", errInfo.Msg)
			}
		}
		preloadUsers, errInfo = repository.PreloadFindGameInfo(ctx, tx, req.RoomID)
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
	// 구조체를 JSON 문자열로 변환 (마샬링)
	message, err := CreateMessage(&roomInfoMsg)
	if err != nil {
		fmt.Println(err)
	}
	msg.Message = message
	sendMessageToClients(req.RoomID, msg)

	return nil
}
