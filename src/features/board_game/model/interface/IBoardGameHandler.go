package _interface

import "github.com/labstack/echo/v4"

type IFindItSoloPlayBoardGameHandler interface {
	FindItSoloPlay(c echo.Context) error
}
