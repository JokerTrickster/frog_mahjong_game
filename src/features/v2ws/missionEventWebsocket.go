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

func MissionEventWebsocket(msg *entity.WSMessage) *entity.ErrorInfo {
	//유저 상태를 변경한다. (대기실로 이동)
	ctx := context.Background()
	uID := msg.UserID
	roomID := msg.RoomID
	decryptedMessage, err := utils.DecryptAES(msg.Message)
	if err != nil {
		return CreateErrorMessage(_errors.ErrCodeBadRequest, _errors.ErrCryptoFailed, "AES 복호화 에러")

	}
	//string to struct
	req := request.ReqV2WSMissionEvent{}
	err = json.Unmarshal([]byte(decryptedMessage), &req)
	if err != nil {
		return CreateErrorMessage(_errors.ErrCodeBadRequest, _errors.ErrUnmarshalFailed, "JSON 언마샬링 에러")
	}

	missionEntity := entity.V2WSMissionEntity{
		RoomID: roomID,
		UserID: uID,
	}
	for _, cardID := range req.Cards {
		missionEntity.Cards = append(missionEntity.Cards, cardID)
	}

	// 비즈니스 로직
	roomInfoMsg := entity.RoomInfo{}
	preloadUsers := []entity.RoomUsers{}
	var errInfo *entity.ErrorInfo
	err = mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		// 미션 정보 생성한다.
		for _, missionID := range req.MissionIDs {
			missionEntity.MissionID = missionID

			userMissionDTO := CreateUserMissionDTO(missionEntity)
			userMissionID, errInfo := repository.MissionCreateUserMission(ctx, tx, userMissionDTO)
			if errInfo != nil {
				return fmt.Errorf("%s", errInfo.Msg)
			}
			userMissionCardDTO := CreateUserMissionCardDTO(missionEntity, int(userMissionID))
			errInfo = repository.MissionCreateUserMissionCard(ctx, tx, userMissionCardDTO)
			if err != nil {
				return fmt.Errorf("%s", errInfo.Msg)
			}
		}
		// 카드 정보 체크 (소유하고 있는지 체크)
		errInfo := repository.MissionFindAllCards(ctx, tx, &missionEntity)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}

		preloadUsers, errInfo = repository.MissionFindAllRoomUsers(ctx, tx, roomID)
		if errInfo != nil {
			return fmt.Errorf("%s", errInfo.Msg)
		}

		return nil
	})
	if err != nil {
		return errInfo
	}

	// 메시지 생성
	roomInfoMsg = *DiscardCreateRoomInfoMSG(ctx, preloadUsers, 1, roomInfoMsg.ErrorInfo, 0)

	// 구조체를 JSON 문자열로 변환 (마샬링)
	message, err := CreateMessage(&roomInfoMsg)
	if err != nil {
		return CreateErrorMessage(_errors.ErrCodeBadRequest, _errors.ErrMarshalFailed, "메시지 생성 에러")
	}
	msg.Message = message
	sendMessageToClients(roomID, msg)
	return nil
}
