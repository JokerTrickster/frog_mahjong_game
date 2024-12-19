package v2ws

import (
	"context"
	"fmt"
	"main/features/v2ws/model/entity"
	"main/features/v2ws/model/request"
	"main/features/v2ws/repository"
	"main/utils"
	"main/utils/db/mysql"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// 함께하기 방 생성 (패스워드 발급) (ws)
// @Router /v2.1/rooms/play/together/ws [get]
// @Summary 함께하기 방 생성 (패스워드 발급) (ws)
// @Description
// @Description ■ errCode with 400
// @Description PARAM_BAD : 파라미터 오류
// @Description USER_NOT_FOUND : 유저를 찾을 수 없음
// @Description USER_ALREADY_EXISTED : 이미 존재하는 유저
// @Description Room_NOT_FOUND : 방을 찾을 수 없음
// @Description Room_FULL : 방이 꽉 참
// @Description Room_USER_NOT_FOUND : 방 유저를 찾을 수 없음
// @Description
// @Description ■ errCode with 500
// @Description INTERNAL_SERVER : 내부 로직 처리 실패
// @Description INTERNAL_DB : DB 처리 실패
// @Description PLAYER_STATE_CHANGE_FAILED : 플레이어 상태 변경 실패
// @Param tkn query string true "access token"
// @Produce json
// @Success 200 {object} boolean
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags ws
func playTogether(c echo.Context) error {
	ws, err := entity.WSUpgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		fmt.Printf("WebSocket upgrade failed: %v\n", err)
		return nil
	}

	req := &request.ReqWSPlayTogether{}
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

	// 비즈니스 로직
	ctx := context.Background()

	// 기존 방과 유저 데이터 삭제
	err = repository.PlayTogetherDeleteRooms(ctx, userID)
	if err != nil {
		fmt.Printf("Failed to delete rooms: %v\n", err)
		return nil
	}

	err = repository.PlayTogetherDeleteRoomUsers(ctx, userID)
	if err != nil {
		fmt.Printf("Failed to delete room users: %v\n", err)
		return nil
	}

	// 대기중인 방이 있는지 체크
	var roomID uint

	err = mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		// 숫자로 이루어진 4자리 랜덤 패스워드 생성
		password := CreateRandomPassword()

		// 방 생성
		roomDTO := CreatePlayTogetherRoomDTO(userID, 2, 15, password)
		newRoomID, err := repository.PlayTogetherInsertOneRoom(ctx, roomDTO)
		if err != nil {
			return err
		}
		roomID = uint(newRoomID)

		// 방에 유저 추가
		err = repository.PlayTogetherAddPlayerToRoom(ctx, tx, roomID)
		if err != nil {
			return err
		}

		// room_user 생성
		roomUserDTO := CreatePlayTogetherRoomUserDTO(userID, int(roomID), "ready")
		err = repository.PlayTogetherInsertOneRoomUser(ctx, tx, roomUserDTO)
		if err != nil {
			return err
		}

		// 아이템 정보 가져오기
		items, err := repository.PlayTogetherFindAllItems(ctx, tx)
		if err != nil {
			return err
		}

		// 유저 아이템 추가
		for _, item := range items {
			userItemDTO := CreatePlayTogetherUserItemDTO(userID, roomID, item)
			err = repository.PlayTogetherInsertOneUserItem(ctx, tx, userItemDTO)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Transaction error: %v\n", err)
		return nil
	}

	defer ws.Close()

	// sessionID 생성
	sessionID := generateSessionID()

	// 초기 메시지 처리
	var initialMsg entity.WSMessage
	err = ws.ReadJSON(&initialMsg)
	if err != nil {
		fmt.Printf("Failed to read initial message: %v\n", err)
		return nil
	}

	initialMsg.UserID = userID
	initialMsg.RoomID = roomID
	initialMsg.SessionID = sessionID

	// 첫 번째 레벨 맵 초기화 (RoomSessions)
	if entity.RoomSessions == nil {
		entity.RoomSessions = make(map[uint][]string)
	}

	// 두 번째 레벨 맵 초기화 (WSClients)
	if entity.WSClients == nil {
		entity.WSClients = make(map[string]*entity.WSClient)
	}

	// 방 세션에 sessionID 추가
	entity.RoomSessions[roomID] = append(entity.RoomSessions[roomID], sessionID)

	// sessionID를 WSClients에 등록
	wsClient := &entity.WSClient{
		RoomID:    roomID,
		UserID:    userID,
		SessionID: sessionID,
		Conn:      ws,
		Closed:    false,
	}
	entity.WSClients[sessionID] = wsClient

	// 메시지 브로드캐스트
	entity.WSBroadcast <- initialMsg

	// Ping/Pong 관리 시작
	go HandlePingPong(wsClient)

	// 메시지 수신 루프
	for {
		var msg entity.WSMessage
		err := ws.ReadJSON(&msg)
		if err != nil {
			fmt.Printf("Error reading message for session %s: %v\n", sessionID, err)

			// 비정상 종료 처리
			wsClient.Closed = true
			AbnormalErrorHandling(roomID, sessionID)
			break
		}

		// 메시지 처리
		msg.RoomID = roomID
		msg.UserID = userID
		msg.SessionID = sessionID
		entity.WSBroadcast <- msg
	}

	// 연결 종료 및 클라이언트 정리
	delete(entity.WSClients, sessionID)
	removeSessionFromRoom(roomID, sessionID)

	if len(entity.RoomSessions[roomID]) == 0 {
		delete(entity.RoomSessions, roomID)
		fmt.Printf("Room %d deleted as it is empty.\n", roomID)
	}

	return nil
}
