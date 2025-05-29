package slime_war

import (
	"context"
	"fmt"
	"main/features/slime_war/model/entity"
	_errors "main/features/slime_war/model/errors"
	"main/features/slime_war/repository"
	"main/utils/db/mysql"
	"sync"

	"gorm.io/gorm"
)

var reconnectTimers sync.Map // 재접속 타이머를 관리하는 맵

// 비정상 연결로 인해 게임 강제 중단
func AbnormalSendErrorMessage(roomID, userID uint, sessionID string) {
	ctx := context.TODO()
	MessageInfoMsg := entity.MessageInfo{}
	preloadUsers := []entity.PreloadUsers{}
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
		users, err := repository.PreloadUsers(ctx, tx, roomID)
		if err != nil {
			return fmt.Errorf("%s", err.Msg)
		}
		preloadUsers = users

		// 에러 메시지 설정
		MessageInfoMsg.ErrorInfo = &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  "상대방이 게임 도중 나가서 강제 중단되었습니다.",
			Type: _errors.ErrGameTerminated,
		}
		return nil
	})

	// 트랜잭션 에러 처리
	if err != nil {
		MessageInfoMsg.ErrorInfo = &entity.ErrorInfo{
			Code: 500,
			Msg:  err.Error(),
			Type: _errors.ErrInternalServer,
		}
	}
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

	cleanupSession(roomID, sessionID, preloadUsers)
}
