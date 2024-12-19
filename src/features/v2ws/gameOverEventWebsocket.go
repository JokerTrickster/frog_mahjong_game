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

func GameOverEventWebsocket(msg *entity.WSMessage) {
	//유저 상태를 변경한다. (대기실로 이동)
	ctx := context.Background()
	uID := msg.UserID
	roomID := msg.RoomID

	GameOverEntity := entity.WSGameOverEntity{
		RoomID: roomID,
		UserID: uID,
	}
	// 비즈니스 로직
	roomInfoMsg := entity.RoomInfo{}
	preloadUsers := []entity.RoomUsers{}
	var err error
	err = mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {

		// 유저 상태 변경
		err = repository.GameOverUpdateRoomUsers(ctx, tx, &GameOverEntity)
		if err != nil {
			return err
		}
		preloadUsers, err = repository.GameOverFindAllRoomUsers(ctx, tx, roomID)
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

	// 유저 상태를 변경한다. (방에 참여)
	if sessionIDs, ok := entity.RoomSessions[msg.RoomID]; ok {
		// 에러 발생 시 이벤트 요청한 유저에게만 메시지를 전달한다.
		if roomInfoMsg.ErrorInfo != nil || err != nil {
			for _, sessionID := range sessionIDs {
				if client, exists := entity.WSClients[sessionID]; exists && client.UserID == msg.UserID {
					// 에러 메시지 전송
					err := client.Conn.WriteJSON(msg)
					if err != nil {
						fmt.Printf("Error sending message to user %d: %v\n", client.UserID, err)

						// 클라이언트를 종료 및 정리
						client.Close()
						delete(entity.WSClients, sessionID)
						removeSessionFromRoom(client.RoomID, sessionID)
					}
				}
			}
		} else {
			// 정상적인 경우 방의 모든 유저에게 메시지 전송
			for _, sessionID := range sessionIDs {
				if client, exists := entity.WSClients[sessionID]; exists {
					err := client.Conn.WriteJSON(msg)
					if err != nil {
						fmt.Printf("Error sending message to user %d: %v\n", client.UserID, err)

						// 클라이언트를 종료 및 정리
						client.Close()
						delete(entity.WSClients, sessionID)
						removeSessionFromRoom(client.RoomID, sessionID)
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
