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

func TimeOutDiscardCardsEventWebsocket(msg *entity.WSMessage) {
	//유저 상태를 변경한다. (대기실로 이동)
	ctx := context.Background()
	uID := msg.UserID
	roomID := msg.RoomID

	//string to struct
	req := request.ReqWSTimeOutDiscardCards{}
	err := json.Unmarshal([]byte(msg.Message), &req)
	if err != nil {
		log.Fatalf("JSON 언마샬링 에러: %s", err)
	}
	TimeOutDiscardCardsEntity := entity.WSTimeOutDiscardCardsEntity{
		RoomID: roomID,
		UserID: uID,
		CardID: uint(req.CardID),
	}

	// 비즈니스 로직
	roomInfoMsg := entity.RoomInfo{}
	doraDTO := &mysql.Cards{}
	preloadUsers := []entity.RoomUsers{}
	err = mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		// 카드 상태 없데이트
		err := repository.TimeOutDiscardCardsUpdateCardState(ctx, tx, &TimeOutDiscardCardsEntity)
		if err != nil {
			return err
		}

		doraDTO, err = repository.TimeOutDiscardCardsFindOneDora(ctx, tx, roomID)
		if err != nil {
			return err
		}
		preloadUsers, err = repository.TimeOutDiscardCardsFindAllRoomUsers(ctx, tx, roomID)
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
	// 유저 상태를 변경한다. (방에 참여)
	if sessionIDs, ok := entity.RoomSessions[msg.RoomID]; ok {
		// 메시지 생성
		// 게임 턴 계산
		playTurn := CalcPlayTurn(req.PlayTurn, len(sessionIDs))
		roomInfoMsg = *CreateRoomInfoMSG(ctx, preloadUsers, playTurn, roomInfoMsg.ErrorInfo, 0)

		// 카드 정보 저장
		doraCardInfo := entity.Card{}
		doraCardInfo.CardID = uint(doraDTO.CardID)

		// 에러 발생 시 이벤트 요청한 유저에게만 메시지 전달
		if roomInfoMsg.ErrorInfo != nil || err != nil {
			for _, sessionID := range sessionIDs {
				if client, exists := entity.WSClients[sessionID]; exists && client.UserID == msg.UserID {
					// 구조체를 JSON 문자열로 변환 (마샬링)
					message, err := CreateMessage(&roomInfoMsg)
					if err != nil {
						fmt.Println(err)
					}
					msg.Message = message
					err = client.Conn.WriteJSON(msg)
					if err != nil {
						log.Printf("Error sending message to user %d: %v", client.UserID, err)
						client.Conn.Close()
						delete(entity.WSClients, sessionID)
					}

					// 방에서도 해당 세션 제거
					removeSessionFromRoom(client.RoomID, sessionID)
				}
			}
		} else {
			// 정상적인 경우 모든 방 유저에게 메시지 브로드캐스트
			for _, sessionID := range sessionIDs {
				if client, exists := entity.WSClients[sessionID]; exists {
					// 각 클라이언트에 맞는 메시지 생성
					filterRoomInfoMsg := Deepcopy(roomInfoMsg)

					// 구조체를 JSON 문자열로 변환 (마샬링)
					message, err := CreateMessage(&filterRoomInfoMsg)
					if err != nil {
						fmt.Println(err)
					}
					msg.Message = message
					err = client.Conn.WriteJSON(msg)
					if err != nil {
						log.Printf("Error broadcasting to user %d: %v", client.UserID, err)
						client.Conn.Close()
						delete(entity.WSClients, sessionID)

						// 방에서도 해당 세션 제거
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
