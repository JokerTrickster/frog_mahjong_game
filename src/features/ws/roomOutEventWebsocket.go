package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"main/features/ws/model/entity"
	_errors "main/features/ws/model/errors"
	"main/features/ws/model/request"
	"main/features/ws/repository"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func RoomOutEventWebsocket(msg *entity.WSMessage) {
	//유저 상태를 변경한다. (대기실로 이동)
	ctx := context.Background()
	uID := msg.UserID
	roomID := msg.RoomID

	//string to struct
	req := request.ReqWSRoomOutEvent{}
	err := json.Unmarshal([]byte(msg.Message), &req)
	if err != nil {
		log.Printf("JSON 언마샬링 에러: %s", err)
	}

	// 비즈니스 로직
	roomInfoMsg := entity.RoomInfo{}
	preloadUsers := []entity.RoomUsers{}
	err = mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		// 요청한 유저가 방장인지 체크
		// 방장이 게임 시작 요청했는지 체크
		err := repository.RoomOutCheckOwner(ctx, tx, uID, roomID)
		if err != nil {
			return err
		}

		// 타겟 유저 데이터 변경 (플레이 상태, 룸ID)
		err = repository.RoomOutUpdateUser(ctx, tx, uint(req.TargetUserID), roomID)
		if err != nil {
			return err
		}

		// 룸 유저 정보 삭제
		err = repository.RoomOutDeleteRoomUser(ctx, tx, uint(req.TargetUserID), roomID)
		if err != nil {
			return err
		}

		// 방 현재 인원을 감소시킨다.
		err = repository.RoomOutUpdateRoom(ctx, tx, roomID)
		if err != nil {
			return err
		}
		preloadUsers, err = repository.RoomOutFindAllRoomUsers(ctx, tx, roomID)
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
	// 메시지 생성
	roomInfoMsg = *CreateRoomInfoMSG(ctx, preloadUsers, 1, roomInfoMsg.ErrorInfo)
	// 구조체를 JSON 문자열로 변환 (마샬링)
	message, err := CreateMessage(&roomInfoMsg)
	if err != nil {
		fmt.Println(err)
	}
	msg.Message = message
	sendMessageToClients(roomID, msg)
}
