package sequence

import (
	"context"
	"fmt"
	"main/features/sequence/model/entity"
	_errors "main/features/sequence/model/errors"
	"main/features/sequence/model/request"
	"main/features/sequence/repository"
	"main/utils"
	"main/utils/db/mysql"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// 랜덤으로 방 매칭 (ws)
// @Router /sequence/v0.1/rooms/match/ws [get]
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

	// 2. 비즈니스 로직
	var messageInfo entity.MessageInfo
	// 기존 데이터 삭제
	newErr = cleanGameInfo(ctx, userID)
	if newErr != nil {
		return fmt.Errorf("%s", newErr.Msg)
	}

	// 대기중인 방 찾기
	rooms, newErr := repository.MatchFindOneWaitingRoom(ctx)
	if newErr != nil && newErr.Msg != gorm.ErrRecordNotFound.Error() {
		messageInfo.ErrorInfo = newErr
		return fmt.Errorf("%s", newErr.Msg)
	}

	// 트랜잭션으로 방 생성/업데이트 처리
	var roomID uint
	err = mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		if rooms == nil {
			// 방 생성
			roomDTO := CreateMatchRoomDTO(userID)
			newRoomID, newErr := repository.MatchInsertOneRoom(ctx, roomDTO)
			if newErr != nil {
				SendWebSocketCloseMessage(ws, newErr.Code, newErr.Msg)
				return fmt.Errorf("%s", newErr.Msg)
			}
			roomID = uint(newRoomID)
			//게임 정보 생성
			roomSettingDTO := CreateRoomSetting(roomID)
			newErr = repository.MatchInsertOneRoomSetting(ctx, tx, roomSettingDTO)
			if newErr != nil {
				SendWebSocketCloseMessage(ws, newErr.Code, newErr.Msg)
				return fmt.Errorf("%s", newErr.Msg)
			}
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
		roomUserDTO := CreateMatchRoomUserDTO(roomID, userID)
		newErr = repository.MatchInsertOneRoomUser(ctx, tx, roomUserDTO)
		if newErr != nil {
			SendWebSocketCloseMessage(ws, newErr.Code, newErr.Msg)
			return fmt.Errorf("%s", newErr.Msg)
		}
		//TODO 유저 코인 1개를 제거한다.

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
