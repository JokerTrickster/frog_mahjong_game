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

func ImportSingleCardEventWebsocket(msg *entity.WSMessage) *entity.ErrorInfo {
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
	req := request.ReqWSImportSingleCard{}
	err = json.Unmarshal([]byte(decryptedMessage), &req)
	if err != nil {
		return CreateErrorMessage(_errors.ErrCodeBadRequest, _errors.ErrUnmarshalFailed, "JSON 언마샬링 에러")
	}

	importSingleCardEntity := entity.WSImportSingleCardEntity{
		RoomID: roomID,
		UserID: uID,
		CardID: uint(req.CardID),
	}

	// 비즈니스 로직
	roomInfoMsg := entity.RoomInfo{}
	preloadUsers := []entity.RoomUsers{}
	var errInfo *entity.ErrorInfo
	// 카드수가 4장 미만인지 체크
	cardCount, errInfo := repository.ImportSingleCardOwnerCardCount(ctx, roomID, uID)
	if errInfo != nil {
		return errInfo
	}
	if cardCount > 3 {
		// 해당 이벤트를 처리하지 않는다.
		return CreateErrorMessage(_errors.ErrCodeBadRequest, "카드는 4장만 소유할 수 있습니다.", _errors.ErrBadRequest)
	}
	// 카드가 이미 선택되었는지 체크
	errInfo = repository.ImportSingleCardFindOneCard(ctx, roomID, importSingleCardEntity.CardID)
	if errInfo != nil {
		return errInfo
	}

	err = mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		// user_bird_cards 테이블에 카드 상태 없데이트
		userBirdCardDTO := CreateUserBirdCardDTO(importSingleCardEntity)
		errInfo := repository.ImportSingleCardCreateCard(ctx, tx, userBirdCardDTO)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}
		// 소유 카드 수 업데이트
		// 유저id로 room_users 테이블에서 찾아서 card_count를 더한 후 업데이트 한다.
		errInfo = repository.ImportSingleCardUpdateRoomUserCardCount(ctx, tx, &importSingleCardEntity)
		if err != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}

		// 현재 참여하고 있는 유저에 대한 정보를 가져와서 메시지 전달한다.
		preloadUsers, errInfo = repository.ImportSingleCardFindAllRoomUsers(ctx, tx, roomID)
		if err != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}
		return nil
	})
	if err != nil {
		return errInfo
	}
	// TODO 트랜잭션 정리 필요
	// 유저 상태를 변경한다. (방에 참여)
	// 메시지 생성
	roomInfoMsg = *CreateRoomInfoMSG(ctx, preloadUsers, req.PlayTurn, roomInfoMsg.ErrorInfo, int(req.CardID))
	// 모든 유저가 카드를 선택했을 경우 처리
	if roomInfoMsg.GameInfo.AllPicked {
		// 카드 상태를 picked -> owned로 변경
		err := mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
			errInfo = repository.ImportSingleCardUpdateAllCardState(ctx, tx, msg.RoomID)
			if errInfo != nil {
				return fmt.Errorf("%s", errInfo.Msg)
			}
			return nil
		})

		if err != nil {
			return errInfo
		}

		// 오픈 카드 업데이트
		err = mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
			errInfo = repository.ImportSingleCardUpdateOpenCards(ctx, tx, msg.RoomID)
			if errInfo != nil {
				return fmt.Errorf("%s", errInfo.Msg)
			}
			return nil
		})

		if err != nil {
			return errInfo
		}
	}

	// 오픈 카드 정보를 가져온다.
	openCards, errInfo := repository.FindAllOpenCards(ctx, int(msg.RoomID))
	if errInfo != nil {
		return errInfo
	}
	roomInfoMsg.GameInfo.OpenCards = openCards

	message, err := CreateMessage(&roomInfoMsg)
	if err != nil {
		return CreateErrorMessage(_errors.ErrCodeBadRequest, _errors.ErrMarshalFailed, "메시지 생성 에러")
	}
	msg.Message = message
	msg.SessionID = ""
	sendMessageToClients(roomID, msg)
	return nil
}
