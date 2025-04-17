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

func TimeOutEventWebsocket(msg *entity.WSMessage) *entity.ErrorInfo {
	//유저 상태를 변경한다. (대기실로 이동)
	ctx := context.Background()
	roomID := msg.RoomID
	req := request.ReqWSTimeOut{}
	err := json.Unmarshal([]byte(msg.Message), &req)
	if err != nil {
		return CreateErrorMessage(_errors.ErrCodeBadRequest, _errors.ErrUnmarshalFailed, "JSON 언마샬링 에러")
	}
	// 비즈니스 로직
	//해당 방이 대기상태인지 체크한다.
	preloadUsers := []entity.PreloadUsers{}
	messageMsg := entity.MessageInfo{}
	var errInfo *entity.ErrorInfo

	err = mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		// 못찾은 좌표를 모두 가져온다.
		// 현재 정답을 맞춘 좌표를 가져온다.
		userCorrectPositionDTOList, errInfo := repository.TimeOutFindCorrectPosition(ctx, tx, roomID, req.Round, req.ImageID)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}

		// 개수만큼 목숨을 줄인다.
		diffLife := 5 - len(userCorrectPositionDTOList)
		// 목숨 차감한다.
		errInfo = repository.TimeOutLifeDecrease(ctx, tx, roomID, diffLife)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
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
	//타이머 아이템 사용
	messageMsg.SlimeWarGameInfo.TimerUsed = true
	if len(preloadUsers) == 2 {
		messageMsg.SlimeWarGameInfo.IsFull = true
	}
	if messageMsg.GameInfo.Life <= 0 {
		msg.Event = "GAME_OVER"
	} else {
		msg.Event = "ROUND_FAIL"
		//TODO 못맞춘 좌표를 보낸다.
		correctPositions, err := repository.TimeOutFindImageCorrectPosition(ctx, int(roomID), req.Round, req.ImageID)
		if err != nil {
			return CreateErrorMessage(_errors.ErrCodeBadRequest, _errors.ErrFetchFailed, "이미지 정답 위치 조회 에러")
		}
		positions := []*entity.Position{}
		for _, correctPosition := range correctPositions {
			p := &entity.Position{
				X: correctPosition.XPosition,
				Y: correctPosition.YPosition,
			}
			positions = append(positions, p)
		}
		messageMsg.GameInfo.FailedPositions = positions
	}
	message, err := CreateMessage(&messageMsg)
	if err != nil {
		return CreateErrorMessage(_errors.ErrCodeBadRequest, _errors.ErrMarshalFailed, "메시지 생성 에러")
	}

	msg.Message = message
	sendMessageToClients(roomID, msg)
	return nil
}
