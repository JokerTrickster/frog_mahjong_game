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

func FailedLoanEventWebsocket(msg *entity.WSMessage) {
	ctx := context.Background()
	uID := msg.UserID
	roomID := msg.RoomID

	//string to struct
	req := request.ReqWSFailedLoan{}
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
		// 소유하고 있는 카드인지 체크
		newErr := repository.FailedLoanCheckCard(ctx, tx, &loanEntity)
		if newErr != nil {
			SendErrorMessage(msg, newErr)
			return fmt.Errorf("%s", newErr.Msg)
		}

		// 카드 정보를 롤백한다.
		newErr = repository.FailedLoanRollbackCard(ctx, tx, &loanEntity)
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

	// 론 가능 여부를 true로 변경
	roomInfoMsg.GameInfo.IsLoanAllowed = true

	// 론 실패한 유저ID 저장
	roomInfoMsg.GameInfo.FailedLoanUserID = uID
	// 구조체를 JSON 문자열로 변환 (마샬링)
	message, err := CreateMessage(&roomInfoMsg)
	if err != nil {
		fmt.Println(err)
	}
	msg.Message = message
	sendMessageToClients(roomID, msg)
}
