package repository

import (
	"context"
	_interface "main/features/room/model/interface"
	"main/features/room/model/response"
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewUserListRoomRepository(gormDB *gorm.DB) _interface.IUserListRoomRepository {
	return &UserListRoomRepository{GormDB: gormDB}
}

func (d *UserListRoomRepository) FindRoomUser(ctx context.Context, roomID uint) ([]response.User, error) {
	roomUserList := make([]response.User, 0)
	err := d.GormDB.Table("users").
		Joins("LEFT JOIN room_users ON users.id = room_users.user_id").
		Select("users.id AS user_id, room_users.id AS room_user_id, room_users.player_state, room_users.turn_number, room_users.owned_card_count, room_users.room_id, room_users.score, users.name AS user_name, users.email AS user_email").
		Where("room_users.room_id = ?", roomID).
		Scan(&roomUserList).Error
	if err != nil {
		return nil, utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), err.Error(), utils.ErrFromMysqlDB)
	}
	return roomUserList, nil
}

func (d *UserListRoomRepository) FindOneRoom(ctx context.Context, roomID uint) (mysql.Rooms, error) {
	room := mysql.Rooms{}
	err := d.GormDB.Table("rooms").Where("id = ?", roomID).Scan(&room).Error
	if err != nil {
		return mysql.Rooms{}, utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), err.Error(), utils.ErrFromMysqlDB)
	}
	return room, nil
}
