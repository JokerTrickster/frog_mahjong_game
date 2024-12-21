package v2ws

import (
	"context"
	"fmt"
	"main/features/v2ws/model/entity"
	_errors "main/features/v2ws/model/errors"
	"main/features/v2ws/repository"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func StartEventWebsocket(msg *entity.WSMessage) {
	//유저 상태를 변경한다. (대기실로 이동)
	ctx := context.Background()
	uID := msg.UserID
	roomID := msg.RoomID

	// 비즈니스 로직
	//해당 방이 대기상태인지 체크한다.
	roomState, err := repository.StartCheckRoomState(ctx, roomID)
	if err != nil {
		fmt.Println(err)
		return
	}
	if roomState != "wait" {
		fmt.Println("게임이 시작되었습니다.")
		return
	}

	roomInfoMsg := entity.RoomInfo{}
	preloadUsers := []entity.RoomUsers{}
	err = mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		// 방장이 게임 시작 요청했는지 체크
		err := repository.StartCheckOwner(ctx, tx, uID, roomID)
		if err != nil {
			return err
		}

		// room 데이터 값 변경 (상태 변경, 시작 시간 추가)
		roomUpdateData := StartUpdateRoom(roomID)
		err = repository.StartUpdateRoom(ctx, tx, roomID, roomUpdateData)
		if err != nil {
			return err
		}

		// 유저들 코인 -1 차감한다.
		err = repository.StartDiffCoin(ctx, tx, roomID)
		if err != nil {
			return err
		}

		// 기존 카드가 있다면 모두 제거한다.
		err = repository.StartDeleteCards(ctx, tx, uID)
		if err != nil {
			return err
		}
		birdCards, err := repository.StartBirdCard(ctx, tx)
		if err != nil {
			return err
		}
		// 카드를 생성한다.
		cards := CreateInitCards(roomID, birdCards)
		err = repository.StartCreateCards(ctx, tx, cards)
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
	preloadUsers, err = repository.StartFindAllRoomUsers(ctx, roomID)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 메시지 생성
	roomInfoMsg = *CreateRoomInfoMSG(ctx, preloadUsers, 1, roomInfoMsg.ErrorInfo, 0)
	openCards, err := repository.StartUpdateCardState(ctx, roomID)
	if err != nil {
		fmt.Println(err)
	}
	roomInfoMsg.GameInfo.OpenCards = openCards

	// 구조체를 JSON 문자열로 변환 (마샬링)
	message, err := CreateMessage(&roomInfoMsg)
	if err != nil {
		fmt.Println(err)
	}
	msg.Message = message
	msg.SessionID = ""
	// 유저 상태를 변경한다. (방에 참여)
	if sessionIDs, ok := entity.RoomSessions[msg.RoomID]; ok {
		// 에러 발생 시 이벤트 요청한 유저에게만 메시지를 전달한다.
		if roomInfoMsg.ErrorInfo != nil || err != nil {
			for _, sessionID := range sessionIDs {
				if client, exists := entity.WSClients[sessionID]; exists && client.UserID == msg.UserID {
					// 메시지 전송
					err := client.Conn.WriteJSON(msg)
					if err != nil {
						fmt.Printf("Error sending message to user %d: %v\n", client.UserID, err)
						client.Conn.Close()

						// 클라이언트와 세션 삭제
						delete(entity.WSClients, sessionID)
						removeSessionFromRoom(client.RoomID, sessionID)
					}
				}
			}
		} else {
			// 정상적인 경우 모든 유저에게 메시지 브로드캐스트
			for _, sessionID := range sessionIDs {
				if client, exists := entity.WSClients[sessionID]; exists {
					// 메시지 전송
					err := client.Conn.WriteJSON(msg)
					if err != nil {
						fmt.Printf("Error broadcasting to user %d: %v\n", client.UserID, err)
						client.Conn.Close()

						// 클라이언트와 세션 삭제
						delete(entity.WSClients, sessionID)
						removeSessionFromRoom(client.RoomID, sessionID)
					}
				}
			}
		}

		// 방이 비어 있으면 삭제
		if len(entity.RoomSessions[msg.RoomID]) == 0 {
			delete(entity.RoomSessions, msg.RoomID)
		}
	}
}
