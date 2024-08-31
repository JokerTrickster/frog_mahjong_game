package handler

import (
	"log"

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
	return nil
}
