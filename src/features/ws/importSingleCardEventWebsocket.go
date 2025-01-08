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

func ImportSingleCardEventWebsocket(msg *entity.WSMessage) *entity.ErrorInfo {
	//유저 상태를 변경한다. (대기실로 이동)
	ctx := context.Background()
	uID := msg.UserID
	roomID := msg.RoomID

	//string to struct
	req := request.ReqWSImportSingleCard{}
	err := json.Unmarshal([]byte(msg.Message), &req)
	if err != nil {
		return CreateErrorMessage(_errors.ErrCodeBadRequest, _errors.ErrUnmarshalFailed, "JSON 언마샬링 에러")
	}

	importSingleCardEntity := entity.WSImportSingleCardEntity{
		RoomID: roomID,
		UserID: uID,
		Cards: &mysql.FrogUserCards{
			CardID: int(req.CardID),
			RoomID: int(roomID),
			UserID: int(uID),
		},
	}

	// 비즈니스 로직
	roomInfoMsg := entity.RoomInfo{}
	preloadUsers := []entity.RoomUsers{}
	var errInfo *entity.ErrorInfo
	err = mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		// 카드 상태 없데이트
		errInfo := repository.ImportSingleCardUpdateCardState(ctx, tx, &importSingleCardEntity)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}
		// 소유 카드 수 업데이트
		// 유저id로 room_users에서 찾아서 card_count를 더한 후 업데이트 한다.
		errInfo = repository.ImportSingleCardUpdateRoomUserCardCount(ctx, tx, &importSingleCardEntity)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}

		// 현재 참여하고 있는 유저에 대한 정보를 가져와서 메시지 전달한다.
		preloadUsers, errInfo = repository.PreloadFindGameInfo(ctx, tx, roomID)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}
		return nil
	})
	if err != nil {
		return errInfo
	}

	// 메시지 생성
	roomInfoMsg = *CreateRoomInfoMSG(ctx, preloadUsers, req.PlayTurn, roomInfoMsg.ErrorInfo)

	message, err := CreateMessage(&roomInfoMsg)
	if err != nil {
		return CreateErrorMessage(_errors.ErrCodeInternal, _errors.ErrMarshalFailed, "메시지 생성 실패")
	}
	msg.Message = message
	msg.SessionID = ""
	sendMessageToClients(roomID, msg)
	return nil
}
