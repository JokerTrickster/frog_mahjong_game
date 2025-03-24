package find_it

import (
	"context"
	"fmt"
	"main/features/find_it/model/entity"
	_errors "main/features/find_it/model/errors"
	"main/features/find_it/model/request"
	"main/features/find_it/repository"
	"main/utils"
	"main/utils/db/mysql"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// 함께하기 방 생성 (패스워드 발급) (ws)
// @Router /find-it/v0.1/rooms/play/together/ws [get]
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
	ctx := context.Background()
	var newErr *entity.ErrorInfo
	ws, err := entity.WSUpgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		SendWebSocketCloseMessage(ws, _errors.ErrCodeBadRequest, err.Error())
		return nil
	}

	req := &request.ReqWSPlayTogether{}
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
	var roomID uint
	// 기존 방과 유저 데이터 삭제
	newErr = cleanGameInfo(ctx, userID)
	if newErr != nil {
		SendWebSocketCloseMessage(ws, newErr.Code, newErr.Msg)
		return fmt.Errorf("%s", newErr.Msg)
	}
	// 대기중인 방이 있는지 체크
	err = mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		// 숫자로 이루어진 4자리 랜덤 패스워드 생성
		password := CreateRandomPassword()

		// 방 생성
		roomDTO := CreatePlayTogetherRoomDTO(userID, 2, 15, password)
		newRoomID, err := repository.PlayTogetherInsertOneRoom(ctx, roomDTO)
		if err != nil {
			SendWebSocketCloseMessage(ws, newErr.Code, newErr.Msg)
			return fmt.Errorf("%s", err.Msg)
		}
		roomID = uint(newRoomID)

		// 방에 유저 추가
		err = repository.PlayTogetherAddPlayerToRoom(ctx, tx, roomID)
		if err != nil {
			SendWebSocketCloseMessage(ws, newErr.Code, newErr.Msg)
			return fmt.Errorf("%s", err.Msg)
		}

		// room_user 생성
		roomUserDTO := CreatePlayTogetherGameRoomUserDTO(userID, int(roomID))
		err = repository.PlayTogetherInsertOneRoomUser(ctx, tx, roomUserDTO)
		if err != nil {
			SendWebSocketCloseMessage(ws, newErr.Code, newErr.Msg)
			return fmt.Errorf("%s", err.Msg)
		}
		return nil
	})
	if err != nil {
		return nil
	}

	// sessionID 생성
	sessionID := generateSessionID()
	// 세션 ID 저장
	newErr = repository.PlayTogetherRedisSessionSet(ctx, sessionID, roomID)
	if newErr != nil {
		SendWebSocketCloseMessage(ws, newErr.Code, newErr.Msg)
		return nil
	}

	registerNewSession(ws, sessionID, roomID, userID)
	return nil
}
