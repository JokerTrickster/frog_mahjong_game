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

func RandomEventWebsocket(msg *entity.WSMessage) *entity.ErrorInfo {
	//유저 상태를 변경한다. (대기실로 이동)
	ctx := context.Background()
	uID := msg.UserID
	roomID := msg.RoomID
	decryptedMessage, err := utils.DecryptAES(msg.Message)
	if err != nil {
		return CreateErrorMessage(_errors.ErrCodeBadRequest, _errors.ErrCryptoFailed, "AES 복호화 에러")
	}
	//string to struct
	req := request.ReqWSRandom{}
	err = json.Unmarshal([]byte(decryptedMessage), &req)
	if err != nil {
		return CreateErrorMessage(_errors.ErrCodeBadRequest, _errors.ErrUnmarshalFailed, "JSON 언마샬링 에러")
	}

	RandomEntity := entity.WSRandomEntity{
		RoomID: roomID,
		UserID: uID,
		Count:  req.Count,
	}

	// 비즈니스 로직
	// TODO 트랭잭션 처리
	roomInfoMsg := entity.RoomInfo{}
	preloadUsers := []entity.RoomUsers{}
	var errInfo *entity.ErrorInfo
	err = mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		// none 카드 중 count만큼 랜덤으로 owned로 변경한다.
		errInfo = repository.RandomUpdateRandomCards(ctx, tx, &RandomEntity)
		if err != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}

		// 소유 카드 수 업데이트
		// 유저id로 room_users에서 찾아서 card_count를 더한 후 업데이트 한다.
		errInfo = repository.RandomUpdateRoomUserCardCount(ctx, tx, &RandomEntity)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}

		// 현재 참여하고 있는 유저에 대한 정보를 가져와서 메시지 전달한다.
		preloadUsers, errInfo = repository.RandomFindAllRoomUsers(ctx, tx, roomID)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}
		return nil
	})
	if err != nil {
		return errInfo
	}

	// 메시지 생성
	roomInfoMsg = *CreateRoomInfoMSG(ctx, preloadUsers, 1, roomInfoMsg.ErrorInfo, 1)

	if roomInfoMsg.GameInfo.AllPicked {
		err := mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
			// 카드 상태를 picked -> owned로 변경
			errInfo = repository.RandomUpdateAllCardState(ctx, tx, msg.RoomID)
			if err != nil {
				return fmt.Errorf("%s", errInfo.Msg)
			}

			// 오픈 카드가 비어 있다면 새로운 카드를 오픈
			errInfo = repository.RandomUpdateOpenCards(ctx, tx, msg.RoomID)
			if err != nil {
				return fmt.Errorf("%s", errInfo.Msg)
			}
			return nil
		})

		// 트랜잭션 에러 처리
		if err != nil {
			return errInfo
		}

		// 오픈 카드 정보를 가져옴
		openCards, errInfo := repository.FindAllOpenCards(ctx, int(msg.RoomID))
		if err != nil {
			fmt.Printf("Error fetching open cards: %v\n", err)
			return errInfo
		}
		roomInfoMsg.GameInfo.OpenCards = openCards
	}

	message, err := CreateMessage(&roomInfoMsg)
	if err != nil {
		return CreateErrorMessage(_errors.ErrCodeInternal, err.Error(), _errors.ErrGameTerminated)
	}
	msg.Message = message
	sendMessageToClients(roomID, msg)
	return nil
}
