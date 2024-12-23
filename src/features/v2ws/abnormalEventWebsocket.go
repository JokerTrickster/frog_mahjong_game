package v2ws

import (
	"context"
	"fmt"
	"main/features/v2ws/model/entity"
	_errors "main/features/v2ws/model/errors"
	"main/features/v2ws/repository"
	"main/utils/db/mysql"
	"sync"

	"gorm.io/gorm"
)

var reconnectTimers sync.Map // 재접속 타이머를 관리하는 맵

// 비정상적인 에러를 처리하는 함수
func AbnormalErrorHandling(roomID, userID uint, sessionID string) {
	ctx := context.TODO()
	roomInfoMsg := entity.RoomInfo{}
	preloadUsers := []entity.RoomUsers{}
	// 비정상적인 유저 상태 처리
	err := mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		abnormalEntity := entity.WSAbnormalEntity{
			RoomID:         roomID,
			UserID:         userID,
			AbnormalUserID: getUserIDFromSessionID(sessionID),
		}

		// 유저 상태 변경
		if err := repository.AbnormalUpdateRoomUsers(ctx, tx, &abnormalEntity); err != nil {
			return fmt.Errorf("%s", err.Msg)
		}

		// 방 유저 정보 조회
		users, err := repository.AbnormalFindAllRoomUsers(ctx, tx, roomID)
		if err != nil {
			return fmt.Errorf("%s", err.Msg)
		}
		preloadUsers = users

		// 에러 메시지 설정
		roomInfoMsg.ErrorInfo = &entity.ErrorInfo{
			Code: 500,
			Msg:  "상대방이 게임 도중 나가서 강제 종료됐습니다.",
			Type: _errors.ErrAbnormalExit,
		}
		return nil
	})

	// 트랜잭션 에러 처리
	if err != nil {
		roomInfoMsg.ErrorInfo = &entity.ErrorInfo{
			Code: 500,
			Msg:  err.Error(),
			Type: _errors.ErrInternalServer,
		}
	}

	// 클라이언트에 메시지 전송
	roomInfoMsg = *CreateRoomInfoMSG(ctx, preloadUsers, 1, roomInfoMsg.ErrorInfo, 0)
	roomInfoMsg.GameInfo.AllReady = false

	message, err := CreateMessage(&roomInfoMsg)
	if err != nil {
		fmt.Println(err)
		return
	}
	msg := entity.WSMessage{
		RoomID:  roomID,
		UserID:  userID,
		Message: message,
	}
	sendMessageToClients(roomID, &msg)

	// 재접속 대기 시작
	waitForReconnection(roomID, sessionID, preloadUsers)
}
