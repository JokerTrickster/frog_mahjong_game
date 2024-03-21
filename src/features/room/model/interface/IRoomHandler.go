package _interface

import "github.com/labstack/echo/v4"

type ICreateRoomHandler interface {
	Create(c echo.Context) error
}
