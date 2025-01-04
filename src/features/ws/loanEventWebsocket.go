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

func LoanEventWebsocket(msg *entity.WSMessage) {
	ctx := context.Background()
	uID := msg.UserID
	roomID := msg.RoomID

	//string to struct
	req := request.ReqWSLoan{}
	err := json.Unmarshal([]byte(msg.Message), &req)
	if err != nil {
		SendErrorMessage(msg, CreateErrorMessage(_errors.ErrCodeBadRequest, _errors.ErrUnmarshalFailed, "JSON 언마샬링 에러"))
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
	err = mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		// loan 가능한지 체크 (마지막으로 버려진 카드인지 체크)
		newErr := repository.LoanCheckLoan(ctx, tx, &loanEntity)
		if newErr != nil {
			SendErrorMessage(msg, newErr)
			return fmt.Errorf("%s", newErr.Msg)
		}

		// loan 하기 (상대방이 버린 카드를 가져온다)
		newErr = repository.LoanCardLoan(ctx, tx, &loanEntity)
		if newErr != nil {
			SendErrorMessage(msg, newErr)
			return fmt.Errorf("%s", newErr.Msg)
		}

		// 룸 유저 카드 수와 상태값을 변경한다.
		newErr = repository.LoanUpdateRoomUserCardCount(ctx, tx, &loanEntity)
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

	//론한 유저에 대한 정보를 게임정보에 저장한다.
	LoanInfo := entity.LoanInfo{
		CardID:       int(req.CardID),
		UserID:       uID,
		TargetUserID: req.TargetUserID,
	}
	roomInfoMsg.GameInfo.LoanInfo = &LoanInfo
	message, err := CreateMessage(&roomInfoMsg)
	if err != nil {
		fmt.Println(err)
	}
	msg.Message = message
	sendMessageToClients(roomID, msg)
}
