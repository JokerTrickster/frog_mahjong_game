package ws

import (
	"context"
	"fmt"
	"main/features/ws/model/entity"
	"main/features/ws/model/request"
	"main/features/ws/repository"
	"main/utils"
	"main/utils/db/mysql"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// 함께하기 방 생성 (패스워드 발급) (ws)
// @Router /v0.1/rooms/play/together/ws [get]
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
		fmt.Println(err)
		return nil
	}

	req := &request.ReqWSPlayTogether{}
	if err := utils.ValidateReq(c, req); err != nil {
		fmt.Println(err)
		return nil
	}

	err = utils.VerifyToken(req.Tkn)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	userID, _, err := utils.ParseToken(req.Tkn)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	// 2. 비즈니스 로직

	ctx := context.Background()
	var roomID uint
	// var roomInfoMsg entity.RoomInfo
	// 대기중인 방이 있는지 체크
	err = mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		//숫자로 이루어진 4개 랜덤값을 생성한다.
		password := CreateRandomPassword()
		// 방 생성
		roomDTO := CreatePlayTogetherRoomDTO(userID, 2, 15, password)
		newRoomID, err := repository.PlayTogetherInsertOneRoom(ctx, roomDTO)
		if err != nil {
			return err
		}
		roomID = uint(newRoomID)
		// room 유저 수 증가
		err = repository.PlayTogetherAddPlayerToRoom(ctx, tx, roomID)
		if err != nil {
			return err
		}
		// 기존에 룸 유저 정보가 있으면 지운다.
		err = repository.PlayTogetherFindOneAndDeleteRoomUser(ctx, tx, userID)
		if err != nil {
			return err
		}
		// room_user 생성
		roomUserDTO := CreatePlayTogetherRoomUserDTO(userID, int(roomID), "ready")
		err = repository.PlayTogetherInsertOneRoomUser(ctx, tx, roomUserDTO)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
		return nil
	}

	// sessionID 생성
	sessionID := generateSessionID()
	// 세션 ID 저장
	newErr := repository.RedisSessionSet(ctx, sessionID, roomID)
	if newErr != nil {
		fmt.Printf("Failed to save session: %v\n", newErr)
		return nil
	}

	registerNewSession(ws, sessionID, roomID, userID)

	return nil
}
