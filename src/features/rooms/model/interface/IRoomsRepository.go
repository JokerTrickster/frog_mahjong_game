package _interface

import (
	"context"
	"main/features/rooms/model/request"
	"main/features/rooms/model/response"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

type ICreateRoomsRepository interface {
	InsertOneRoom(ctx context.Context, RoomsDTO mysql.Rooms) (int, error)
	InsertOneRoomUser(ctx context.Context, RoomsUserDTO mysql.RoomUsers) error
	FindOneAndUpdateUser(ctx context.Context, uID uint, RoomID uint) error
}

type IJoinRoomsRepository interface {
	FindOneRoom(ctx context.Context, req *request.ReqJoin) (mysql.Rooms, error)
	FindOneAndUpdateRoom(ctx context.Context, RoomID uint) error
	FindOneAndUpdateUser(ctx context.Context, uID uint, RoomsID uint) error
	InsertOneRoomUser(ctx context.Context, RoomUserDTO mysql.RoomUsers) error
}

type IOutRoomsRepository interface {
	FindOneAndDeleteRoomUser(ctx context.Context, tx *gorm.DB, uID uint, RoomID uint) error
	FindOneAndUpdateRoom(ctx context.Context, tx *gorm.DB,RoomID uint) (mysql.Rooms, error)
	FindOneAndUpdateUser(ctx context.Context, tx *gorm.DB,uID uint) error
	FindOneAndDeleteRoom(ctx context.Context,tx *gorm.DB, RoomID uint) error
	FindOneRoomUser(ctx context.Context, tx *gorm.DB,RoomID uint) (mysql.RoomUsers, error)
	ChangeRoomOnwer(ctx context.Context,tx *gorm.DB, RoomID uint, ownerID uint) error
	FindOneUser(ctx context.Context,tx *gorm.DB,uID uint) (mysql.Users, error)
}

type IReadyRoomsRepository interface {
	FindOneAndUpdateRoomUser(ctx context.Context, uID uint, req *request.ReqReady) error
}
type IListRoomsRepository interface {
	FindRoomList(ctx context.Context, page int, pageSize int) ([]mysql.Rooms, error)
	CountRoomList(ctx context.Context) (int, error)
}
type IUserListRoomsRepository interface {
	FindRoomUser(ctx context.Context, RoomID uint) ([]response.User, error)
	FindOneRoom(ctx context.Context, RoomID uint) (mysql.Rooms, error)
}
