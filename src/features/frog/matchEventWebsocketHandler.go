package ws

import (
	"context"
	"fmt"
	"main/features/frog/model/entity"
	_errors "main/features/frog/model/errors"
	"main/features/frog/model/request"
	"main/features/frog/repository"
	"main/utils"
	"main/utils/db/mysql"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// 랜덤으로 방 매칭 (ws)
// @Router /frog/v0.1/rooms/match/ws [get]
// @Summary 랜덤으로 방 매칭 (ws)
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
// @Param timer query int true "timer"
// @Param count query int true "count"
// @Produce json
// @Success 200 {object} boolean
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags frog
func match(c echo.Context) error {
	ctx := context.Background()
	var errInfo *entity.ErrorInfo
	ws, err := entity.WSUpgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		SendWebSocketCloseMessage(ws, _errors.ErrCodeBadRequest, err.Error())
		return nil
	}

	req := &request.ReqWSMatch{}
	if err := utils.ValidateReq(c, req); err != nil {
		SendWebSocketCloseMessage(ws, _errors.ErrCodeBadRequest, err.Error())
		return nil
	}

	err = utils.VerifyToken(req.Tkn)
	if err != nil {
		SendWebSocketCloseMessage(ws, _errors.ErrCodeBadRequest, err.Error())
		return nil
	}

	userID, _, err := utils.ParseToken(req.Tkn)
	if err != nil {
		SendWebSocketCloseMessage(ws, _errors.ErrCodeBadRequest, err.Error())
		return nil
	}
	// 비즈니스 로직
	// var roomInfoMsg entity.RoomInfo
	var roomID uint
	// 기존 유저에 게임 정보를 모두 제거한다.
	errInfo = cleanGameInfo(ctx, userID)
	if errInfo != nil {
		SendWebSocketCloseMessage(ws, errInfo.Code, errInfo.Msg)
		return nil
	}
	// 대기중인 방이 있는지 체크
	rooms, errInfo := repository.MatchFindOneWaitingRoom(ctx)
	if errInfo != nil {
		SendWebSocketCloseMessage(ws, errInfo.Code, errInfo.Msg)
		return nil
	}
	err = mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {

		if rooms.ID == 0 && err == nil {
			//대기 방이 없는 경우
			// 방 생성
			roomDTO := CreateMatchRoomDTO(userID)
			newRoomID, errInfo := repository.MatchInsertOneRoom(ctx, roomDTO)
			if errInfo != nil {
				SendWebSocketCloseMessage(ws, errInfo.Code, errInfo.Msg)
				return fmt.Errorf("%s", errInfo.Msg)
			}

			roomID = uint(newRoomID)
		} else {
			roomID = rooms.ID
		}
		// room 유저 수 증가
		errInfo = repository.MatchFindOneAndUpdateRoom(ctx, tx, roomID)
		if errInfo != nil {
			SendWebSocketCloseMessage(ws, errInfo.Code, errInfo.Msg)
			return fmt.Errorf("%s", errInfo.Msg)
		}
		// room_user 생성
		roomUserDTO := CreateMatchRoomUserDTO(userID, int(roomID))
		errInfo = repository.MatchInsertOneRoomUser(ctx, tx, roomUserDTO)
		if errInfo != nil {
			SendWebSocketCloseMessage(ws, errInfo.Code, errInfo.Msg)
			return fmt.Errorf("%s", errInfo.Msg)
		}

		return nil
	})
	if err != nil {
		return nil
	}

	// 세션 ID 생성
	sessionID := generateSessionID()

	// 세션 ID 저장
	errInfo = repository.RedisSessionSet(ctx, sessionID, roomID)
	if errInfo != nil {
		SendWebSocketCloseMessage(ws, errInfo.Code, errInfo.Msg)
		return fmt.Errorf("%s", errInfo.Msg)
	}

	// defer ws.Close()

	// 3. 새로운 세션 등록
	registerNewSession(ws, sessionID, roomID, userID)
	return nil
}
