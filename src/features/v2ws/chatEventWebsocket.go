package v2ws

import (
	"context"
	"log"
	"main/features/v2ws/model/entity"
	_errors "main/features/v2ws/model/errors"
	"main/features/v2ws/model/request"
	"main/features/v2ws/repository"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func ChatEventWebsocket(msg *entity.WSMessage) {
	//유저 상태를 변경한다. (대기실로 이동)
	ctx := context.Background()
	uID := msg.UserID
	roomID := msg.RoomID
	//string to struct
	req := request.ReqWSChat{
		UserID:  uID,
		RoomID:  roomID,
		Name:    msg.Name,
		Message: msg.Message,
	}

	// 비즈니스 로직
	ChatInfo := entity.ChatInfo{}
	var ChatID uint
	err := mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {

		// 채팅 DTO 만들기
		chatDTO := CreateChatDTO(req)

		// 채팅 정보 저장
		chatID, err := repository.ChatInsertOneChat(ctx, tx, chatDTO)
		if err != nil {
			return err
		}
		ChatID = chatID

		return nil
	})
	if err != nil {
		ChatInfo.ErrorInfo = &entity.ErrorInfo{
			Code: 500,
			Msg:  err.Error(),
			Type: _errors.ErrInternalServer,
		}
	}
	// 유저 상태를 변경한다. (방에 참여)
	if sessionIDs, ok := entity.RoomSessions[msg.RoomID]; ok {
		// 에러 발생 시 이벤트 요청한 유저에게만 메시지를 전달한다.
		if ChatInfo.ErrorInfo != nil || err != nil {
			for _, sessionID := range sessionIDs {
				if client, exists := entity.WSClients[sessionID]; exists && client.UserID == msg.UserID {
					// 메시지 설정
					msg.Message = req.Message
					msg.ChatID = ChatID

					// 메시지 전송
					err := client.Conn.WriteJSON(msg)
					if err != nil {
						log.Printf("Error sending message to user %d: %v", client.UserID, err)

						// 클라이언트를 종료 및 정리
						client.Close()
						delete(entity.WSClients, sessionID)
						removeSessionFromRoom(client.RoomID, sessionID)
					}
				}
			}
		} else {
			// 정상적인 경우 모든 유저에게 메시지 브로드캐스트
			for _, sessionID := range sessionIDs {
				if client, exists := entity.WSClients[sessionID]; exists {
					// 메시지 설정
					msg.Message = req.Message
					msg.ChatID = ChatID

					// 메시지 전송
					err := client.Conn.WriteJSON(msg)
					if err != nil {
						log.Printf("Error broadcasting to user %d: %v", client.UserID, err)

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
