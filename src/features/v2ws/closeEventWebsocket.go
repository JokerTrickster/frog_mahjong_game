package v2ws

import (
	"context"
	"fmt"
	"main/features/v2ws/model/entity"
	_errors "main/features/v2ws/model/errors"
	"main/features/v2ws/model/request"
	"main/features/v2ws/repository"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func CloseEventWebsocket(msg *entity.WSMessage) {
	//유저 상태를 변경한다. (대기실로 이동)
	ctx := context.Background()
	uID := msg.UserID
	roomID := msg.RoomID
	req := request.ReqWSClose{
		RoomID: uint(msg.RoomID),
	}

	//비즈니스 로직
	roomInfoMsg := entity.RoomInfo{}
	preloadUsers := []entity.RoomUsers{}
	err := mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		// RoomsID에 해당하는 userID를 삭제한다.
		err := repository.CloseFindOneAndDeleteRoomUser(ctx, tx, uID, req.RoomID)
		if err != nil {
			roomInfoMsg.ErrorInfo = err
			ErrorHandling(msg, &roomInfoMsg)
			return fmt.Errorf("%s", err.Msg)
		}
		// Rooms 현재 인원수를 -1한다.
		roomDTO, err := repository.CloseFindOneAndUpdateRoom(ctx, tx, req.RoomID)
		if err != nil {
			roomInfoMsg.ErrorInfo = err
			ErrorHandling(msg, &roomInfoMsg)
			return fmt.Errorf("%s", err.Msg)
		}
		// user에 rooms_id를 1로 바꾸고 state를 wait으로 변경한다.
		err = repository.CloseFindOneAndUpdateUser(ctx, tx, uID)
		if err != nil {
			roomInfoMsg.ErrorInfo = err
			ErrorHandling(msg, &roomInfoMsg)
			return fmt.Errorf("%s", err.Msg)
		}

		if roomDTO.CurrentCount == 0 {
			// 방 삭제
			err = repository.CloseFindOneAndDeleteRoom(ctx, tx, req.RoomID)
			if err != nil {
				roomInfoMsg.ErrorInfo = err
				ErrorHandling(msg, &roomInfoMsg)
				return fmt.Errorf("%s", err.Msg)
			}

		} else if roomDTO.CurrentCount == 1 {
			// 인원이 1명이면 남아 있는 유저를 방장으로 변경
			//방장이 나가면 다른 유저 중 한명을 방장으로 변경
			//룸에 남아있는 유저 정보를 가져온다.
			roomUserDTO, err := repository.CloseFindOneRoomUser(ctx, tx, req.RoomID)
			if err != nil {
				roomInfoMsg.ErrorInfo = err
				ErrorHandling(msg, &roomInfoMsg)
				return fmt.Errorf("%s", err.Msg)
			}
			userDTO, err := repository.CloseFindOneUser(ctx, tx, uint(roomUserDTO.UserID))
			if err != nil {
				roomInfoMsg.ErrorInfo = err
				ErrorHandling(msg, &roomInfoMsg)
				return fmt.Errorf("%s", err.Msg)
			}

			// 해당 유저를 방장으로 업데이트 한다.

			//방장으로 변경하기 위해 업데이트해야 될 부분들
			// rooms -> owner 변경
			err = repository.CloseChangeRoomOnwer(ctx, tx, req.RoomID, userDTO.ID)
			if err != nil {
				roomInfoMsg.ErrorInfo = err
				ErrorHandling(msg, &roomInfoMsg)
				return fmt.Errorf("%s", err.Msg)
			}
		}
		preloadUsers, err = repository.CloseFindAllRoomUsers(ctx, tx, req.RoomID)
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
	}

	// 메시지 생성
	roomInfoMsg = *CreateRoomInfoMSG(ctx, preloadUsers, 1, roomInfoMsg.ErrorInfo, 0)

	message, err := CreateMessage(&roomInfoMsg)
	if err != nil {
		fmt.Println(err)
		return
	}
	msg.Message = message
	sendMessageToClients(roomID, msg)
}
