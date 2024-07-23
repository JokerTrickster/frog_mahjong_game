package repository

import (
	"context"
	_interface "main/features/rooms/model/interface"
	"main/features/rooms/model/response"
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewUserListRoomsRepository(gormDB *gorm.DB) _interface.IUserListRoomsRepository {
	return &UserListRoomsRepository{GormDB: gormDB}
}

func (d *UserListRoomsRepository) FindRoomUser(ctx context.Context, RoomsID uint) ([]response.User, error) {
	RoomUserList := make([]response.User, 0)
	err := d.GormDB.Table("users").
		Joins("LEFT JOIN room_users ON users.id = room_users.user_id").
		Select("users.id AS user_id, room_users.id AS room_user_id, room_users.player_state, room_users.turn_number, room_users.owned_card_count, room_users.Room_id, room_users.score, users.name AS user_name, users.email AS user_email").
		Where("room_users.room_id = ?", RoomsID).
		Scan(&RoomUserList).Error
	if err != nil {
		return nil, utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), err.Error(), utils.ErrFromMysqlDB)
	}
	return RoomUserList, nil
}

func (d *UserListRoomsRepository) FindOneRoom(ctx context.Context, RoomID uint) (mysql.Rooms, error) {
	Room := mysql.Rooms{}
	err := d.GormDB.Table("rooms").Where("id = ?", RoomID).Scan(&Room).Error
	if err != nil {
		return mysql.Rooms{}, utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), err.Error(), utils.ErrFromMysqlDB)
	}
	return Room, nil
}
