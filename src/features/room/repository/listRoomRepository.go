package repository

import (
	"context"
	_interface "main/features/room/model/interface"
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewListRoomRepository(gormDB *gorm.DB) _interface.IListRoomRepository {
	return &ListRoomRepository{GormDB: gormDB}
}

func (d *ListRoomRepository) FindRoomList(ctx context.Context, page int, pageSize int) ([]mysql.Rooms, error) {
	var rooms []mysql.Rooms
	// 페이지네이션 처리
	err := d.GormDB.Where("deleted_at IS NULL").Limit(pageSize).Offset((page - 1) * pageSize).Find(&rooms).Error
	if err != nil {
		return nil, utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), err.Error(), utils.ErrFromMysqlDB)
	}

	return rooms, nil
}

func (d *ListRoomRepository) CountRoomList(ctx context.Context) (int, error) {
	var count int64
	err := d.GormDB.Model(&mysql.Rooms{}).Where("deleted_at IS NULL").Count(&count).Error
	if err != nil {
		return 0, utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), err.Error(), utils.ErrFromMysqlDB)
	}

	return int(count), nil
}
