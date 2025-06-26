package _interface

import "github.com/labstack/echo/v4"

type IListProfilesHandler interface {
	List(c echo.Context) error
}

type IUploadProfilesHandler interface {
	Upload(c echo.Context) error
}

type IUpdateProfilesHandler interface {
	Update(c echo.Context) error
}
