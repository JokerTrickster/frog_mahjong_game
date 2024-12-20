package v2ws

import (
	"context"
	"fmt"
	"main/features/v2ws/model/entity"
	"main/features/v2ws/model/request"
	"main/features/v2ws/repository"
	"main/utils"
	"main/utils/db/mysql"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// 함께하기 참여 (패스워드 필수) (ws)
// @Router /v2.1/rooms/join/play/ws [get]
// @Summary 함께하기 참여 (패스워드 필수) (ws)
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
// @Param password query string true "password"
// @Produce json
// @Success 200 {object} boolean
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags ws
func joinPlay(c echo.Context) error {
	ws, err := entity.WSUpgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		fmt.Printf("WebSocket upgrade failed: %v\n", err)
		return nil
	}

	req := &request.ReqWSJoinPlay{}
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
	// 재접속 확인
	// 유저 상태가 abnormal 이면 해당 roomID를 가지고 온다.
	if req.SessionID != "" {
		roomID, _ := repository.JoinRedisSessionGet(context.Background(), req.SessionID)
		if roomID != 0 {
			// 기존 연결 복구
			restoreSession(ws, req.SessionID, roomID, userID)
			// 연결한 유저에게 메시지 정보를 전달해야 된다.
			return nil
		}
	}

	// 비즈니스 로직
	// 대기중인 방이 있는지 체크
	ctx := context.Background()

	// 기존 방과 유저 데이터 삭제
	err = repository.JoinPlayDeleteRooms(ctx, userID)
	if err != nil {
		fmt.Printf("Failed to delete rooms: %v\n", err)
		return nil
	}

	err = repository.JoinPlayDeleteRoomUsers(ctx, userID)
	if err != nil {
		fmt.Printf("Failed to delete room users: %v\n", err)
		return nil
	}

	var roomID uint

	// 방 찾기
	rooms, err := repository.JoinPlayFindOneWaitingRoom(ctx, req.Password)
	if err != nil {
		message := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "비밀번호를 잘못 입력했습니다.")
		ws.WriteMessage(websocket.CloseMessage, message)
		return nil
	}

	// 트랜잭션으로 방 정보 및 유저 데이터 갱신
	err = mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		roomID = rooms.ID

		// 방 유저 수 증가
		err = repository.JoinPlayFindOneAndUpdateRoom(ctx, tx, roomID)
		if err != nil {
			return err
		}

		// 방 유저 정보 추가
		roomUserDTO := CreateMatchRoomUserDTO(userID, int(roomID))
		err = repository.JoinPlayInsertOneRoomUser(ctx, tx, roomUserDTO)
		if err != nil {
			return err
		}

		// 아이템 정보 가져오기
		items, err := repository.JoinFindAllItems(ctx, tx)
		if err != nil {
			return err
		}

		// 유저 아이템 추가
		for _, item := range items {
			userItemDTO := CreateJoinUserItemDTO(userID, roomID, item)
			err = repository.JoinInsertOneUserItem(ctx, tx, userItemDTO)
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

	// sessionID 생성
	sessionID := generateSessionID()
	// 세션 ID 저장
	err = repository.PlayTogetherRedisSessionSet(ctx, sessionID, roomID)
	if err != nil {
		fmt.Printf("Failed to save session: %v\n", err)
		return nil
	}

	registerNewSession(ws, sessionID, roomID, userID)
	return nil
}
