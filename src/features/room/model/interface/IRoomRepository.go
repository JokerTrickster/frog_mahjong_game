package _interface

import (
	"context"
	"main/features/room/model/request"
	"main/utils/db/mysql"
)

type ICreateRoomRepository interface {
	InsertOneRoom(ctx context.Context, roomDTO mysql.Rooms) (int, error)
	InsertOneRoomUser(ctx context.Context, roomUserDTO mysql.RoomUsers) error
	FindOneAndUpdateUser(ctx context.Context, uID uint, roomID uint) error
}

type IJoinRoomRepository interface {
	FindOneRoom(ctx context.Context, req *request.ReqJoin) (mysql.Rooms, error)
	FindOneAndUpdateRoom(ctx context.Context, roomID uint) error
	FindOneAndUpdateUser(ctx context.Context, uID uint, roomID uint) error
	InsertOneRoomUser(ctx context.Context, roomUserDTO mysql.RoomUsers) error
}

type IOutRoomRepository interface {
	FindOneAndDeleteRoomUser(ctx context.Context, uID uint, roomID uint) error
	FindOneAndUpdateRoom(ctx context.Context, roomID uint) error
	FindOneAndUpdateUser(ctx context.Context, uID uint) error
}

type IReadyRoomRepository interface {
	FindOneAndUpdateRoomUser(ctx context.Context, uID uint, req *request.ReqReady) error
}
type IListRoomRepository interface {
	FindRoomList(ctx context.Context, page int, pageSize int) ([]mysql.Rooms, error)
	CountRoomList(ctx context.Context) (int, error)
}
type IUserListRoomRepository interface {
}
