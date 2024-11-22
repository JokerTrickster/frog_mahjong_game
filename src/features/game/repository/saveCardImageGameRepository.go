package repository

import (
	"context"
	_interface "main/features/game/model/interface"
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewSaveCardImageGameRepository(gormDB *gorm.DB) _interface.ISaveCardImageGameRepository {
	return &SaveCardImageGameRepository{GormDB: gormDB}
}

func (d *SaveCardImageGameRepository) FindOneUpdateCardImage(c context.Context, cardName, newFileName string) error {
	err := d.GormDB.Model(&mysql.BirdCards{}).Where("image = ?", cardName).Update("image", newFileName).Error
	if err != nil {
		return utils.ErrorMsg(c, utils.ErrInternalDB, utils.Trace(), err.Error(), utils.ErrFromMysqlDB)
	}
	return nil
}
