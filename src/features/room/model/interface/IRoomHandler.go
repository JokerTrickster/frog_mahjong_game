package _interface

import "github.com/labstack/echo/v4"

type ICreateRoomHandler interface {
	Create(c echo.Context) error
}

type IJoinRoomHandler interface {
	Join(c echo.Context) error
}

type IOutRoomHandler interface {
	Out(c echo.Context) error
}

type IReadyRoomHandler interface {
	Ready(c echo.Context) error
}

type IListRoomHandler interface {
	List(c echo.Context) error
}

type IUserListRoomHandler interface {
	UserList(c echo.Context) error
}
