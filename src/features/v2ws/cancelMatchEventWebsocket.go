package v2ws

import (
	"context"
	"fmt"
	"log"
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

func CancelMatchEventWebsocket(msg *entity.WSMessage) {
	ctx := context.Background()
	uID := msg.UserID
	roomID := msg.RoomID

	//비즈니스 로직
	roomInfoMsg := entity.RoomInfo{}
	preloadUsers := []entity.RoomUsers{}
	err := mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		// 룸 유저 정보를 삭제한다.
		err := repository.CancelMatchDeleteOneRoomUser(ctx, tx, roomID, uID)
		if err != nil {
			return err
		}

		// 유저 정보를 업데이트 한다.
		err = repository.CancelMatchFindOneAndUpdateUser(ctx, tx, uID)
		if err != nil {
			return err
		}

		// 방 정보를 업데이트 한다. (방이 비어있으면 방을 삭제한다.)
		roomDTO, err := repository.CancelMatchFindOneAndUpdateRoom(ctx, tx, roomID)
		if err != nil {
			return err
		}
		//
		//방장이 나가면 다른 유저 중 한명을 방장으로 변경
		if roomDTO.CurrentCount != 0 && roomDTO.OwnerID == int(uID) {
			//룸 유저 정보를 가져온다.
			roomUserID, err := repository.CancelMatchFindOneRoomUser(ctx, tx, roomID)
			if err != nil {
				return err
			}
			//해당 유저ID를 방장으로 변경한다.
			err = repository.CancelMatchUpdateRoomOwner(ctx, tx, roomID, roomUserID)
			if err != nil {
				return err
			}
		}
		preloadUsers, err = repository.CancelMatchFindAllRoomUsers(ctx, tx, roomID)
		if err != nil {
			return err
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

	// 구조체를 JSON 문자열로 변환 (마샬링)
	message, err := CreateMessage(&roomInfoMsg)
	if err != nil {
		fmt.Println(err)
	}
	msg.Message = message
	msg.RoomID = roomID
	// 방 유저들에게 메시지 전달
	if sessionIDs, ok := entity.RoomSessions[msg.RoomID]; ok {
		// 에러 발생 시 이벤트 요청한 유저에게만 메시지를 전달한다.
		if roomInfoMsg.ErrorInfo != nil || err != nil {
			for _, sessionID := range sessionIDs {
				if client, exists := entity.WSClients[sessionID]; exists && client.UserID == msg.UserID {
					_ = client.Conn.WriteJSON(msg)

					// 클라이언트를 종료 및 정리
					client.Close()
					delete(entity.WSClients, sessionID)
					removeSessionFromRoom(client.RoomID, sessionID)
				}
			}
		} else {
			// 정상적인 경우 방의 모든 유저에게 메시지 전달
			for _, sessionID := range sessionIDs {
				if client, exists := entity.WSClients[sessionID]; exists {
					if client.UserID == msg.UserID {
						// 방에서 나간 유저 처리
						client.Close()
						delete(entity.WSClients, sessionID)
						removeSessionFromRoom(client.RoomID, sessionID)
					} else {
						// 나머지 유저에게 메시지 전달
						err := client.Conn.WriteJSON(msg)
						if err != nil {
							log.Printf("Error sending message to user %d: %v", client.UserID, err)
							client.Close()
							delete(entity.WSClients, sessionID)
							removeSessionFromRoom(client.RoomID, sessionID)
						}
					}
				}
			}
		}

		// 방이 비어 있으면 삭제
		if len(entity.RoomSessions[msg.RoomID]) == 0 {
			delete(entity.RoomSessions, msg.RoomID)
		}
	}

}
