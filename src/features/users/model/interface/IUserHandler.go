package _interface

import "github.com/labstack/echo/v4"

type IGetUsersHandler interface {
	Get(c echo.Context) error
}
