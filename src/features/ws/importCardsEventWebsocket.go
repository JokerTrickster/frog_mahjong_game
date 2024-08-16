package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"main/features/ws/model/entity"
	"main/features/ws/model/request"
	"main/features/ws/repository"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func ImportCardsEventWebsocket(msg *entity.WSMessage) {
	//유저 상태를 변경한다. (대기실로 이동)
	ctx := context.Background()
	uID := msg.UserID
	roomID := msg.RoomID

	//string to struct
	req := request.ReqWSImportCards{}
	err := json.Unmarshal([]byte(msg.Message), &req)
	if err != nil {
		log.Fatalf("JSON 언마샬링 에러: %s", err)
	}
	importCardsEntity := entity.WSImportCardsEntity{
		RoomID: roomID,
		UserID: uID,
	}
	for _, card := range req.Cards {
		importCardsEntity.Cards = append(importCardsEntity.Cards, &mysql.Cards{
			CardID: int(card.CardID),
			RoomID: int(roomID),
			UserID: int(uID),
		})
	}

	// 비즈니스 로직
	roomInfoMsg := entity.RoomInfo{}
	doraDTO := &mysql.Cards{}
	preloadUsers := []entity.RoomUsers{}
	err = mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		// 카드 상태 없데이트
		err := repository.ImportCardsUpdateCardState(ctx, tx, &importCardsEntity)
		if err != nil {
			return err
		}
		// 소유 카드 수 업데이트
		// 유저id로 room_users에서 찾아서 card_count를 더한 후 업데이트 한다.
		err = repository.ImportCardsUpdateRoomUserCardCount(ctx, tx, &importCardsEntity)
		if err != nil {
			return err
		}
		doraDTO, err = repository.ImportCardsFindOneDora(ctx, tx, roomID)
		if err != nil {
			return err
		}
		preloadUsers, err = repository.ImportCardsFindAllRoomUsers(ctx, tx, roomID)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		roomInfoMsg.ErrorInfo = &entity.ErrorInfo{
			Code: 500,
			Msg:  err.Error(),
		}
	}

	// 메시지 생성

	//게임턴 계산
	playTurn := CalcPlayTurn(req.PlayTurn, len(entity.WSClients[msg.RoomID]))
	roomInfoMsg = *CreateRoomInfoMSG(ctx, preloadUsers, playTurn, roomInfoMsg.ErrorInfo)

	//카드 정보 저장
	doraCardInfo := entity.Card{}
	doraCardInfo.CardID = uint(doraDTO.CardID)
	roomInfoMsg.GameInfo.Dora = &doraCardInfo
	// 구조체를 JSON 문자열로 변환 (마샬링)
	message, err := CreateMessage(&roomInfoMsg)
	if err != nil {
		fmt.Println(err)
	}
	msg.Message = message

	//유저 상태를 변경한다. (방에 참여)
	if clients, ok := entity.WSClients[msg.RoomID]; ok {
		//에러 발생시 이벤트 요청한 유저에게만 메시지를 전달한다.
		if roomInfoMsg.ErrorInfo != nil || err != nil {
			for client := range clients {
				if clients[client].UserID == msg.UserID {
					err := client.WriteJSON(msg)
					if err != nil {
						fmt.Printf("error: %v", err)
						client.Close()
						delete(clients, client)
					}
				}
			}
		} else {
			for client := range clients {
				err := client.WriteJSON(msg)
				if err != nil {
					fmt.Printf("error: %v", err)
					client.Close()
					delete(clients, client)
				}
			}
		}
	}
}
