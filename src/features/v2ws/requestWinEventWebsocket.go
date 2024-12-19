package v2ws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"main/features/v2ws/model/entity"
	_errors "main/features/v2ws/model/errors"
	"main/features/v2ws/model/request"
	"main/features/v2ws/repository"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func RequestWinEventWebsocket(msg *entity.WSMessage) {
	//유저 상태를 변경한다. (대기실로 이동)
	ctx := context.Background()
	uID := msg.UserID
	roomID := msg.RoomID

	//string to struct
	req := request.ReqV2WSWinEvent{}
	err := json.Unmarshal([]byte(msg.Message), &req)
	if err != nil {
		log.Fatalf("JSON 언마샬링 에러: %s", err)
	}

	requestWinEntity := entity.V2WSRequestWinEntity{
		RoomID: roomID,
		UserID: uID,
	}
	for _, card := range req.Cards {
		requestWinEntity.Cards = append(requestWinEntity.Cards, int(card.CardID))
	}

	// 비즈니스 로직
	roomInfoMsg := entity.RoomInfo{}
	preloadUsers := []entity.RoomUsers{}
	err = mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		// 카드 정보 체크 (소유하고 있는지 체크)
		_, err := repository.RequestWinFindAllCards(ctx, tx, &requestWinEntity)
		if err != nil {
			return err
		}

		// 유저 상태 변경
		err = repository.RequestWinUpdateRoomUsers(ctx, tx, &requestWinEntity)
		if err != nil {
			return err
		}
		preloadUsers, err = repository.RequestWinFindAllRoomUsers(ctx, tx, roomID)
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
	roomInfoMsg.GameInfo.Winner = uID

	//승리 유저 카드 정보 순서 저장
	cards := []*entity.Card{}
	for _, card := range req.Cards {
		cards = append(cards, &entity.Card{
			CardID: card.CardID,
			UserID: uID,
		})
	}
	for i := 0; i < len(roomInfoMsg.Users); i++ {
		if roomInfoMsg.Users[i].ID == uID {
			roomInfoMsg.Users[i].Cards = cards
			break
		}
	}

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
					// 메시지 전송
					err := client.Conn.WriteJSON(msg)
					if err != nil {
						fmt.Printf("Error sending message to user %d: %v\n", client.UserID, err)
						client.Conn.Close()

						// 클라이언트와 세션 삭제
						delete(entity.WSClients, sessionID)
						removeSessionFromRoom(client.RoomID, sessionID)
					}
				}
			}
		} else {
			// 정상적인 경우 모든 유저에게 메시지 브로드캐스트
			for _, sessionID := range sessionIDs {
				if client, exists := entity.WSClients[sessionID]; exists {
					// 메시지 전송
					err := client.Conn.WriteJSON(msg)
					if err != nil {
						fmt.Printf("Error broadcasting to user %d: %v\n", client.UserID, err)
						client.Conn.Close()

						// 클라이언트와 세션 삭제
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
