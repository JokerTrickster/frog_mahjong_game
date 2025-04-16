package slime_war

import (
	"context"
	"encoding/json"
	"fmt"
	"main/features/slime_war/model/entity"
	_errors "main/features/slime_war/model/errors"
	"main/features/slime_war/model/request"
	"main/features/slime_war/repository"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func HintItemEventWebsocket(msg *entity.WSMessage) *entity.ErrorInfo {
	//유저 상태를 변경한다. (대기실로 이동)
	ctx := context.Background()
	roomID := msg.RoomID
	req := request.ReqWSHintItem{}
	err := json.Unmarshal([]byte(msg.Message), &req)
	if err != nil {
		return CreateErrorMessage(_errors.ErrCodeBadRequest, _errors.ErrUnmarshalFailed, "JSON 언마샬링 에러")
	}
	// 비즈니스 로직
	//해당 방이 대기상태인지 체크한다.
	preloadUsers := []entity.PreloadUsers{}
	messageMsg := entity.MessageInfo{}
	var errInfo *entity.ErrorInfo
	position := &entity.Position{}
	// 힌트 아이템 사용 가능한지 체크
	roomSettings, errInfo := repository.HintItemCheck(ctx, roomID)
	if errInfo != nil {
		return CreateErrorMessage(_errors.ErrCodeBadRequest, _errors.ErrHintItemFailed, "힌트 아이템 사용 불가")
	}
	if roomSettings.ItemHintCount == 0 {
		return CreateErrorMessage(_errors.ErrCodeBadRequest, _errors.ErrHintItemFailed, "힌트 아이템 사용 불가")
	}

	err = mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		// 힌트 아이템 1 감소
		errInfo := repository.HintItemDecrease(ctx, tx, roomSettings)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}
		// 힌트로 좌표를 가져온다. (랜덤)
		// 현재 정답을 맞춘 좌표를 가져온다.
		userCorrectPositionDTOList, errInfo := repository.HintItemFindCorrectPosition(ctx, tx, roomID, req.Round, req.ImageID)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}

		// 찾은 좌표 ID 리스트를 만든다.
		var correctPositionIDList []uint
		for _, userCorrectPositionDTO := range userCorrectPositionDTOList {
			correctPositionIDList = append(correctPositionIDList, uint(userCorrectPositionDTO.CorrectPositionID))
		}

		// 못찾은 좌표 하나를 가져온다.
		imageCorrectPositionDTO, errInfo := repository.HintItemFindOneCorrectPosition(ctx, tx, uint(req.ImageID), correctPositionIDList)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}

		position.X = imageCorrectPositionDTO.XPosition
		position.Y = imageCorrectPositionDTO.YPosition

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
	//힌트 좌표 추가
	messageMsg.GameInfo.HintPosition = position
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
