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

func ItemChangeEventWebsocket(msg *entity.WSMessage) {
	//유저 상태를 변경한다. (대기실로 이동)
	ctx := context.Background()
	uID := msg.UserID
	roomID := msg.RoomID
	// 복호화 후 JSON 언마샬링
	decryptedMessage, err := utils.DecryptAES(msg.Message)
	if err != nil {
		log.Fatalf("AES 복호화 에러: %s", err)
	}
	//string to struct
	req := request.ReqWSItemChange{}
	err = json.Unmarshal([]byte(decryptedMessage), &req)
	if err != nil {
		log.Fatalf("JSON 언마샬링 에러: %s", err)
	}

	ItemChangeEntity := entity.WSItemChangeEntity{
		RoomID: roomID,
		UserID: uID,
		ItemID: req.ItemID,
	}
	// 비즈니스 로직
	roomInfoMsg := entity.RoomInfo{}
	preloadUsers := []entity.RoomUsers{}

	// 아이템 사용 가능한지 체크
	newErr := repository.ItemChangeCheck(ctx, ItemChangeEntity)
	if newErr != nil {
		roomInfoMsg.ErrorInfo = newErr
		ErrorHandling(msg, &roomInfoMsg)
		return
	}

	err = mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		// 아이템 사용 (오픈된 카드를 교체한다.)
		err := repository.ItemChange(ctx, tx, ItemChangeEntity)
		if err != nil {
			roomInfoMsg.ErrorInfo = err
			ErrorHandling(msg, &roomInfoMsg)
			return fmt.Errorf("%s", err.Msg)
		}
		// 아이템 사용 횟수를 -1 감소한다.
		err = repository.ItemChangeConsumeUserItems(ctx, tx, ItemChangeEntity)
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
	// 현재 참여하고 있는 유저에 대한 정보를 가져와서 메시지 전달한다.
	preloadUsers, newEer := repository.ItemChangeFindAllRoomUsers(ctx, roomID)
	if newEer != nil {
		return
	}
	roomInfoMsg = *CreateRoomInfoMSG(ctx, preloadUsers, 1, roomInfoMsg.ErrorInfo, 0)

	message, err := CreateMessage(&roomInfoMsg)
	if err != nil {
		fmt.Println(err)
		return
	}
	msg.Message = message
	sendMessageToClients(roomID, msg)
}
