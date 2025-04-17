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

	// 클라이언트에 메시지 전송
	MessageInfoMsg = *CreateMessageInfoMSG(ctx, preloadUsers, 1, MessageInfoMsg.ErrorInfo, 0)
	MessageInfoMsg.SlimeWarGameInfo.AllReady = false

	message, err := CreateMessage(&MessageInfoMsg)
	if err != nil {
		fmt.Println(err)
		return
	}
	msg := entity.WSMessage{
		RoomID:  roomID,
		UserID:  userID,
		Message: message,
		Event:   "DISCONNECT",
	}
	sendMessageToClients(roomID, &msg)

	// 재접속 대기 시작
	// waitForReconnection(roomID, sessionID, preloadUsers)
}
