package ws

import (
	"context"
	"fmt"
	"main/features/ws/model/entity"
	_errors "main/features/ws/model/errors"
	"main/features/ws/model/request"
	"main/features/ws/repository"
	"main/utils"
	"main/utils/db/mysql"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// 함께하기 참여 (패스워드 필수) (ws)
// @Router /v0.1/rooms/join/play/ws [get]
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
	ctx := context.Background()
	var errInfo *entity.ErrorInfo
	ws, err := entity.WSUpgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		SendWebSocketCloseMessage(ws, _errors.ErrCodeBadRequest, err.Error())
		return nil
	}

	req := &request.ReqWSJoinPlay{}
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
	// 재접속 확인
	// 유저 상태가 abnormal 이면 해당 roomID를 가지고 온다.
	if req.SessionID != "" {
		roomID, _ := repository.MatchRedisSessionGet(context.Background(), req.SessionID)
		if roomID != 0 {
			_, errInfo = repository.JoinPlayFindOneRoom(ctx, roomID)
			if errInfo.Msg == gorm.ErrRecordNotFound.Error() {
				_ = repository.RedisSessionDelete(ctx, req.SessionID)
			} else {
				// 기존 연결 복구
				if client, exists := entity.WSClients[req.SessionID]; exists {
					closeAndRemoveClient(client, req.SessionID, roomID)
				}
				restoreSession(ws, req.SessionID, roomID, userID)
				// 연결한 유저에게 메시지 정보를 전달해야 된다.

				// 기존 유저 상태 변경
				err := repository.JoinPlayPlayerStateUpdate(context.Background(), roomID, userID)
				if err != nil {
					SendWebSocketCloseMessage(ws, err.Code, err.Msg)
					return nil
				}
				return nil
			}
		}
	}
	// 비즈니스 로직
	// 대기중인 방이 있는지 체크
	// var roomInfoMsg entity.RoomInfo
	var roomID uint
	// 기존 유저에 게임 정보를 모두 제거한다.
	errInfo = cleanGameInfo(ctx, userID)
	if errInfo != nil {
		SendWebSocketCloseMessage(ws, errInfo.Code, errInfo.Msg)
		return nil
	}
	rooms, errInfo := repository.JoinPlayFindOneWaitingRoom(ctx, req.Password)
	if errInfo != nil {
		SendWebSocketCloseMessage(ws, _errors.ErrCodeBadRequest, "비밀번호를 잘못 입력했습니다.")
		return nil
	}
	err = mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		roomID = rooms.ID
		// room 유저 수 증가
		errInfo = repository.JoinPlayFindOneAndUpdateRoom(ctx, tx, roomID)
		if errInfo != nil {
			SendWebSocketCloseMessage(ws, errInfo.Code, errInfo.Msg)
			return fmt.Errorf("%s", errInfo.Msg)
		}
		// room_user 생성
		roomUserDTO := CreateMatchRoomUserDTO(userID, int(roomID))
		errInfo = repository.JoinPlayInsertOneRoomUser(ctx, tx, roomUserDTO)
		if errInfo != nil {
			SendWebSocketCloseMessage(ws, errInfo.Code, errInfo.Msg)
			return fmt.Errorf("%s", errInfo.Msg)
		}
		return nil
	})
	if err != nil {
		return nil
	}

	// sessionID 생성
	sessionID := generateSessionID()
	// 세션 ID 저장
	errInfo = repository.RedisSessionSet(ctx, sessionID, roomID)
	if errInfo != nil {
		SendWebSocketCloseMessage(ws, errInfo.Code, errInfo.Msg)
		return nil
	}

	registerNewSession(ws, sessionID, roomID, userID)

	return nil
}
