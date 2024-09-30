package repository

import (
	"context"
	_errors "main/features/rooms/model/errors"
	_interface "main/features/rooms/model/interface"
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewMetaRoomsRepository(gormDB *gorm.DB) _interface.IMetaRoomsRepository {
	return &MetaRoomsRepository{GormDB: gormDB}
}

func (g *MetaRoomsRepository) FindAllTimeMeta(ctx context.Context) ([]mysql.Times, error) {
	var timeDTO []mysql.Times
	if err := g.GormDB.WithContext(ctx).Find(&timeDTO).Error; err != nil {
		return nil, utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), utils.HandleError(_errors.ErrServerError.Error()+err.Error()), utils.ErrFromMysqlDB)
	}
	return timeDTO, nil
}
