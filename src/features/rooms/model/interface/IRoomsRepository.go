package _interface

import (
	"context"
	"main/features/rooms/model/request"
	"main/features/rooms/model/response"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

type ICreateRoomsRepository interface {
	InsertOneRoom(ctx context.Context, tx *gorm.DB, RoomsDTO mysql.Rooms) (int, error)
	InsertOneRoomUser(ctx context.Context, tx *gorm.DB, RoomsUserDTO mysql.RoomUsers) error
	FindOneAndUpdateUser(ctx context.Context, tx *gorm.DB, uID uint, RoomID uint) error
}
type IV02CreateRoomsRepository interface {
	InsertOneRoom(ctx context.Context, tx *gorm.DB, RoomsDTO mysql.Rooms) (int, error)
	InsertOneRoomUser(ctx context.Context, tx *gorm.DB, RoomsUserDTO mysql.RoomUsers) error
	FindOneAndUpdateUser(ctx context.Context, tx *gorm.DB, uID uint, RoomID uint) error
}

type IJoinPlayRoomsRepository interface {
	FindOneRoom(ctx context.Context, req *request.ReqJoinPlay) error
}
type IV02JoinRoomsRepository interface {
}

type IOutRoomsRepository interface {
	FindOneAndDeleteRoomUser(ctx context.Context, tx *gorm.DB, uID uint, RoomID uint) error
	FindOneAndUpdateRoom(ctx context.Context, tx *gorm.DB, RoomID uint) (mysql.Rooms, error)
	FindOneAndUpdateUser(ctx context.Context, tx *gorm.DB, uID uint) error
	FindOneAndDeleteRoom(ctx context.Context, tx *gorm.DB, RoomID uint) error
	FindOneRoomUser(ctx context.Context, tx *gorm.DB, RoomID uint) (mysql.RoomUsers, error)
	ChangeRoomOnwer(ctx context.Context, tx *gorm.DB, RoomID uint, ownerID uint) error
	FindOneUser(ctx context.Context, tx *gorm.DB, uID uint) (mysql.Users, error)
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
type IMetaRoomsRepository interface {
	FindAllTimeMeta(ctx context.Context) ([]mysql.Times, error)
}
