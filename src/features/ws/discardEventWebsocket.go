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

func DiscardCardsEventWebsocket(msg *entity.WSMessage) {
	//유저 상태를 변경한다. (대기실로 이동)
	ctx := context.Background()
	uID := msg.UserID
	roomID := msg.RoomID

	//string to struct
	req := request.ReqWSDiscardCards{}
	err := json.Unmarshal([]byte(msg.Message), &req)
	if err != nil {
		SendErrorMessage(msg, CreateErrorMessage(_errors.ErrCodeBadRequest, _errors.ErrUnmarshalFailed, "JSON 언마샬링 에러"))
		return
	}
	DiscardCardsEntity := entity.WSDiscardCardsEntity{
		RoomID: roomID,
		UserID: uID,
		CardID: uint(req.CardID),
	}

	// 비즈니스 로직
	roomInfoMsg := entity.RoomInfo{}
	preloadUsers := []entity.RoomUsers{}
	err = mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		// 카드 상태 없데이트
		newErr := repository.DiscardCardsUpdateCardState(ctx, tx, &DiscardCardsEntity)
		if newErr != nil {
			SendErrorMessage(msg, newErr)
			return fmt.Errorf("%s", newErr.Msg)
		}
		// 소유 카드 수 업데이트
		// 유저id로 room_users에서 찾아서 card_count를 뺀 후 업데이트 한다.
		newErr = repository.DiscardCardsUpdateRoomUserCardCount(ctx, tx, &DiscardCardsEntity)
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
	//게임턴 계산
	playTurn := CalcPlayTurn(req.PlayTurn, len(entity.RoomSessions[msg.RoomID]))
	roomInfoMsg = *CreateRoomInfoMSG(ctx, preloadUsers, playTurn, roomInfoMsg.ErrorInfo)

	// 론 가능 여부를 true로 변경
	roomInfoMsg.GameInfo.IsLoanAllowed = true
	// 구조체를 JSON 문자열로 변환 (마샬링)
	message, err := CreateMessage(&roomInfoMsg)
	if err != nil {
		fmt.Println(err)
	}
	msg.Message = message
	sendMessageToClients(msg.RoomID, msg)
}
