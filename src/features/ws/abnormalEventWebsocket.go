package ws

import (
	"context"
	"fmt"
	"main/features/ws/model/entity"
	_errors "main/features/ws/model/errors"
	"main/features/ws/repository"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func AbnormalErrorHandling(roomID, userID uint) {
	// 비정상적인 에러 발생했으므로 비정상적 에러 처리하는 로직 실행

	//business logic

	roomInfoMsg := entity.RoomInfo{}
	preloadUsers := []entity.RoomUsers{}
	msg := entity.WSMessage{
		RoomID: roomID,
		UserID: userID,
	}
	ctx := context.TODO()
	err := mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		//모든 유저 게임 종료 처리하고 대기 상태로 변경한다.
		abnormalEntity := entity.WSAbnormalEntity{
			RoomID:         roomID,
			AbnormalUserID: userID,
		}
		// 비정상적인 유저 삭제처리

		// 카드 정보 모두 삭제
		err := repository.AbnormalDeleteAllCards(ctx, tx, &abnormalEntity)
		if err != nil {
			return err
		}
		// 방 삭제 처리
		err = repository.AbnormalDeleteRoom(ctx, tx, &abnormalEntity)
		if err != nil {
			return err
		}

		// 유저 상태 변경
		err = repository.AbnormalUpdateUsers(ctx, tx, &abnormalEntity)
		if err != nil {
			return err
		}
		preloadUsers, err = repository.AbnormalFindAllRoomUsers(ctx, tx, roomID)
		if err != nil {
			return err
		}
		// 에러 메시지에 상대방이 게임 도중 나가서 강제 종료됐다는 에러 메시지 표시한다.
		roomInfoMsg.ErrorInfo = &entity.ErrorInfo{
			Code: 500,
			Msg:  "상대방이 게임 도중 나가서 강제 종료됐습니다.",
			Type: _errors.ErrAbnormalExit,
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
	roomInfoMsg.GameInfo.AllReady = false

	// 구조체를 JSON 문자열로 변환 (마샬링)
	message, err := CreateMessage(&roomInfoMsg)
	if err != nil {
		fmt.Println(err)
	}
	msg.Message = message
	//방 유저들에게 메시지 전달
	if clients, ok := entity.WSClients[msg.RoomID]; ok {
		//에러 발생시 이벤트 요청한 유저에게만 메시지를 전달한다.
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				fmt.Printf("message send error: %v", err)
				client.Close()
				delete(clients, client)
			}
			client.Close()
			delete(clients, client)
		}
	}
	if len(entity.WSClients[msg.RoomID]) == 0 {
		delete(entity.WSClients, roomID)
	}
}
