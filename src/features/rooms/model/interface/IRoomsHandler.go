package _interface

import "github.com/labstack/echo/v4"

type ICreateRoomsHandler interface {
	Create(c echo.Context) error
}

type IJoinRoomsHandler interface {
	Join(c echo.Context) error
}

type IOutRoomsHandler interface {
	Out(c echo.Context) error
}

type IReadyRoomsHandler interface {
	Ready(c echo.Context) error
}

type IListRoomsHandler interface {
	List(c echo.Context) error
}

type IUserListRoomsHandler interface {
	UserList(c echo.Context) error
}
