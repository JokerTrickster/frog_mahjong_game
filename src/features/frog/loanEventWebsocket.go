package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"main/features/frog/model/entity"
	_errors "main/features/frog/model/errors"
	"main/features/frog/model/request"
	"main/features/frog/repository"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func LoanEventWebsocket(msg *entity.WSMessage) *entity.ErrorInfo {
	ctx := context.Background()
	uID := msg.UserID
	roomID := msg.RoomID
	//string to struct
	req := request.ReqWSLoan{}
	err := json.Unmarshal([]byte(msg.Message), &req)
	if err != nil {
		return CreateErrorMessage(_errors.ErrCodeBadRequest, _errors.ErrUnmarshalFailed, "JSON 언마샬링 에러")
	}
	loanEntity := entity.WSLoanEntity{
		RoomID:       roomID,
		CardID:       req.CardID,
		TargetUserID: req.TargetUserID,
		UserID:       uID,
	}

	// 비즈니스 로직
	roomInfoMsg := entity.RoomInfo{}
	preloadUsers := []entity.RoomUsers{}
	var errInfo *entity.ErrorInfo
	err = mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		// loan 가능한지 체크 (마지막으로 버려진 카드인지 체크)
		errInfo = repository.LoanCheckLoan(ctx, tx, &loanEntity)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}

		// loan 하기 (상대방이 버린 카드를 가져온다)
		errInfo = repository.LoanCardLoan(ctx, tx, &loanEntity)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}

		// 룸 유저 카드 수와 상태값을 변경한다.
		errInfo = repository.LoanUpdateRoomUserCardCount(ctx, tx, &loanEntity)
		if errInfo != nil {
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

	//론한 유저에 대한 정보를 게임정보에 저장한다.
	LoanInfo := entity.LoanInfo{
		CardID:       int(req.CardID),
		UserID:       uID,
		TargetUserID: req.TargetUserID,
	}
	roomInfoMsg.GameInfo.LoanInfo = &LoanInfo
	message, err := CreateMessage(&roomInfoMsg)
	if err != nil {
		return CreateErrorMessage(_errors.ErrCodeInternal, err.Error(), _errors.ErrGameTerminated)
	}
	msg.Message = message
	sendMessageToClients(roomID, msg)
	return nil
}
