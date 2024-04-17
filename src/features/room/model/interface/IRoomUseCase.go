package _interface

import (
	"context"
	"main/features/room/model/request"
	"main/features/room/model/response"
)

type ICreateRoomUseCase interface {
	Create(c context.Context, uID uint, email string, req *request.ReqCreate) (response.ResCreateRoom, error)
}

type IJoinRoomUseCase interface {
	Join(c context.Context, uID uint, email string, req *request.ReqJoin) (response.ResJoinRoom, error)
}

type IOutRoomUseCase interface {
	Out(c context.Context, uID uint, req *request.ReqOut) error
}

type IReadyRoomUseCase interface {
	Ready(c context.Context, uID uint, req *request.ReqReady) error
}

type IListRoomUseCase interface {
	List(c context.Context, page int, pageSize int) (response.ResListRoom, error)
}

type IUserListRoomUseCase interface {
	UserList(c context.Context, roomID uint) (response.ResUserListRoom, error)
}
