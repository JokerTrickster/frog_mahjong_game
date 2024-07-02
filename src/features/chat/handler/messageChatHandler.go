package handler

import (
	"log"
	"net/http"

	_interface "main/features/chat/model/interface"

	"main/utils"

	"github.com/labstack/echo/v4"
)

type MessageChatHandler struct {
	UseCase _interface.IMessageChatUseCase
}

func NewMessageChatHandler(c *echo.Echo, useCase _interface.IMessageChatUseCase) _interface.IMessageChatHandler {
	handler := &MessageChatHandler{
		UseCase: useCase,
	}
	c.GET("/v0.1/chat/message", handler.Message)
	return handler
}

// 챗 메시지 보내기
// @Router /v0.1/chat/message [get]
// @Summary 챗 메시지 보내기
// @Description
// @Description ■ errCode with 400
// @Description PARAM_BAD : 파라미터 오류
// @Description NOT_OWNER : 방장이 시작 요청을 하지 않음
// @Description NOT_FIRST_PLAYER : 첫 플레이어가 아님
// @Description
// @Description ■ errCode with 500
// @Description INTERNAL_SERVER : 내부 로직 처리 실패
// @Description INTERNAL_DB : DB 처리 실패
// @Produce json
// @Success 200 {object} boolean
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags chat
func (d *MessageChatHandler) Message(c echo.Context) error {

	ws, err := utils.Upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()
	var initialMsg utils.Message
	err = ws.ReadJSON(&initialMsg)
	if err != nil {
		log.Printf("error: %v", err)
		return err
	}
	utils.Clients[ws] = utils.Client{
		Conn: ws,
		Name: initialMsg.Username,
	}
	for {
		var msg utils.Message
		err := ws.ReadJSON(&msg)

		if err != nil {
			log.Printf("error: %v", err)
			delete(utils.Clients, ws)
			break
		}
		msg.Username = initialMsg.Username
		utils.Broadcast <- msg
	}
	return c.JSON(http.StatusOK, true)
}
