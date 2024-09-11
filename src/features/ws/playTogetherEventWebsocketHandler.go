package ws

import (
	"context"
	"fmt"
	"log"
	"main/features/ws/model/entity"
	"main/features/ws/model/request"
	"main/features/ws/repository"
	"main/utils"
	"main/utils/db/mysql"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// 방 매칭 참여 (ws)
// @Router /v0.1/rooms/play/together/ws [get]
// @Summary 방 참여 (ws)
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
// @Tags ws
func playTogether(c echo.Context) error {
	ws, err := entity.WSUpgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	req := &request.ReqWSMatch{}
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

	// 대기중인 방이 있는지 체크
	ctx := context.Background()
	var roomID uint

	err = mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		//숫자로 이루어진 6개 랜덤값을 생성한다.
		password := CreateRandomPassword()
		// 방 생성
		roomDTO := CreatePlayTogetherRoomDTO(userID, req.Count, req.Timer, password)
		newRoomID, err := repository.PlayTogetherInsertOneRoom(ctx, roomDTO)
		if err != nil {
			return err
		}
		roomID = uint(newRoomID)
		// room 유저 수 증가
		err = repository.PlayTogetherFindOneAndUpdateRoom(ctx, tx, roomID)
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

	defer ws.Close()
	var initialMsg entity.WSMessage
	err = ws.ReadJSON(&initialMsg)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	initialMsg.UserID = userID
	// 첫 번째 레벨 맵 초기화
	if entity.WSClients == nil {
		entity.WSClients = make(map[uint]map[*websocket.Conn]*entity.WSClient)
	}

	// 두 번째 레벨 맵 초기화
	if entity.WSClients[roomID] == nil {
		entity.WSClients[roomID] = make(map[*websocket.Conn]*entity.WSClient)
	}
	wsClient := &entity.WSClient{
		RoomID: roomID,
		UserID: userID,
		Conn:   ws,
		Closed: false,
	}
	entity.WSClients[roomID][ws] = wsClient
	entity.WSBroadcast <- initialMsg
	go HandlePingPong(wsClient)

	for {
		var msg entity.WSMessage
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(entity.WSClients[roomID], ws)
			break
		}
		msg.RoomID = roomID
		msg.UserID = userID
		entity.WSBroadcast <- msg
	}

	return nil
}
