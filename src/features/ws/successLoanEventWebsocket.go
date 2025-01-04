package ws

import (
	"context"
	"encoding/json"
	"fmt"
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
		SendErrorMessage(msg, CreateErrorMessage(_errors.ErrCodeBadRequest, _errors.ErrUnmarshalFailed, "JSON 언마샬링 에러"))
		return
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
	preloadUsers := []entity.RoomUsers{}
	err = mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		// 카드 정보 체크 (소유하고 있는지 체크)
		cards, newErr := repository.SuccessFindAllCards(ctx, tx, &successEntity)
		if newErr != nil {
			SendErrorMessage(msg, newErr)
			return fmt.Errorf("%s", newErr.Msg)
		}
		// 카드 정보로 점수 체크한다.
		newErr = CalcScore(cards, successEntity.Score)
		if newErr != nil {
			SendErrorMessage(msg, newErr)
			return fmt.Errorf("%s", newErr.Msg)
		}
		// 론인 경우 해당 유저에 코인 차감한다.
		newErr = repository.SuccessLoanDiffCoin(ctx, tx, &successEntity)
		if newErr != nil {
			SendErrorMessage(msg, newErr)
			return fmt.Errorf("%s", newErr.Msg)
		}
		// 론인 경우 해당 유저에 코인 추가한다.
		newErr = repository.SuccessLoanAddCoin(ctx, tx, &successEntity)
		if newErr != nil {
			SendErrorMessage(msg, newErr)
			return fmt.Errorf("%s", newErr.Msg)
		}

		// 유저 상태 변경
		newErr = repository.SuccessUpdateRoomUsers(ctx, tx, &successEntity)
		if newErr != nil {
			SendErrorMessage(msg, newErr)
			return fmt.Errorf("%s", newErr.Msg)
		}
		preloadUsers, newErr = repository.PreloadFindGameInfo(ctx, tx, roomID)
		if newErr != nil {
			SendErrorMessage(msg, newErr)
			return fmt.Errorf("%s", newErr.Msg)
		}

		return nil
	})
	if err != nil {
		return
	}
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
	message, err := CreateMessage(&roomInfoMsg)
	if err != nil {
		fmt.Println(err)
	}
	msg.Message = message
	sendMessageToClients(roomID, msg)
}
