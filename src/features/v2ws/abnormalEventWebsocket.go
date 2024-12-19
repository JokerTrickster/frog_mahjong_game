package v2ws

import (
	"context"
	"fmt"
	"main/features/v2ws/model/entity"
	_errors "main/features/v2ws/model/errors"
	"main/features/v2ws/repository"
	"main/utils/db/mysql"
	"sync"
	"time"

	"gorm.io/gorm"
)

var reconnectTimers sync.Map // 재접속 타이머를 관리하는 맵

// 비정상적인 에러를 처리하는 함수
func AbnormalErrorHandling(roomID uint, sessionID string) {
	ctx := context.TODO()
	roomInfoMsg := entity.RoomInfo{}
	preloadUsers := []entity.RoomUsers{}

	// 비정상적인 유저 상태 처리
	err := mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		abnormalEntity := entity.WSAbnormalEntity{
			RoomID:         roomID,
			AbnormalUserID: getUserIDFromSessionID(sessionID),
		}

		// 유저 상태 변경
		if err := repository.AbnormalUpdateUsers(ctx, tx, &abnormalEntity); err != nil {
			return err
		}

		// 방 유저 정보 조회
		users, err := repository.AbnormalFindAllRoomUsers(ctx, tx, roomID)
		if err != nil {
			return err
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
	fmt.Println("비정상 에러 핸들러 호출 세션 ID : ", sessionID)

	// 클라이언트에 메시지 전송
	roomInfoMsg = *CreateRoomInfoMSG(ctx, preloadUsers, 1, roomInfoMsg.ErrorInfo, 0)
	roomInfoMsg.GameInfo.AllReady = false
	sendMessageToClients(roomID, roomInfoMsg)

	// 재접속 대기 시작
	waitForReconnection(roomID, sessionID, preloadUsers)
}

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

// 클라이언트에 메시지 전송
func sendMessageToClients(roomID uint, roomInfoMsg entity.RoomInfo) {
	message, err := CreateMessage(&roomInfoMsg)
	if err != nil {
		fmt.Printf("Failed to create message: %v\n", err)
		return
	}

	// 방에 있는 모든 클라이언트에 메시지 전송
	if sessionIDs, ok := entity.RoomSessions[roomID]; ok {
		for _, sessionID := range sessionIDs {
			if client, exists := entity.WSClients[sessionID]; exists {
				if err := client.Conn.WriteJSON(entity.WSMessage{Message: message}); err != nil {
					fmt.Printf("Failed to send message to session %s: %v\n", sessionID, err)
					client.Close()
					delete(entity.WSClients, sessionID)
					removeSessionFromRoom(roomID, sessionID)
				}
			}
		}
	}
}

// 세션 정리 (재접속 실패 시 호출)
func cleanupSession(roomID uint, sessionID string, preloadUsers []entity.RoomUsers) {
	ctx := context.TODO()
	fmt.Printf("Cleaning up session %s for room %d...\n", sessionID, roomID)

	err := mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		abnormalEntity := entity.WSAbnormalEntity{
			RoomID:         roomID,
			AbnormalUserID: getUserIDFromSessionID(sessionID),
		}

		// 카드 삭제
		if err := repository.AbnormalDeleteAllCards(ctx, tx, &abnormalEntity); err != nil {
			return fmt.Errorf("Failed to delete cards: %w", err)
		}

		// 방 삭제
		if err := repository.AbnormalDeleteRoom(ctx, tx, &abnormalEntity); err != nil {
			return fmt.Errorf("Failed to delete room: %w", err)
		}

		return nil
	})

	// 에러 처리
	if err != nil {
		fmt.Printf("Cleanup error: %v\n", err)
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
