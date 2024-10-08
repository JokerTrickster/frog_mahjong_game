package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"main/features/ws/model/entity"
	_errors "main/features/ws/model/errors"
	"main/features/ws/model/request"
	"main/features/ws/repository"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func SuccessLoanEventWebsocket(msg *entity.WSMessage) {
	ctx := context.Background()
	uID := msg.UserID
	roomID := msg.RoomID

	//string to struct
	req := request.ReqWSSuccessEvent{}
	err := json.Unmarshal([]byte(msg.Message), &req)
	if err != nil {
		log.Fatalf("JSON 언마샬링 에러: %s", err)
	}
	successEntity := entity.WSSuccessEntity{
		RoomID: roomID,
		UserID: uID,
		Score:  req.Score,
		LoanInfo: &entity.ReqSuccessLoanInfo{
			TargetUserID: req.LoanInfo.TargetUserID,
			CardID:       req.LoanInfo.CardID,
		},
	}
	for _, card := range req.Cards {
		successEntity.Cards = append(successEntity.Cards, int(card.CardID))
	}

	// 비즈니스 로직
	roomInfoMsg := entity.RoomInfo{}
	doraDTO := &mysql.Cards{}
	preloadUsers := []entity.RoomUsers{}
	err = mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		// 카드 정보 체크 (소유하고 있는지 체크)
		cards, err := repository.SuccessFindAllCards(ctx, tx, &successEntity)
		if err != nil {
			return err
		}
		// 카드 정보로 점수 체크한다.
		err = CalcScore(cards, successEntity.Score)
		if err != nil {
			return err
		}
		// 론인 경우 해당 유저에 코인 차감한다.
		err = repository.SuccessLoanDiffCoin(ctx, tx, &successEntity)
		if err != nil {
			return err
		}
		// 론인 경우 해당 유저에 코인 추가한다.
		err = repository.SuccessLoanAddCoin(ctx, tx, &successEntity)
		if err != nil {
			return err
		}

		// 유저 상태 변경
		err = repository.SuccessUpdateRoomUsers(ctx, tx, &successEntity)
		if err != nil {
			return err
		}
		preloadUsers, err = repository.SuccessFindAllRoomUsers(ctx, tx, roomID)
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

	//유저들에게 메시지 전송한다.
	if clients, ok := entity.WSClients[msg.RoomID]; ok {
		// 메시지 생성
		roomInfoMsg = *CreateRoomInfoMSG(ctx, preloadUsers, req.PlayTurn, roomInfoMsg.ErrorInfo)

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

		// 론 가능 여부를 true로 변경
		roomInfoMsg.GameInfo.IsLoanAllowed = true

		//도라 카드 정보 저장
		doraCardInfo := entity.Card{}
		doraCardInfo.CardID = uint(doraDTO.CardID)
		roomInfoMsg.GameInfo.Dora = &doraCardInfo

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
