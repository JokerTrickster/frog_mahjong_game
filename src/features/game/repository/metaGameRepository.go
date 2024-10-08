package repository

import (
	"context"
	_interface "main/features/game/model/interface"
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewMetaGameRepository(gormDB *gorm.DB) _interface.IMetaGameRepository {
	return &MetaGameRepository{GormDB: gormDB}
}

func (d *MetaGameRepository) FindAllReportCategory(c context.Context) ([]mysql.Categories, error) {
	var categoryList []mysql.Categories
	err := d.GormDB.Model(&categoryList).Where("type = ?", "report").Find(&categoryList).Error
	if err != nil {
		return nil, utils.ErrorMsg(c, utils.ErrInternalDB, utils.Trace(), utils.HandleError(err.Error()), utils.ErrFromMysqlDB)
	}
	return categoryList, nil
}
