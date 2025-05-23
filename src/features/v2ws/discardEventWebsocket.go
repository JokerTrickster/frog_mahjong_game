package v2ws

import (
	"context"
	"encoding/json"
	"fmt"
	"main/features/v2ws/model/entity"
	_errors "main/features/v2ws/model/errors"
	"main/features/v2ws/model/request"
	"main/features/v2ws/repository"
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func DiscardCardsEventWebsocket(msg *entity.WSMessage) *entity.ErrorInfo {
	//유저 상태를 변경한다. (대기실로 이동)
	ctx := context.Background()
	uID := msg.UserID
	roomID := msg.RoomID
	// 복호화 후 JSON 언마샬링
	decryptedMessage, err := utils.DecryptAES(msg.Message)
	if err != nil {
		return CreateErrorMessage(_errors.ErrCodeBadRequest, _errors.ErrCryptoFailed, "AES 복호화 에러")
	}
	//string to struct
	req := request.ReqWSDiscardCards{}
	err = json.Unmarshal([]byte(decryptedMessage), &req)
	if err != nil {
		return CreateErrorMessage(_errors.ErrCodeBadRequest, _errors.ErrUnmarshalFailed, "JSON 언마샬링 에러")
	}
	DiscardCardsEntity := entity.WSDiscardCardsEntity{
		RoomID: roomID,
		UserID: uID,
		CardID: uint(req.CardID),
	}

	// 비즈니스 로직
	// 보유 카드수가 4장인지 체크
	roomInfoMsg := entity.RoomInfo{}
	preloadUsers := []entity.RoomUsers{}
	var errInfo *entity.ErrorInfo
	cardCount, errInfo := repository.DiscardCardsOwnerCardCount(ctx, roomID, uID)
	if errInfo != nil {
		return errInfo
	}
	if cardCount != 4 {
		return CreateErrorMessage(_errors.ErrCodeBadRequest, "카드를 4장 선택해주세요.", _errors.ErrBadRequest)
	}

	err = mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		// 카드 상태 없데이트
		errInfo = repository.DiscardCardsUpdateCardState(ctx, tx, &DiscardCardsEntity)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}
		// 소유 카드 수 업데이트
		// 유저id로 room_users에서 찾아서 card_count를 뺀 후 업데이트 한다.
		errInfo = repository.DiscardCardsUpdateRoomUserCardCount(ctx, tx, &DiscardCardsEntity)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}

		preloadUsers, errInfo = repository.DiscardCardsFindAllRoomUsers(ctx, tx, roomID)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}
		return nil
	})
	if err != nil {
		return errInfo
	}
	// 유저 상태를 변경한다. (방에 참여)
	if sessionIDs, ok := entity.RoomSessions[msg.RoomID]; ok {
		// 게임 턴 계산
		playTurn := CalcPlayTurn(req.PlayTurn, len(sessionIDs))
		roomInfoMsg := *DiscardCreateRoomInfoMSG(ctx, preloadUsers, playTurn, roomInfoMsg.ErrorInfo, int(req.CardID))

		// 모든 유저가 카드를 선택했을 때
		if roomInfoMsg.GameInfo.AllPicked {
			// 카드 상태를 picked -> owned로 변경
			err := mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
				errInfo := repository.DiscardCardUpdateAllCardState(ctx, tx, msg.RoomID)
				if err != nil {
					return fmt.Errorf("%s", errInfo.Msg)
				}
				preloadUsers, errInfo = repository.DiscardCardsFindAllRoomUsers(ctx, tx, msg.RoomID)
				if err != nil {
					return fmt.Errorf("%s", errInfo.Msg)
				}
				return nil
			})

			// 에러 처리
			if err != nil {
				return errInfo
			}

			// 게임 상태 갱신
			roomInfoMsg = *DiscardCreateRoomInfoMSG(ctx, preloadUsers, playTurn, roomInfoMsg.ErrorInfo, int(req.CardID))
			roomInfoMsg.GameInfo.AllPicked = true
		}

		message, err := CreateMessage(&roomInfoMsg)
		if err != nil {
			return CreateErrorMessage(_errors.ErrCodeInternal, err.Error(), _errors.ErrGameTerminated)
		}
		msg.Message = message
		sendMessageToClients(roomID, msg)
	}
	return nil
}
