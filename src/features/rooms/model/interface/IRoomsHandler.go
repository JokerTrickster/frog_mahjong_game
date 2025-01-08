package _interface

import "github.com/labstack/echo/v4"

type ICreateRoomsHandler interface {
	Create(c echo.Context) error
}

type IV02CreateRoomsHandler interface {
	V02Create(c echo.Context) error
}
type IJoinPlayRoomsHandler interface {
	JoinPlay(c echo.Context) error
}

type IV02JoinRoomsHandler interface {
	V02Join(c echo.Context) error
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

type IMetaRoomsHandler interface {
	Meta(c echo.Context) error
}

type ICheckSessionRoomsHandler interface {
	CheckSession(c echo.Context) error
}