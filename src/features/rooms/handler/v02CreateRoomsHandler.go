package handler

import (
	"fmt"
	"log"
	_interface "main/features/rooms/model/interface"
	"main/features/rooms/model/request"
	"main/utils"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type V02CreateRoomsHandler struct {
	UseCase _interface.IV02CreateRoomsUseCase
}

func NewV02CreateRoomsHandler(c *echo.Echo, useCase _interface.IV02CreateRoomsUseCase) _interface.IV02CreateRoomsHandler {
	handler := &V02CreateRoomsHandler{
		UseCase: useCase,
	}
	c.GET("/v0.2/rooms/create", handler.V02Create)
	return handler
}

// 방 생성 (ws)
// @Router /v0.2/rooms/create [get]
// @Summary 방 생성 (ws)
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
// @Param tkn header string false "accessToken"
// @Produce json
// @Success 200 {object} boolean
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags rooms
func (d *V02CreateRoomsHandler) V02Create(c echo.Context) error {
	ws, err := utils.WSUpgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	req := &request.ReqV02Create{}
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

	defer ws.Close()
	var initialMsg utils.WSMessage
	err = ws.ReadJSON(&initialMsg)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	roomID := initialMsg.RoomID
	// 첫 번째 레벨 맵 초기화
	if utils.WSClients == nil {
		utils.WSClients = make(map[uint]map[*websocket.Conn]utils.WSClient)
	}

	// 두 번째 레벨 맵 초기화
	if utils.WSClients[roomID] == nil {
		utils.WSClients[roomID] = make(map[*websocket.Conn]utils.WSClient)
	}
	utils.WSClients[roomID][ws] = utils.WSClient{
		RoomID: roomID,
		UserID: userID,
		Conn:   ws,
	}
	for {
		var msg utils.WSMessage
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(utils.WSClients[roomID], ws)
			break
		}
		msg.RoomID = roomID
		utils.WSBroadcast <- msg
	}

	return nil
}
