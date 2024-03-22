package _interface

import (
	"context"
	"main/features/room/model/request"
)

type ICreateRoomUseCase interface {
	Create(c context.Context, uID uint, email string, req *request.ReqCreate) error
}

type IJoinRoomUseCase interface {
	Join(c context.Context, uID uint, email string, req *request.ReqJoin) error
}
