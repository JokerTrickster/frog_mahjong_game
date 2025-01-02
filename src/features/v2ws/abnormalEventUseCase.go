package v2ws

import (
	"context"
	"fmt"
	"main/features/v2ws/model/entity"
	_errors "main/features/v2ws/model/errors"
	"main/features/v2ws/repository"
	"main/utils/db/mysql"
	"time"

	"gorm.io/gorm"
)

// 재접속 대기 타이머 (30초)
func waitForReconnection(roomID uint, sessionID string, preloadUsers []entity.RoomUsers) {
	fmt.Printf("Waiting for session %s to reconnect in room %d...\n", sessionID, roomID)

	// 타이머 설정
	timer := time.AfterFunc(30*time.Second, func() {
		fmt.Printf("Session %s in room %d failed to reconnect. Cleaning up.\n", sessionID, roomID)
		cleanupSession(roomID, sessionID, preloadUsers)
	})

	// 타이머 저장
	reconnectTimers.Store(sessionID, timer)
}

// 세션 정리 (재접속 실패 시 호출)
func cleanupSession(roomID uint, sessionID string, preloadUsers []entity.RoomUsers) {
	ctx := context.TODO()

	msg := &entity.WSMessage{
		RoomID: roomID,
		Event:  "ERROR",
	}
	errMsg := CreateErrorMessage(_errors.ErrCodeInternal, _errors.ErrAbnormalExit, "상대방이 연결이 끊겼습니다. 강제로 게임을 종료합니다.")
	msg.Message = errMsg
	sendMessageToClients(roomID, msg)

	fmt.Printf("Cleaning up session %s for room %d...\n", sessionID, roomID)

	err := mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		abnormalEntity := entity.WSAbnormalEntity{
			RoomID:         roomID,
			AbnormalUserID: getUserIDFromSessionID(sessionID),
		}

		// 카드 삭제
		if err := repository.AbnormalDeleteAllCards(ctx, tx, &abnormalEntity); err != nil {
			return fmt.Errorf("Failed to delete cards: %s", err.Msg)
		}

		// 방 삭제
		if err := repository.AbnormalDeleteRoom(ctx, tx, &abnormalEntity); err != nil {
			return fmt.Errorf("Failed to delete room: %s", err.Msg)
		}
		// 방 유저 정보 삭제
		if err := repository.AbnormalDeleteRoomUsers(ctx, tx, &abnormalEntity); err != nil {
			return fmt.Errorf("Failed to delete room users: %s", err.Msg)
		}

		return nil
	})
	// 에러 처리
	if err != nil {
		fmt.Printf("Cleanup error: %v\n", err)
	}
	// 레디스 세션값 삭제
	newErr := repository.RedisSessionDelete(ctx, sessionID)
	if newErr != nil {
		fmt.Printf("Failed to delete session: %v\n", newErr.Msg)
	}

	// 클라이언트 연결 종료 및 제거
	if sessionIDs, ok := entity.RoomSessions[roomID]; ok {
		for _, id := range sessionIDs {
			if client, exists := entity.WSClients[id]; exists {
				client.Close()
				delete(entity.WSClients, id)
			}
		}
	}

	// 방 삭제
	if sessionIDs, ok := entity.RoomSessions[roomID]; ok && len(sessionIDs) == 0 {
		delete(entity.RoomSessions, roomID)
		fmt.Printf("Room %d deleted as it is empty.\n", roomID)
	}

	// 타이머 삭제
	reconnectTimers.Delete(sessionID)
}

// 유저 ID를 sessionID로부터 가져오는 함수
func getUserIDFromSessionID(sessionID string) uint {
	if client, exists := entity.WSClients[sessionID]; exists {
		return client.UserID
	}
	return 0
}
