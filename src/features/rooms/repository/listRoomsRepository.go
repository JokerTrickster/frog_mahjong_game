package repository

import (
	"context"
	_interface "main/features/rooms/model/interface"
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewListRoomsRepository(gormDB *gorm.DB) _interface.IListRoomsRepository {
	return &ListRoomsRepository{GormDB: gormDB}
}

func (d *ListRoomsRepository) FindRoomList(ctx context.Context, page int, pageSize int) ([]mysql.Rooms, error) {
	var Rooms []mysql.Rooms
	// 페이지네이션 처리
	err := d.GormDB.Where("deleted_at IS NULL and id != 1").Limit(pageSize).Offset((page - 1) * pageSize).Find(&Rooms).Error
	if err != nil {
		return nil, utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), err.Error(), utils.ErrFromMysqlDB)
	}

	return Rooms, nil
}

func (d *ListRoomsRepository) CountRoomList(ctx context.Context) (int, error) {
	var count int64
	err := d.GormDB.Model(&mysql.Rooms{}).Where("deleted_at IS NULL").Count(&count).Error
	if err != nil {
		return 0, utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), err.Error(), utils.ErrFromMysqlDB)
	}

	return int(count), nil
}
