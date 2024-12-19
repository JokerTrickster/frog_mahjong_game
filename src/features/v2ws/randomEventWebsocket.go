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

func RandomEventWebsocket(msg *entity.WSMessage) {
	//유저 상태를 변경한다. (대기실로 이동)
	ctx := context.Background()
	uID := msg.UserID
	roomID := msg.RoomID

	//string to struct
	req := request.ReqWSRandom{}
	err := json.Unmarshal([]byte(msg.Message), &req)
	if err != nil {
		log.Fatalf("JSON 언마샬링 에러: %s", err)
	}

	RandomEntity := entity.WSRandomEntity{
		RoomID: roomID,
		UserID: uID,
		Count:  req.Count,
	}

	// 비즈니스 로직
	roomInfoMsg := entity.RoomInfo{}
	preloadUsers := []entity.RoomUsers{}
	err = mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		// none 카드 중 count만큼 랜덤으로 owned로 변경한다.
		err := repository.RandomUpdateRandomCards(ctx, tx, &RandomEntity)
		if err != nil {
			return err
		}

		// 소유 카드 수 업데이트
		// 유저id로 room_users에서 찾아서 card_count를 더한 후 업데이트 한다.
		err = repository.RandomUpdateRoomUserCardCount(ctx, tx, &RandomEntity)
		if err != nil {
			return err
		}

		// 현재 참여하고 있는 유저에 대한 정보를 가져와서 메시지 전달한다.
		preloadUsers, err = repository.RandomFindAllRoomUsers(ctx, tx, roomID)
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
	// 유저 상태를 변경한다. (방에 참여)
	if sessionIDs, ok := entity.RoomSessions[msg.RoomID]; ok {
		// 메시지 생성
		roomInfoMsg := *CreateRoomInfoMSG(ctx, preloadUsers, 1, roomInfoMsg.ErrorInfo, 1)

		if roomInfoMsg.GameInfo.AllPicked {
			err := mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
				// 카드 상태를 picked -> owned로 변경
				err := repository.RandomUpdateAllCardState(ctx, tx, msg.RoomID)
				if err != nil {
					fmt.Printf("Error updating card state: %v\n", err)
					return err
				}

				// 오픈 카드가 비어 있다면 새로운 카드를 오픈
				err = repository.RandomUpdateOpenCards(ctx, tx, msg.RoomID)
				if err != nil {
					fmt.Printf("Error updating open cards: %v\n", err)
					return err
				}

				return nil
			})

			// 트랜잭션 에러 처리
			if err != nil {
				fmt.Printf("Transaction error: %v\n", err)
				return
			}

			// 오픈 카드 정보를 가져옴
			openCards, err := repository.FindAllOpenCards(ctx, int(msg.RoomID))
			if err != nil {
				fmt.Printf("Error fetching open cards: %v\n", err)
				return
			}
			roomInfoMsg.GameInfo.OpenCards = openCards
		}

		// 방의 모든 유저에게 메시지 전달
		for _, sessionID := range sessionIDs {
			if client, exists := entity.WSClients[sessionID]; exists {
				filterRoomInfoMsg := Deepcopy(roomInfoMsg)

				// 구조체를 JSON 문자열로 변환 (마샬링)
				message, err := CreateMessage(&filterRoomInfoMsg)
				if err != nil {
					fmt.Printf("Error creating message: %v\n", err)
					continue
				}
				msg.Message = message

				// 메시지 전송
				err = client.Conn.WriteJSON(msg)
				if err != nil {
					fmt.Printf("Error sending message to user %d: %v\n", client.UserID, err)
					client.Close()

					// 클라이언트를 종료 및 정리
					delete(entity.WSClients, sessionID)
					removeSessionFromRoom(client.RoomID, sessionID)
				}
			}
		}

		// 방이 비어 있으면 삭제
		if len(entity.RoomSessions[msg.RoomID]) == 0 {
			delete(entity.RoomSessions, msg.RoomID)
		}
	}

}
