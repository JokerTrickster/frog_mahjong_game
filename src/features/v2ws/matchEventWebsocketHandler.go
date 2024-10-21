package v2ws

import (
	"context"
	"fmt"
	"log"
	"main/features/v2ws/model/entity"
	"main/features/v2ws/model/request"
	"main/features/v2ws/repository"
	"main/utils"
	"main/utils/db/mysql"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// 랜덤으로 방 매칭 (ws)
// @Router /v2.1/rooms/match/ws [get]
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
// @Tags ws
func match(c echo.Context) error {
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

	rooms, err := repository.MatchFindOneWaitingRoom(ctx, uint(req.Count), uint(req.Timer))
	err = mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {

		if rooms.ID == 0 && err == nil {
			//대기 방이 없는 경우
			// 방 생성
			roomDTO := CreateMatchRoomDTO(userID, req.Count, req.Timer)
			newRoomID, err := repository.MatchInsertOneRoom(ctx, roomDTO)
			if err != nil {
				return err
			}
			roomID = uint(newRoomID)
		} else {
			roomID = rooms.ID
		}
		// room 유저 수 증가
		err = repository.MatchFindOneAndUpdateRoom(ctx, tx, roomID)
		if err != nil {
			return err
		}
		// 기존에 룸 유저 정보가 있으면 지운다.
		err = repository.MatchFindOneAndDeleteRoomUser(ctx, tx, userID)
		if err != nil {
			return err
		}
		// room_user 생성
		roomUserDTO := CreateMatchRoomUserDTO(userID, int(roomID), "ready")
		err = repository.MatchInsertOneRoomUser(ctx, tx, roomUserDTO)
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
