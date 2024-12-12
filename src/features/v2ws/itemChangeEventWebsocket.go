package v2ws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"main/features/v2ws/model/entity"
	_errors "main/features/v2ws/model/errors"
	"main/features/v2ws/model/request"
	"main/features/v2ws/repository"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func ItemChangeEventWebsocket(msg *entity.WSMessage) {
	//유저 상태를 변경한다. (대기실로 이동)
	ctx := context.Background()
	uID := msg.UserID
	roomID := msg.RoomID

	//string to struct
	req := request.ReqWSItemChange{}
	err := json.Unmarshal([]byte(msg.Message), &req)
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
		ErrorHandling(msg, roomID, uID, &roomInfoMsg)
		return
	}

	err = mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		// 아이템 사용 (오픈된 카드를 교체한다.)
		err = repository.ItemChange(ctx, tx, ItemChangeEntity)
		if err != nil {
			return err
		}
		// 아이템 사용 횟수를 -1 감소한다.
		err = repository.ItemChangeConsumeUserItems(ctx, tx, ItemChangeEntity)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		roomInfoMsg.ErrorInfo = &entity.ErrorInfo{
			Code: 500,
			Msg:  err.Error(),
			Type: _errors.ErrInternalServer,
		}
	}
	// 현재 참여하고 있는 유저에 대한 정보를 가져와서 메시지 전달한다.
	preloadUsers, err = repository.ItemChangeFindAllRoomUsers(ctx, roomID)
	if err != nil {
		roomInfoMsg.ErrorInfo = &entity.ErrorInfo{
			Code: 500,
			Msg:  err.Error(),
			Type: _errors.ErrInternalServer,
		}
		ErrorHandling(msg, roomID, uID, &roomInfoMsg)
		return
	}
	//유저 상태를 변경한다. (방에 참여)
	if clients, ok := entity.WSClients[msg.RoomID]; ok {
		// 메시지 생성
		roomInfoMsg = *CreateRoomInfoMSG(ctx, preloadUsers, 1, roomInfoMsg.ErrorInfo, 1)

		for client := range clients {
			filterRoomInfoMsg := Deepcopy(roomInfoMsg)

			// 구조체를 JSON 문자열로 변환 (마샬링)
			message, err := CreateMessage(&filterRoomInfoMsg)
			if err != nil {
				fmt.Println(err)
			}
			msg.Message = message
			err = client.WriteJSON(msg)
			if err != nil {
				fmt.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
