package _interface

import (
	"context"
	"main/features/rooms/model/request"
	"main/features/rooms/model/response"
)

type ICreateRoomsUseCase interface {
	Create(c context.Context, uID uint, email string, req *request.ReqCreate) (response.ResCreateRoom, error)
}
type IV02CreateRoomsUseCase interface {
	V02Create(c context.Context, uID uint, email string, req *request.ReqV02Create) (response.ResV02CreateRoom, error)
}

type IJoinRoomsUseCase interface {
	Join(c context.Context, uID uint, email string, req *request.ReqJoin) (response.ResJoinRoom, error)
}

type IOutRoomsUseCase interface {
	Out(c context.Context, uID uint, req *request.ReqOut) error
}

type IReadyRoomsUseCase interface {
	Ready(c context.Context, uID uint, req *request.ReqReady) error
}

type IListRoomsUseCase interface {
	List(c context.Context, page int, pageSize int) (response.ResListRoom, error)
}

type IUserListRoomsUseCase interface {
	UserList(c context.Context, RoomID uint) (response.ResUserListRoom, error)
}
