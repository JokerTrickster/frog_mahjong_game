package repository

import (
	"context"
	_interface "main/features/game/model/interface"
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewV2ReportGameRepository(gormDB *gorm.DB) _interface.IV2ReportGameRepository {
	return &V2ReportGameRepository{GormDB: gormDB}
}

func (d *V2ReportGameRepository) SaveReport(c context.Context, V2ReportDTO *mysql.Reports) error {
	err := d.GormDB.Create(V2ReportDTO).Error
	if err != nil {
		return utils.ErrorMsg(c, utils.ErrInternalDB, utils.Trace(), err.Error(), utils.ErrFromMysqlDB)
	}
	return nil
}
