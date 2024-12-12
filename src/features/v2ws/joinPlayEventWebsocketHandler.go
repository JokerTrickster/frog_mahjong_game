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
		fmt.Println(err)
		return nil
	}

	req := &request.ReqWSJoinPlay{}
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

	// 비즈니스 로직
	// 대기중인 방이 있는지 체크
	ctx := context.Background()

	// rooms에 owner_id가 userID인 데이터 모두 삭제
	err = repository.JoinPlayDeleteRooms(ctx, userID)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	// room_users 에 user_id가 userID인 데이터 모두 삭제
	err = repository.JoinPlayDeleteRoomUsers(ctx, userID)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	var roomID uint

	rooms, err := repository.JoinPlayFindOneWaitingRoom(ctx, req.Password)
	if err != nil {
		message := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "비밀번호를 잘못 입력했습니다.")
		ws.WriteMessage(websocket.CloseMessage, message)
		return nil
	}
	err = mysql.Transaction(mysql.GormMysqlDB, func(tx *gorm.DB) error {
		roomID = rooms.ID
		// room 유저 수 증가
		err = repository.JoinPlayFindOneAndUpdateRoom(ctx, tx, roomID)
		if err != nil {
			return err
		}

		// room_user 생성
		roomUserDTO := CreateMatchRoomUserDTO(userID, int(roomID), "ready")
		err = repository.JoinPlayInsertOneRoomUser(ctx, tx, roomUserDTO)
		if err != nil {
			return err
		}


		// 아이템 정보들을 가져온다.
		items, err := repository.JoinFindAllItems(ctx, tx)
		if err != nil {
			return err
		}
		for _, item := range items {
			// user_items 아이템 정보 생성
			userItemDTO := CreateJoinUserItemDTO(userID, roomID, item)
			err = repository.JoinInsertOneUserItem(ctx, tx, userItemDTO)
			if err != nil {
				return err
			}
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
