package ws

import (
	"context"
	"fmt"
	"main/features/frog/model/entity"
	_errors "main/features/frog/model/errors"
	"main/features/frog/repository"
	"main/utils/db/mysql"
	"sync"

	"gorm.io/gorm"
)

var reconnectTimers sync.Map

func AbnormalSendErrorMessage(roomID, userID uint, sessionID string) {
	// 비정상적인 에러 발생했으므로 비정상적 에러 처리하는 로직 실행

	//business logic

	roomInfoMsg := entity.RoomInfo{}
	preloadUsers := []entity.RoomUsers{}

	ctx := context.TODO()
	err := mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		//모든 유저 게임 종료 처리하고 대기 상태로 변경한다.
		abnormalEntity := entity.WSAbnormalEntity{
			RoomID:         roomID,
			AbnormalUserID: userID,
		}
		// 유저 상태 변경
		errInfo := repository.AbnormalUpdateUsers(ctx, tx, &abnormalEntity)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}
		preloadUsers, errInfo = repository.PreloadFindGameInfo(ctx, tx, roomID)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}
		// 에러 메시지에 상대방이 게임 도중 나가서 강제 종료됐다는 에러 메시지 표시한다.
		roomInfoMsg.ErrorInfo = &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  "상대방이 게임 도중 나가서 강제 중단되었습니다.",
			Type: _errors.ErrGameTerminated,
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
	msg := entity.WSMessage{
		RoomID:  roomID,
		UserID:  userID,
		Message: message,
	}
	sendMessageToClients(roomID, &msg)
	// `entity.WSClients`에서 sessionID가 존재하는지 확인
	if _, exists := entity.WSClients[sessionID]; !exists {
		fmt.Printf("Session %s does not exist in WSClients. Skipping cleanup.\n", sessionID)
		return
	}

	// `entity.RoomSessions`에서 roomID가 존재하는지 확인
	if _, exists := entity.RoomSessions[roomID]; !exists {
		fmt.Printf("Room %d does not exist in RoomSessions. Skipping cleanup.\n", roomID)
		return
	}

	fmt.Printf("Session %s in room %d failed to reconnect. Cleaning up.\n", sessionID, roomID)
	cleanupSession(roomID, sessionID, preloadUsers)
}
