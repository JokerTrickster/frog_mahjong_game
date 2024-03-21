package _interface

import (
	"context"
	"main/utils/db/mysql"
)

type ICreateRoomRepository interface {
	InsertOneRoom(ctx context.Context, roomDTO mysql.Rooms) (int, error)
	InsertOneRoomUser(ctx context.Context, roomUserDTO mysql.RoomUsers) error
}
