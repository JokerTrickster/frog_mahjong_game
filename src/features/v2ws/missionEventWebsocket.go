package v2ws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"main/features/v2ws/model/entity"
	"main/features/v2ws/model/request"
	"main/features/v2ws/repository"
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func MissionEventWebsocket(msg *entity.WSMessage) {
	//유저 상태를 변경한다. (대기실로 이동)
	ctx := context.Background()
	uID := msg.UserID
	roomID := msg.RoomID
	decryptedMessage, err := utils.DecryptAES(msg.Message)
	if err != nil {
		log.Fatalf("AES 복호화 에러: %s", err)
	}
	//string to struct
	req := request.ReqV2WSMissionEvent{}
	err = json.Unmarshal([]byte(decryptedMessage), &req)
	if err != nil {
		log.Fatalf("JSON 언마샬링 에러: %s", err)
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

	err = mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		// 미션 정보 생성한다.
		for _, missionID := range req.MissionIDs {
			missionEntity.MissionID = missionID

			userMissionDTO := CreateUserMissionDTO(missionEntity)
			userMissionID, err := repository.MissionCreateUserMission(ctx, tx, userMissionDTO)
			if err != nil {
				roomInfoMsg.ErrorInfo = err
				ErrorHandling(msg, &roomInfoMsg)
				return fmt.Errorf("%s", err.Msg)
			}
			userMissionCardDTO := CreateUserMissionCardDTO(missionEntity, int(userMissionID))
			err = repository.MissionCreateUserMissionCard(ctx, tx, userMissionCardDTO)
			if err != nil {
				roomInfoMsg.ErrorInfo = err
				ErrorHandling(msg, &roomInfoMsg)
				return fmt.Errorf("%s", err.Msg)
			}
		}
		// 카드 정보 체크 (소유하고 있는지 체크)
		err := repository.MissionFindAllCards(ctx, tx, &missionEntity)
		if err != nil {
			roomInfoMsg.ErrorInfo = err
			ErrorHandling(msg, &roomInfoMsg)
			return fmt.Errorf("%s", err.Msg)
		}

		preloadUsers, err = repository.MissionFindAllRoomUsers(ctx, tx, roomID)
		if err != nil {
			roomInfoMsg.ErrorInfo = err
			ErrorHandling(msg, &roomInfoMsg)
			return fmt.Errorf("%s", err.Msg)
		}

		return nil
	})
	if err != nil {
		return
	}

	// 메시지 생성
	roomInfoMsg = *DiscardCreateRoomInfoMSG(ctx, preloadUsers, 1, roomInfoMsg.ErrorInfo, 0)

	// 구조체를 JSON 문자열로 변환 (마샬링)
	message, err := CreateMessage(&roomInfoMsg)
	if err != nil {
		fmt.Println(err)
	}
	msg.Message = message
	sendMessageToClients(roomID, msg)

}
