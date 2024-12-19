package v2ws

import (
	"context"
	"fmt"
	"log"
	"main/features/v2ws/model/entity"
	"main/features/v2ws/model/request"
	"main/features/v2ws/repository"
	"main/utils"
	"main/utils/db/mysql"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// 랜덤으로 방 매칭 (ws)
// @Router /v2.1/rooms/match/ws [get]
func match(c echo.Context) error {
	ws, err := entity.WSUpgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		fmt.Printf("WebSocket upgrade failed: %v\n", err)
		return nil
	}

	req := &request.ReqWSMatch{}
	if err := utils.ValidateReq(c, req); err != nil {
		fmt.Printf("Invalid request: %v\n", err)
		return nil
	}

	err = utils.VerifyToken(req.Tkn)
	if err != nil {
		fmt.Printf("Token verification failed: %v\n", err)
		return nil
	}

	userID, _, err := utils.ParseToken(req.Tkn)
	if err != nil {
		fmt.Printf("Failed to parse token: %v\n", err)
		return nil
	}
	// 1. 재접속 확인
	// 유저 상태가 abnormal 이면 해당 roomID를 가지고 온다.
	if req.SessionID != "" {
		fmt.Println("세션이 있지? ", req.SessionID)
		roomID, _ := repository.MatchRedisSessionGet(context.Background(), req.SessionID)
		if roomID != 0 {
			fmt.Println("기존 연결 복구")
			// 기존 연결 복구
			restoreSession(ws, req.SessionID, roomID, userID)
			return nil
		}
		fmt.Println("세션 정보가 없습니다.========== ", roomID)
	}
	// 2. 비즈니스 로직
	ctx := context.Background()

	// 기존 데이터 삭제
	err = repository.MatchDeleteRooms(ctx, userID)
	if err != nil {
		fmt.Printf("Failed to delete rooms: %v\n", err)
		return nil
	}

	err = repository.MatchDeleteRoomUsers(ctx, userID)
	if err != nil {
		fmt.Printf("Failed to delete room users: %v\n", err)
		return nil
	}

	// 대기중인 방 찾기
	rooms, err := repository.MatchFindOneWaitingRoom(ctx, uint(req.Count), uint(req.Timer))
	if err != nil && err != gorm.ErrRecordNotFound {
		fmt.Printf("Failed to find waiting room: %v\n", err)
		return nil
	}
	// 트랜잭션으로 방 생성/업데이트 처리
	var roomID uint
	err = mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		if rooms.ID == 0 {
			// 방 생성
			roomDTO := CreateMatchRoomDTO(userID, req.Count, req.Timer)
			newRoomID, err := repository.MatchInsertOneRoom(ctx, roomDTO)
			if err != nil {
				return err
			}
			roomID = uint(newRoomID)
			fmt.Println("새로 방을 만들었습니다. 방 번호는", roomID)
		} else {
			roomID = rooms.ID
		}

		// 방 유저 정보 업데이트
		err = repository.MatchFindOneAndUpdateRoom(ctx, tx, roomID)
		if err != nil {
			return err
		}

		// room_user 생성
		roomUserDTO := CreateMatchRoomUserDTO(userID, int(roomID), "ready")
		err = repository.MatchInsertOneRoomUser(ctx, tx, roomUserDTO)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		fmt.Printf("Transaction error: %v\n", err)
		return nil
	}

	// 세션 ID 생성
	var sessionID string
	sessionID = generateSessionID()
	// 세션 ID 저장
	err = repository.MatchRedisSessionSet(ctx, sessionID, roomID)
	if err != nil {
		fmt.Printf("Failed to save session: %v\n", err)
		return nil
	}

	// defer ws.Close()

	// 3. 새로운 세션 등록
	registerNewSession(ws, sessionID, roomID, userID)
	return nil
}

