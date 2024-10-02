package _interface

import "github.com/labstack/echo/v4"

type IGetUsersHandler interface {
	Get(c echo.Context) error
}

type IListUsersHandler interface {
	List(c echo.Context) error
}
type IUpdateUsersHandler interface {
	Update(c echo.Context) error
}

type IDeleteUsersHandler interface {
	Delete(c echo.Context) error
}

type IListProfilesUsersHandler interface {
	ListProfiles(c echo.Context) error
}
