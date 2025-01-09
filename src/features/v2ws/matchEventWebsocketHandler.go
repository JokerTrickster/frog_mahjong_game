package v2ws

import (
	"context"
	"fmt"
	"main/features/v2ws/model/entity"
	_errors "main/features/v2ws/model/errors"
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
	ctx := context.Background()
	var newErr *entity.ErrorInfo
	ws, err := entity.WSUpgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		SendWebSocketCloseMessage(ws, _errors.ErrCodeBadRequest, err.Error())
		return nil
	}
	req := &request.ReqWSMatch{}
	if err := utils.ValidateReq(c, req); err != nil {
		SendWebSocketCloseMessage(ws, _errors.ErrCodeBadRequest, err.Error())
		return err
	}

	//토큰 검증
	err = utils.VerifyToken(req.Tkn)
	if err != nil {
		SendWebSocketCloseMessage(ws, _errors.ErrCodeBadRequest, err.Error())
		return err
	}

	userID, _, err := utils.ParseToken(req.Tkn)
	if err != nil {
		SendWebSocketCloseMessage(ws, _errors.ErrCodeBadRequest, err.Error())
		return err
	}

	// 재접속 확인
	// 유저 상태가 abnormal 이면 해당 roomID를 가지고 온다.
	if req.SessionID != "" {
		roomID, _ := repository.MatchRedisSessionGet(context.Background(), req.SessionID)
		if roomID != 0 {
			// 기존 연결 복구
			if client, exists := entity.WSClients[req.SessionID]; exists {
				closeAndRemoveClient(client, req.SessionID, roomID)
			}

			restoreSession(ws, req.SessionID, roomID, userID)
			// 연결한 유저에게 메시지 정보를 전달해야 된다.
			//기존 유저 상태 변경
			err := repository.MatchPlayerStateUpdate(context.Background(), roomID, userID)
			if err != nil {
				return fmt.Errorf("%s", err.Msg)
			}
			return nil
		}
	}
	// 2. 비즈니스 로직
	var roomInfoMsg entity.RoomInfo
	// 기존 데이터 삭제
	newErr = cleanGameInfo(ctx, userID)
	if newErr != nil {
		return fmt.Errorf("%s", newErr.Msg)
	}

	// 대기중인 방 찾기
	rooms, newErr := repository.MatchFindOneWaitingRoom(ctx, uint(req.Count), uint(req.Timer))
	if newErr != nil && newErr.Msg != gorm.ErrRecordNotFound.Error() {
		roomInfoMsg.ErrorInfo = newErr
		return fmt.Errorf("%s", newErr.Msg)
	}
	// 트랜잭션으로 방 생성/업데이트 처리
	var roomID uint
	err = mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		if rooms == nil {
			// 방 생성
			roomDTO := CreateMatchRoomDTO(userID, req.Count, req.Timer)
			newRoomID, newErr := repository.MatchInsertOneRoom(ctx, roomDTO)
			if newErr != nil {
				SendWebSocketCloseMessage(ws, newErr.Code, newErr.Msg)
				return fmt.Errorf("%s", newErr.Msg)
			}
			roomID = uint(newRoomID)
		} else {
			roomID = rooms.ID
		}
		// 방 유저 정보 업데이트
		newErr = repository.MatchFindOneAndUpdateRoom(ctx, tx, roomID)
		if newErr != nil {
			SendWebSocketCloseMessage(ws, newErr.Code, newErr.Msg)
			return fmt.Errorf("%s", newErr.Msg)
		}

		// room_user 생성
		roomUserDTO := CreateMatchRoomUserDTO(userID, int(roomID))
		newErr = repository.MatchInsertOneRoomUser(ctx, tx, roomUserDTO)
		if newErr != nil {
			SendWebSocketCloseMessage(ws, newErr.Code, newErr.Msg)
			return fmt.Errorf("%s", newErr.Msg)
		}
		// 아이템 정보들을 가져온다.
		items, newErr := repository.MatchFindAllItems(ctx, tx)
		if newErr != nil {
			SendWebSocketCloseMessage(ws, newErr.Code, newErr.Msg)
			return fmt.Errorf("%s", newErr.Msg)
		}
		for _, item := range items {
			// user_items 아이템 정보 생성
			userItemDTO := CreateMatchUserItemDTO(userID, roomID, item)
			newErr = repository.MatchInsertOneUserItem(ctx, tx, userItemDTO)
			if newErr != nil {
				SendWebSocketCloseMessage(ws, newErr.Code, newErr.Msg)
				return fmt.Errorf("%s", newErr.Msg)
			}
		}

		return nil
	})
	if err != nil {
		fmt.Printf("Transaction error: %v\n", err)
		return err
	}

	// 세션 ID 생성
	sessionID := generateSessionID()

	// 세션 ID 저장
	newErr = repository.MatchRedisSessionSet(ctx, sessionID, roomID)
	if newErr != nil {
		SendWebSocketCloseMessage(ws, newErr.Code, newErr.Msg)
		return fmt.Errorf("%s", newErr.Msg)
	}

	// defer ws.Close()

	// 3. 새로운 세션 등록
	registerNewSession(ws, sessionID, roomID, userID)
	return nil
}
