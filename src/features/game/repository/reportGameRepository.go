package repository

import (
	"context"
	_interface "main/features/game/model/interface"
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewReportGameRepository(gormDB *gorm.DB) _interface.IReportGameRepository {
	return &ReportGameRepository{GormDB: gormDB}
}

func (d *ReportGameRepository) SaveReport(c context.Context, reportDTO *mysql.Reports) error {
	err := d.GormDB.Create(reportDTO).Error
	if err != nil {
		return utils.ErrorMsg(c, utils.ErrInternalDB, utils.Trace(), err.Error(), utils.ErrFromMysqlDB)
	}
	return nil
}
