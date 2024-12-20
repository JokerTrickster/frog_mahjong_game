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

// 랜덤으로 방 매칭 (ws)
// @Router /v2.1/rooms/match/ws [get]
func match(c echo.Context) error {
	ws, err := entity.WSUpgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		fmt.Printf("WebSocket upgrade failed: %v\n", err)
		return err
	}

	req := &request.ReqWSMatch{}
	if err := utils.ValidateReq(c, req); err != nil {
		fmt.Printf("Invalid request: %v\n", err)
		return err
	}

	//토큰 검증
	err = utils.VerifyToken(req.Tkn)
	if err != nil {
		fmt.Printf("Token verification failed: %v\n", err)
		return err
	}

	userID, _, err := utils.ParseToken(req.Tkn)
	if err != nil {
		fmt.Printf("Failed to parse token: %v\n", err)
		return err
	}

	// 재접속 확인
	// 유저 상태가 abnormal 이면 해당 roomID를 가지고 온다.
	if req.SessionID != "" {
		roomID, _ := repository.MatchRedisSessionGet(context.Background(), req.SessionID)
		if roomID != 0 {
			// 기존 연결 복구
			restoreSession(ws, req.SessionID, roomID, userID)
			// 연결한 유저에게 메시지 정보를 전달해야 된다.
			return nil
		}
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
			utils.LogInfo(fmt.Sprintf("Room %d created by User %d.", roomID, userID))
		} else {
			roomID = rooms.ID
		}

		// 방 유저 정보 업데이트
		err = repository.MatchFindOneAndUpdateRoom(ctx, tx, roomID)
		if err != nil {
			return err
		}

		// room_user 생성
		roomUserDTO := CreateMatchRoomUserDTO(userID, int(roomID))
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
	sessionID := generateSessionID()
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
