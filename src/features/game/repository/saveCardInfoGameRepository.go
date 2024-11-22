package repository

import (
	"context"
	_interface "main/features/game/model/interface"
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewSaveCardInfoGameRepository(gormDB *gorm.DB) _interface.ISaveCardInfoGameRepository {
	return &SaveCardInfoGameRepository{GormDB: gormDB}
}

func (d *SaveCardInfoGameRepository) SaveCardInfo(c context.Context, birdCardsDTO []mysql.BirdCards) error {
	err := d.GormDB.Create(birdCardsDTO).Error
	if err != nil {
		return utils.ErrorMsg(c, utils.ErrInternalDB, utils.Trace(), err.Error(), utils.ErrFromMysqlDB)
	}
	return nil
}
