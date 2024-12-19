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

func PlayTogetherEventWebsocket(msg *entity.WSMessage) {
	ctx := context.Background()
	uID := msg.UserID

	//string to struct
	req := request.ReqWSPlayTogetherEvent{}
	err := json.Unmarshal([]byte(msg.Message), &req)
	if err != nil {
		log.Fatalf("JSON 언마샬링 에러: %s", err)
	}

	//비즈니스 로직
	roomInfoMsg := entity.RoomInfo{}
	preloadUsers := []entity.RoomUsers{}
	roomID, err := repository.PlayTogetherFindOneRoomUsers(ctx, uID)
	if err != nil {
		log.Fatalf("방 유저 정보 조회 에러: %s", err)
	}
	err = mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		// 방 정보를 업데이트 한다. (타이머, 인원 수)
		err := repository.PlayTogetherFindOneAndUpdateRoom(ctx, tx, roomID, uint(req.Count), uint(req.Timer))
		if err != nil {
			return err
		}

		//유저 정보를 업데이트 한다.
		err = repository.PlayTogetherFindOneAndUpdateUser(ctx, tx, uID, roomID)
		if err != nil {
			return err
		}
		// 미션 3개를 생성한다.
		err = repository.PlayTogetherCreateMissions(ctx, tx, roomID)
		if err != nil {
			return err
		}
		preloadUsers, err = repository.PlayTogetherFindAllRoomUsers(ctx, tx, roomID)
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
		if roomInfoMsg.ErrorInfo.Msg == "방이 꽉 찼습니다." {
			roomInfoMsg.ErrorInfo.Type = _errors.ErrRoomFull
		} else if roomInfoMsg.ErrorInfo.Msg == "비밀번호가 일치하지 않습니다." {
			roomInfoMsg.ErrorInfo.Type = _errors.ErrWrongPassword
		} else if roomInfoMsg.ErrorInfo.Msg == "게임 중인 방입니다." {
			roomInfoMsg.ErrorInfo.Type = _errors.ErrGameInProgress
		}
	}

	// 메시지 생성
	roomInfoMsg = *CreateRoomInfoMSG(ctx, preloadUsers, 1, roomInfoMsg.ErrorInfo, 0)
	roomInfoMsg.GameInfo.AllReady = false

	if len(preloadUsers) == req.Count {
		roomInfoMsg.GameInfo.IsFull = true
		roomInfoMsg.GameInfo.AllReady = true
	}
	// 구조체를 JSON 문자열로 변환 (마샬링)
	message, err := CreateMessage(&roomInfoMsg)
	if err != nil {
		fmt.Println(err)
	}
	msg.Message = message
	msg.RoomID = roomID

	// 방 유저들에게 메시지 전달
	if sessionIDs, ok := entity.RoomSessions[msg.RoomID]; ok {
		// 에러 발생 시 이벤트 요청한 유저에게만 메시지를 전달
		if roomInfoMsg.ErrorInfo != nil || err != nil {
			for _, sessionID := range sessionIDs {
				if client, exists := entity.WSClients[sessionID]; exists && client.UserID == msg.UserID {
					// 메시지 전송
					if err := client.Conn.WriteJSON(msg); err != nil {
						fmt.Printf("Error sending message to session %s: %v\n", sessionID, err)

						// 클라이언트 종료 및 세션 제거
						client.Close()
						delete(entity.WSClients, sessionID)
						removeSessionFromRoom(msg.RoomID, sessionID)
					}
				}
			}
		} else {
			// 정상적인 경우 방의 모든 유저에게 메시지 전송
			for _, sessionID := range sessionIDs {
				if client, exists := entity.WSClients[sessionID]; exists {
					// 메시지 전송
					if err := client.Conn.WriteJSON(msg); err != nil {
						fmt.Printf("Error sending message to session %s: %v\n", sessionID, err)

						// 클라이언트 종료 및 세션 제거
						client.Close()
						delete(entity.WSClients, sessionID)
						removeSessionFromRoom(msg.RoomID, sessionID)
					}
				}
			}
		}

		// 방이 비어 있으면 삭제
		if len(entity.RoomSessions[msg.RoomID]) == 0 {
			delete(entity.RoomSessions, msg.RoomID)
			fmt.Printf("Room %d deleted as it is empty.\n", msg.RoomID)
		}
	}

}
