package ws

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"main/features/ws/model/entity"
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
	err := mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {

		// 방장이 게임 시작 요청했는지 체크
		err := repository.StartCheckOwner(ctx, tx, uID, roomID)
		if err != nil {
			return err
		}

		// 방에 있는 유저들이 모두 레디 상태인지 확인
		roomUsers, err := repository.StartCheckReady(ctx, tx, roomID)
		if err != nil {
			return err
		}
		if allReady := CheckRoomUsersReady(roomUsers); !allReady {
			return errors.New("모든 유저가 준비하지 않았습니다.")
		}

		// room user 데이터 변경 (대기 -> 플레이, 플레이 순번 랜덤으로 생성)
		updatedRoomUsers, err := StartUpdateRoomUsers(roomUsers)
		if err != nil {
			return err
		}
		// room user 데이터 변경 (대기 -> 플레이)
		err = repository.StartUpdateRoomUser(ctx, tx, updatedRoomUsers)
		if err != nil {
			return err
		}

		// room 데이터 상태 변경 (대기 -> 플레이)
		err = repository.StartUpdateRoom(ctx, tx, roomID, "play")
		if err != nil {
			return err
		}

		// cards 데이터 생성
		cards := CreateInitCards(roomID)
		err = repository.StartCreateCards(ctx, tx, roomID, cards)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		roomInfoMsg.ErrorInfo = &entity.ErrorInfo{
			Code: 500,
			Msg:  err.Error(),
		}
	}

	// 메시지 생성
	// 현재 참여하고 있는 유저에 대한 정보를 가져와서 메시지 전달한다.
	preloadUsers, err := repository.ReadyFindAllRoomUsers(ctx, uint(msg.RoomID))
	if err != nil {
		log.Println(err)
	}
	//유저 정보 저장
	for _, roomUser := range preloadUsers {
		user := entity.User{
			ID:          uint(roomUser.UserID),
			PlayerState: roomUser.PlayerState,
			Coin:        roomUser.User.Coin,
			Name:        roomUser.User.Name,
			Email:       roomUser.User.Email,
			TurnNumber:  roomUser.TurnNumber,
		}
		if roomUser.Room.OwnerID == roomUser.UserID {
			user.IsOwner = true
		}
		roomInfoMsg.Users = append(roomInfoMsg.Users, &user)
	}
	//게임 정보 저장
	gameInfo := entity.GameInfo{
		PlayTurn: 1,
		AllReady: true,
	}
	roomInfoMsg.GameInfo = &gameInfo

	// 구조체를 JSON 문자열로 변환 (마샬링)
	jsonData, err := json.Marshal(roomInfoMsg)
	if err != nil {
		log.Fatalf("JSON 마샬링 에러: %s", err)
	}

	// JSON 바이트 배열을 문자열로 변환
	jsonString := string(jsonData)
	msg.Message = jsonString

	//유저 상태를 변경한다. (방에 참여)
	if clients, ok := entity.WSClients[msg.RoomID]; ok {
		//에러 발생시 이벤트 요청한 유저에게만 메시지를 전달한다.
		if err != nil {
			for client := range clients {
				if clients[client].UserID == msg.UserID {
					err := client.WriteJSON(msg)
					if err != nil {
						log.Printf("error: %v", err)
						client.Close()
						delete(clients, client)
					}
				}
			}
		} else {
			for client := range clients {
				err := client.WriteJSON(msg)
				if err != nil {
					log.Printf("error: %v", err)
					client.Close()
					delete(clients, client)
				}
			}
		}
	}
}
