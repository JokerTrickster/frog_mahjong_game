package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"main/features/frog/model/entity"
	_errors "main/features/frog/model/errors"
	"main/features/frog/model/request"
	"main/features/frog/repository"
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func SuccessLoanEventWebsocket(msg *entity.WSMessage) *entity.ErrorInfo {
	ctx := context.Background()
	uID := msg.UserID
	roomID := msg.RoomID
	decryptedMessage, err := utils.DecryptAES(msg.Message)
	if err != nil {
		return CreateErrorMessage(_errors.ErrCodeBadRequest, _errors.ErrCryptoFailed, "AES 복호화 에러")
	}
	//string to struct
	req := request.ReqWSSuccessEvent{}
	err = json.Unmarshal([]byte(decryptedMessage), &req)
	if err != nil {
		return CreateErrorMessage(_errors.ErrCodeBadRequest, _errors.ErrUnmarshalFailed, "JSON 언마샬링 에러")
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
	var errInfo *entity.ErrorInfo
	err = mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		// 카드 정보 체크 (소유하고 있는지 체크)
		cards, errInfo := repository.SuccessFindAllCards(ctx, tx, &successEntity)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}
		// 카드 정보로 점수 체크한다.
		errInfo = CalcScore(cards, successEntity.Score)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}

		// 유저 상태 변경
		errInfo = repository.SuccessUpdateRoomUsers(ctx, tx, &successEntity)
		if errInfo != nil {
			SendErrorMessage(msg, errInfo)
			return fmt.Errorf("%s", errInfo.Msg)
		}
		preloadUsers, errInfo = repository.PreloadFindGameInfo(ctx, tx, roomID)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}

		return nil
	})
	if err != nil {
		return errInfo
	}
	// 메시지 생성
	gameRoomSettings, errInfo := repository.FindOneFrogCurrentRound(ctx, roomID)
	if errInfo != nil {
		return errInfo
	}
	roomInfoMsg = *CreateRoomInfoMSG(ctx, preloadUsers, gameRoomSettings.CurrentRound, roomInfoMsg.ErrorInfo)

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
		return CreateErrorMessage(_errors.ErrCodeInternal, err.Error(), _errors.ErrGameTerminated)
	}
	msg.Message = message
	sendMessageToClients(roomID, msg)
	return nil
}
