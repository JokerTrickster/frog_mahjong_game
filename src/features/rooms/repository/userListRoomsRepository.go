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
		Joins("LEFT JOIN Room_users ON users.id = Room_users.user_id").
		Select("users.id AS user_id, Room_users.id AS Room_user_id, Room_users.player_state, Room_users.turn_number, Room_users.owned_card_count, Room_users.Room_id, Room_users.score, users.name AS user_name, users.email AS user_email").
		Where("Room_users.Rooms_id = ?", RoomsID).
		Scan(&RoomUserList).Error
	if err != nil {
		return nil, utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), err.Error(), utils.ErrFromMysqlDB)
	}
	return RoomUserList, nil
}

func (d *UserListRoomsRepository) FindOneRoom(ctx context.Context, RoomID uint) (mysql.Rooms, error) {
	Room := mysql.Rooms{}
	err := d.GormDB.Table("Rooms").Where("id = ?", RoomID).Scan(&Room).Error
	if err != nil {
		return mysql.Rooms{}, utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), err.Error(), utils.ErrFromMysqlDB)
	}
	return Room, nil
}