// 기존 연결 복구
func restoreSession(ws *websocket.Conn, sessionID string, roomID uint, userID uint) {
	// 타이머 취소
	if timer, ok := reconnectTimers.Load(sessionID); ok {
		timer.(*time.Timer).Stop()
		reconnectTimers.Delete(sessionID)
		fmt.Printf("Reconnection successful for session %s in room %d. Timer canceled.\n", sessionID, roomID)
	}
	if client, ok := entity.WSClients[sessionID]; ok {
		fmt.Println("기존 연결 해제========")
		fmt.Println(client.RoomID, client.SessionID, client.UserID)
		// 기존 연결 닫기
		client.Closed = true
		client.Conn.Close()

		// 새로운 연결로 갱신
		client.Conn = ws
		client.Closed = false
		entity.WSClients[sessionID] = client

		fmt.Printf("User %d reconnected to Room %d with Session %s.\n", userID, roomID, sessionID)

		// 핑/퐁 핸들링 재시작
		go HandlePingPong(client)

		// 메시지 처리 루프 시작
		go readMessages(ws, sessionID, roomID, userID)
	} else {
		fmt.Println("세션 정보가 없습니다.========== ", sessionID)
		// 새로운 세션으로 등록
		entity.WSClients[sessionID] = &entity.WSClient{
			Conn:      ws,
			SessionID: sessionID,
			RoomID:    roomID,
			UserID:    userID,
			Closed:    false,
		}
		fmt.Println(len(entity.WSClients))
		entity.RoomSessions[roomID] = append(entity.RoomSessions[roomID], sessionID)
		fmt.Printf("New connection established for Session %s in Room %d by User %d.\n", sessionID, roomID, userID)

		// 핑/퐁 및 메시지 처리 시작
		go HandlePingPong(entity.WSClients[sessionID])
		go readMessages(ws, sessionID, roomID, userID)
	}
}

// 새로운 세션 등록
func registerNewSession(ws *websocket.Conn, sessionID string, roomID uint, userID uint) {
	// 세션 등록
	wsClient := &entity.WSClient{
		SessionID: sessionID,
		RoomID:    roomID,
		UserID:    userID,
		Conn:      ws,
		Closed:    false,
	}
	entity.WSClients[sessionID] = wsClient

	// 방에 세션 추가
	entity.RoomSessions[roomID] = append(entity.RoomSessions[roomID], sessionID)

	fmt.Printf("User %d joined Room %d with Session %s.\n", userID, roomID, sessionID)

	// 핑/퐁 핸들링 시작
	go HandlePingPong(wsClient)

	// 메시지 처리 루프 시작
	go readMessages(ws, sessionID, roomID, userID)

}

// 메시지 읽기 및 처리
func readMessages(ws *websocket.Conn, sessionID string, roomID uint, userID uint) {
	client := entity.WSClients[sessionID]
	defer func() {
		// 연결 종료 시 세션 정리
		client.Closed = true
		ws.Close()
		delete(entity.WSClients, sessionID)
		removeSessionFromRoom(roomID, sessionID)
		fmt.Println("Session", sessionID, "closed. Read loop stopped.")

		// 방 삭제 여부 확인
		if len(entity.RoomSessions[roomID]) == 0 {
			delete(entity.RoomSessions, roomID)
			log.Printf("Room %d deleted as it has no active sessions.", roomID)
		}
	}()

	for {
		if client.Closed {
			log.Printf("Session %s is closed. Stopping read loop.", sessionID)
			return
		}

		var msg entity.WSMessage
		err := ws.ReadJSON(&msg)
		if err != nil {
			if closeErr, ok := err.(*websocket.CloseError); ok {
				if closeErr.Code == websocket.CloseNormalClosure {
					log.Printf("Session %s closed normally (Code 1000).", sessionID)
					break
				}
				log.Printf("Session %s closed with error: %v", sessionID, closeErr)
			} else {
				log.Printf("Error reading message for session %s: %v", sessionID, err)
			}
			break
		}

		// 메시지 브로드캐스트
		msg.RoomID = roomID
		msg.UserID = userID
		msg.SessionID = sessionID
		entity.WSBroadcast <- msg
	}
}
