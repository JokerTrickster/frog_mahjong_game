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

func ImportCardsEventWebsocket(msg *entity.WSMessage) *entity.ErrorInfo {
	//유저 상태를 변경한다. (대기실로 이동)
	ctx := context.Background()
	uID := msg.UserID
	roomID := msg.RoomID

	//string to struct
	req := request.ReqWSImportCards{}
	err := json.Unmarshal([]byte(msg.Message), &req)
	if err != nil {
		return CreateErrorMessage(_errors.ErrCodeBadRequest, _errors.ErrUnmarshalFailed, "JSON 언마샬링 에러")
	}
	importCardsEntity := entity.WSImportCardsEntity{
		RoomID: roomID,
		UserID: uID,
	}
	for _, card := range req.Cards {
		importCardsEntity.Cards = append(importCardsEntity.Cards, &mysql.FrogUserCards{
			CardID: int(card.CardID),
			RoomID: int(roomID),
			UserID: int(uID),
		})
	}

	// 비즈니스 로직
	roomInfoMsg := entity.RoomInfo{}
	preloadUsers := []entity.RoomUsers{}
	var errInfo *entity.ErrorInfo
	err = mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		// 카드 상태 없데이트
		errInfo := repository.ImportCardsUpdateCardState(ctx, tx, &importCardsEntity)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}
		// 소유 카드 수 업데이트
		// 유저id로 room_users에서 찾아서 card_count를 더한 후 업데이트 한다.
		errInfo = repository.ImportCardsUpdateRoomUserCardCount(ctx, tx, &importCardsEntity)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}
		errInfo = repository.UpdateRound(ctx, tx, roomID)

		preloadUsers, errInfo = repository.PreloadFindGameInfo(ctx, tx, roomID)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}
		return nil
	})
	if err != nil {
		return errInfo
	}
	//메시지 생성
	//게임턴 계산
	gameRoomSettings, errInfo := repository.FindOneFrogCurrentRound(ctx, roomID)
	if errInfo != nil {
		return errInfo
	}
	roomInfoMsg = *CreateRoomInfoMSG(ctx, preloadUsers, gameRoomSettings.CurrentRound, roomInfoMsg.ErrorInfo)
	message, err := CreateMessage(&roomInfoMsg)
	if err != nil {
		return CreateErrorMessage(_errors.ErrCodeInternal, err.Error(), _errors.ErrGameTerminated)
	}
	msg.Message = message
	msg.SessionID = ""
	sendMessageToClients(msg.RoomID, msg)
	return nil
}
