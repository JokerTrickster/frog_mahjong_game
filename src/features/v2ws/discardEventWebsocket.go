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
	}

	//유저 상태를 변경한다. (방에 참여)
	if clients, ok := entity.WSClients[msg.RoomID]; ok {
		// 메시지 생성
		//게임턴 계산
		playTurn := CalcPlayTurn(req.PlayTurn, len(entity.WSClients[msg.RoomID]))
		roomInfoMsg = *DiscardCreateRoomInfoMSG(ctx, preloadUsers, playTurn, roomInfoMsg.ErrorInfo, int(req.CardID))
		if roomInfoMsg.GameInfo.AllPicked == true {
			// 카드 상태 picked -> owned 로 변경한다.
			// 모든 유저가 카드를 선택했을 때, 모든 유저의 카드 상태를 picked -> owned 로 변경한다.
			err = mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
				err = repository.DiscardCardUpdateAllCardState(ctx, tx, roomID)
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
			}
			roomInfoMsg = *DiscardCreateRoomInfoMSG(ctx, preloadUsers, playTurn, roomInfoMsg.ErrorInfo, int(req.CardID))
			roomInfoMsg.GameInfo.AllReady = true
		}
		//에러 발생시 이벤트 요청한 유저에게만 메시지를 전달한다.
		if roomInfoMsg.ErrorInfo != nil || err != nil {
			for client := range clients {
				if clients[client].UserID == msg.UserID {
					// 구조체를 JSON 문자열로 변환 (마샬링)
					message, err := CreateMessage(&roomInfoMsg)
					if err != nil {
						fmt.Println(err)
					}
					msg.Message = message
					err = client.WriteJSON(msg)
					if err != nil {
						log.Printf("error: %v", err)
						client.Close()
						delete(clients, client)
					}
				}
			}
		} else {
			for client := range clients {
				filterRoomInfoMsg := Deepcopy(roomInfoMsg)

				// 구조체를 JSON 문자열로 변환 (마샬링)
				message, err := CreateMessage(&filterRoomInfoMsg)
				if err != nil {
					fmt.Println(err)
				}
				msg.Message = message
				err = client.WriteJSON(msg)
				if err != nil {
					log.Printf("error: %v", err)
					client.Close()
					delete(clients, client)
				}
			}
		}
	}
}
