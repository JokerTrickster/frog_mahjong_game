package _interface

import "github.com/labstack/echo/v4"

type IMessageChatHandler interface {
	Message(c echo.Context) error
}

type IAuthChatHandler interface {
	Auth(c echo.Context) error
}
