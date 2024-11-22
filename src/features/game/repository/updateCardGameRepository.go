package repository

import (
	"context"
	_interface "main/features/game/model/interface"
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewUpdateCardGameRepository(gormDB *gorm.DB) _interface.IUpdateCardGameRepository {
	return &UpdateCardGameRepository{GormDB: gormDB}
}

func (d *UpdateCardGameRepository) UpdateCard(c context.Context, birdCardsDTO mysql.BirdCards) error {
	err := d.GormDB.WithContext(c).Model(&mysql.BirdCards{}).Where("id = ?", birdCardsDTO.ID).Updates(&birdCardsDTO).Error
	if err != nil {
		return utils.ErrorMsg(c, utils.ErrInternalDB, utils.Trace(), err.Error(), utils.ErrFromMysqlDB)
	}
	return nil
}
