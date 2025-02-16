package find_it

import (
	"context"
	"encoding/json"
	"fmt"
	"main/features/find_it/model/entity"
	_errors "main/features/find_it/model/errors"
	"main/features/find_it/model/request"
	"main/features/find_it/repository"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func SubmitPositionEventWebsocket(msg *entity.WSMessage) *entity.ErrorInfo {
	//유저 상태를 변경한다. (대기실로 이동)
	ctx := context.Background()
	uID := msg.UserID
	roomID := msg.RoomID
	req := request.ReqWSSubmitPosition{}
	err := json.Unmarshal([]byte(msg.Message), &req)
	if err != nil {
		return CreateErrorMessage(_errors.ErrCodeBadRequest, _errors.ErrUnmarshalFailed, "JSON 언마샬링 에러")
	}
	// 비즈니스 로직
	//해당 방이 대기상태인지 체크한다.
	preloadUsers := []entity.PreloadUsers{}
	messageMsg := entity.MessageInfo{}
	var errInfo *entity.ErrorInfo
	correctID := 0
	err = mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		// 해당 좌표가 정답인지 아닌지 체크
		var errInfo *entity.ErrorInfo
		correctID, errInfo = repository.SubmitPositionCheck(ctx, tx, req.ImageID, req.XPosition, req.YPosition)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}
		if correctID != 0 {
			// 정답일 경우
			UserCorrectPositionDTO := CreateUserCorrectPosition(roomID, uID, req.Round, req.ImageID, correctID)
			errInfo = repository.SubmitPositionCorrectSave(ctx, tx, UserCorrectPositionDTO)
			if errInfo != nil {
				return fmt.Errorf("%s", errInfo.Msg)
			}
		} else {
			// 목숨 1개를 줄인다.
			errInfo = repository.SubmitPositionLifeDecrease(ctx, tx, roomID, uID)
			if errInfo != nil {
				return fmt.Errorf("%s", errInfo.Msg)
			}

		}

		//TODO 30라운드 이미지를 선택해서 각 라운드마다 이미지를 만든다.
		preloadUsers, errInfo = repository.PreloadUsers(ctx, tx, roomID)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}
		return nil
	})
	if err != nil {
		return errInfo
	}

	// 메시지 생성
	messageMsg = *CreateMessageInfoMSG(ctx, preloadUsers, 1, messageMsg.ErrorInfo, 0)
	if correctID == 0 {
		wrongPosition := entity.Position{
			X: req.XPosition,
			Y: req.YPosition,
		}
		messageMsg.GameInfo.WrongPosition = wrongPosition
	}
	if len(preloadUsers) == 2 {
		messageMsg.GameInfo.IsFull = true
	}

	message, err := CreateMessage(&messageMsg)
	if err != nil {
		return CreateErrorMessage(_errors.ErrCodeBadRequest, _errors.ErrMarshalFailed, "메시지 생성 에러")
	}
	msg.Message = message
	sendMessageToClients(roomID, msg)
	return nil
}
