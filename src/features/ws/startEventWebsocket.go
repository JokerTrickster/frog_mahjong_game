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

func StartEventWebsocket(msg *entity.WSMessage) {
	//유저 상태를 변경한다. (대기실로 이동)
	ctx := context.Background()
	uID := msg.UserID
	roomID := msg.RoomID

	// 비즈니스 로직
	roomInfoMsg := entity.RoomInfo{}
	preloadUsers := []entity.RoomUsers{}
	err := mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {

		// 방장이 게임 시작 요청했는지 체크
		err := repository.StartCheckOwner(ctx, tx, uID, roomID)
		if err != nil {
			return err
		}

		roomUsers, err := repository.StartFindRoomUsers(ctx, tx, roomID)
		if err != nil {
			return err
		}
		// room user 데이터 변경 (플레이 순번 랜덤으로 생성)
		updatedRoomUsers, err := StartUpdateRoomUsers(roomUsers)
		if err != nil {
			return err
		}
		// room user 데이터 변경 (플레이 순번 랜덤으로 생성)
		err = repository.StartUpdateRoomUser(ctx, tx, updatedRoomUsers)
		if err != nil {
			return err
		}

		// room 데이터 상태 변경 (대기 -> 플레이)
		err = repository.StartUpdateRoom(ctx, tx, roomID, "play")
		if err != nil {
			return err
		}
		// 카드 정보 가져온다.
		cards, err := repository.StartFindCards(ctx, tx)
		if err != nil {
			return err
		}
		// cards 데이터 생성
		userCards := CreateInitCards(roomID, cards)
		err = repository.StartCreateCards(ctx, tx, userCards)
		if err != nil {
			return err
		}
		preloadUsers, err = repository.StartFindAllRoomUsers(ctx, tx, roomID)
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

	//유저 상태를 변경한다. (방에 참여)
	if clients, ok := entity.WSClients[msg.RoomID]; ok {
		//에러 발생시 이벤트 요청한 유저에게만 메시지를 전달한다.
		if roomInfoMsg.ErrorInfo != nil || err != nil {
			for client := range clients {
				if clients[client].UserID == msg.UserID {
					err := client.WriteJSON(msg)
					if err != nil {
						fmt.Printf("error: %v", err)
						client.Close()
						delete(clients, client)
					}
				}
			}
		} else {
			for client := range clients {
				err := client.WriteJSON(msg)
				if err != nil {
					fmt.Printf("error: %v", err)
					client.Close()
					delete(clients, client)
				}
			}
		}
	}
}
