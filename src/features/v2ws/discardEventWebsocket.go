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

func DiscardCardsEventWebsocket(msg *entity.WSMessage) {
	//유저 상태를 변경한다. (대기실로 이동)
	ctx := context.Background()
	uID := msg.UserID
	roomID := msg.RoomID

	//string to struct
	req := request.ReqWSDiscardCards{}
	err := json.Unmarshal([]byte(msg.Message), &req)
	if err != nil {
		log.Fatalf("JSON 언마샬링 에러: %s", err)
	}
	DiscardCardsEntity := entity.WSDiscardCardsEntity{
		RoomID: roomID,
		UserID: uID,
		CardID: uint(req.CardID),
	}

	// 비즈니스 로직
	// 보유 카드수가 4장인지 체크
	cardCount, err := repository.DiscardCardsOwnerCardCount(ctx, roomID, uID)
	if err != nil {
		fmt.Println(err)
		return
	}
	if cardCount != 4 {
		fmt.Println("보유 카드수가 4장이 아닙니다.")
		return
	}

	roomInfoMsg := entity.RoomInfo{}
	preloadUsers := []entity.RoomUsers{}
	err = mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		// 카드 상태 없데이트
		err := repository.DiscardCardsUpdateCardState(ctx, tx, &DiscardCardsEntity)
		if err != nil {
			return err
		}
		// 소유 카드 수 업데이트
		// 유저id로 room_users에서 찾아서 card_count를 뺀 후 업데이트 한다.
		err = repository.DiscardCardsUpdateRoomUserCardCount(ctx, tx, &DiscardCardsEntity)
		if err != nil {
			return err
		}

		preloadUsers, err = repository.DiscardCardsFindAllRoomUsers(ctx, tx, roomID)
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
		ErrorHandling(msg, roomID, uID, &roomInfoMsg)
	}
	// 유저 상태를 변경한다. (방에 참여)
	if sessionIDs, ok := entity.RoomSessions[msg.RoomID]; ok {
		// 게임 턴 계산
		playTurn := CalcPlayTurn(req.PlayTurn, len(sessionIDs))
		roomInfoMsg := *DiscardCreateRoomInfoMSG(ctx, preloadUsers, playTurn, roomInfoMsg.ErrorInfo, int(req.CardID))

		// 모든 유저가 카드를 선택했을 때
		if roomInfoMsg.GameInfo.AllPicked {
			// 카드 상태를 picked -> owned로 변경
			err := mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
				err := repository.DiscardCardUpdateAllCardState(ctx, tx, msg.RoomID)
				if err != nil {
					return err
				}
				preloadUsers, err = repository.DiscardCardsFindAllRoomUsers(ctx, tx, msg.RoomID)
				if err != nil {
					return err
				}
				return nil
			})

			// 에러 처리
			if err != nil {
				roomInfoMsg.ErrorInfo = &entity.ErrorInfo{
					Code: 500,
					Msg:  err.Error(),
					Type: _errors.ErrInternalServer,
				}
				ErrorHandling(msg, msg.RoomID, msg.UserID, &roomInfoMsg)
			}

			// 게임 상태 갱신
			roomInfoMsg = *DiscardCreateRoomInfoMSG(ctx, preloadUsers, playTurn, roomInfoMsg.ErrorInfo, int(req.CardID))
			roomInfoMsg.GameInfo.AllPicked = true
		}

		// 방의 모든 유저에게 메시지 전송
		for _, sessionID := range sessionIDs {
			if client, exists := entity.WSClients[sessionID]; exists {
				filterRoomInfoMsg := Deepcopy(roomInfoMsg)

				// 구조체를 JSON 문자열로 변환 (마샬링)
				message, err := CreateMessage(&filterRoomInfoMsg)
				if err != nil {
					fmt.Println(err)
					continue
				}

				msg.Message = message
				err = client.Conn.WriteJSON(msg)
				if err != nil {
					log.Printf("Error sending message to user %d: %v", client.UserID, err)
					client.Close()

					// 클라이언트 정리
					delete(entity.WSClients, sessionID)
					removeSessionFromRoom(client.RoomID, sessionID)
				}
			}
		}

		// 방이 비어 있으면 삭제
		if len(entity.RoomSessions[msg.RoomID]) == 0 {
			delete(entity.RoomSessions, msg.RoomID)
		}
	}

}
