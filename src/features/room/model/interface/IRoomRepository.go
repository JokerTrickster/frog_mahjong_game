package _interface

import (
	"context"
	"main/features/room/model/request"
	"main/utils/db/mysql"
)

type ICreateRoomRepository interface {
	InsertOneRoom(ctx context.Context, roomDTO mysql.Rooms) (int, error)
	InsertOneRoomUser(ctx context.Context, roomUserDTO mysql.RoomUsers) error
}

type IJoinRoomRepository interface {
	FindOneRoom(ctx context.Context, req *request.ReqJoin) (mysql.Rooms, error)
	FindOneAndUpdateRoom(ctx context.Context, roomID uint) error
	FindOneAndUpdateUser(ctx context.Context, uID uint, roomID uint) error
	InsertOneRoomUser(ctx context.Context, roomUserDTO mysql.RoomUsers) error
}