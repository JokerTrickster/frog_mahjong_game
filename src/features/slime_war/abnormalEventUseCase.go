package slime_war

import (
	"context"
	"fmt"
	"main/features/slime_war/model/entity"
	_errors "main/features/slime_war/model/errors"
	"main/features/slime_war/repository"
	"main/utils/db/mysql"
	"time"

	"gorm.io/gorm"
)

// 재접속 대기 타이머 (20초)
func waitForReconnection(roomID uint, sessionID string, preloadUsers []entity.PreloadUsers) {
	fmt.Printf("Waiting for session %s to reconnect in room %d...\n", sessionID, roomID)

	// 타이머 설정
	timer := time.AfterFunc(3*time.Second, func() {
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
	})

	// 타이머 저장
	reconnectTimers.Store(sessionID, timer)
}

// 세션 정리 (재접속 실패 시 호출)
func cleanupSession(roomID uint, sessionID string, preloadUsers []entity.PreloadUsers) {
	ctx := context.TODO()

	msg := &entity.WSMessage{
		RoomID: roomID,
		Event:  "ERROR",
	}
	errMsg := CreateErrorMessage(_errors.ErrCodeInternal, _errors.ErrAbnormalExit, "상대방이 연결이 끊겼습니다. 강제로 게임을 종료합니다.")
	MessageInfo := &entity.MessageInfo{} // Initialize MessageInfo
	MessageInfo.ErrorInfo = errMsg
	sendMsg, _ := CreateMessage(MessageInfo)
	msg.Message = sendMsg
	sendMessageToClients(roomID, msg)
	var errInfo *entity.ErrorInfo
	err := mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		abnormalEntity := entity.WSAbnormalEntity{
			RoomID:         roomID,
			AbnormalUserID: getUserIDFromSessionID(sessionID),
		}

		// 방 카드 삭제
		if errInfo = repository.AbnormalDeleteAllCards(ctx, tx, &abnormalEntity); errInfo != nil {
			return fmt.Errorf("Failed to delete cards: %s", errInfo.Msg)
		}

		// 방 맵 삭제
		if errInfo = repository.AbnormalDeleteAllMaps(ctx, tx, &abnormalEntity); errInfo != nil {
			return fmt.Errorf("Failed to delete maps: %s", errInfo.Msg)
		}

		// 방 게임 셋팅 삭제
		if errInfo = repository.AbnormalDeleteGameRoomSetting(ctx, tx, &abnormalEntity); errInfo != nil {
			return fmt.Errorf("Failed to delete game room setting: %s", errInfo.Msg)
		}

		// 방 유저 삭제
		if errInfo = repository.AbnormalDeleteRoomUsers(ctx, tx, &abnormalEntity); errInfo != nil {
			return fmt.Errorf("Failed to delete room users: %s", errInfo.Msg)
		}

		// 방 삭제
		if errInfo = repository.AbnormalDeleteRoom(ctx, tx, &abnormalEntity); errInfo != nil {
			return fmt.Errorf("Failed to delete room: %s", errInfo.Msg)
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
}

// 유저 ID를 sessionID로부터 가져오는 함수
func getUserIDFromSessionID(sessionID string) uint {
	if client, exists := entity.WSClients[sessionID]; exists {
		return client.UserID
	}
	return 0
}
