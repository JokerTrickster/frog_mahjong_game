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
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func RequestWinEventWebsocket(msg *entity.WSMessage) {
	//유저 상태를 변경한다. (대기실로 이동)
	ctx := context.Background()
	uID := msg.UserID
	roomID := msg.RoomID
	decryptedMessage, err := utils.DecryptAES(msg.Message)
	if err != nil {
		log.Fatalf("AES 복호화 에러: %s", err)
	}
	//string to struct
	req := request.ReqV2WSWinEvent{}
	err = json.Unmarshal([]byte(decryptedMessage), &req)
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
			roomInfoMsg.ErrorInfo = err
			ErrorHandling(msg, &roomInfoMsg)
			return fmt.Errorf("%s", err.Msg)
		}

		// 유저 상태 변경
		err = repository.RequestWinUpdateRoomUsers(ctx, tx, &requestWinEntity)
		if err != nil {
			roomInfoMsg.ErrorInfo = err
			ErrorHandling(msg, &roomInfoMsg)
			return fmt.Errorf("%s", err.Msg)
		}
		preloadUsers, err = repository.RequestWinFindAllRoomUsers(ctx, tx, roomID)
		if err != nil {
			roomInfoMsg.ErrorInfo = err
			ErrorHandling(msg, &roomInfoMsg)
			return fmt.Errorf("%s", err.Msg)
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

	message, err := CreateMessage(&roomInfoMsg)
	if err != nil {
		fmt.Println(err)
		return
	}
	msg.Message = message
	sendMessageToClients(roomID, msg)
}
