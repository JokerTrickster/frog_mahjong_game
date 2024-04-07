package _interface

import "github.com/labstack/echo/v4"

type IStartGameHandler interface {
	Start(c echo.Context) error
}

type IDoraGameHandler interface {
	Dora(c echo.Context) error
}

type IOwnershipGameHandler interface {
	Ownership(c echo.Context) error
}
