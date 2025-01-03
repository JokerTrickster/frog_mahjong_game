package repository

import (
	"context"
	_interface "main/features/game/model/interface"
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewV2ListCardGameRepository(gormDB *gorm.DB) _interface.IV2ListCardGameRepository {
	return &V2ListCardGameRepository{GormDB: gormDB}
}

func (d *V2ListCardGameRepository) FindAllBirdCard(c context.Context) ([]*mysql.BirdCards, error) {
	var cards []*mysql.BirdCards
	err := d.GormDB.Model(&cards).Find(&cards).Error
	if err != nil {
		return nil, utils.ErrorMsg(c, utils.ErrInternalServer, utils.Trace(), utils.HandleError(err.Error()), utils.ErrFromMysqlDB)
	}
	return cards, nil
}

func (d *V2ListCardGameRepository) CountAllBirdCard(c context.Context) (int, error) {
	var count int64
	err := d.GormDB.Model(&mysql.BirdCards{}).Count(&count).Error
	if err != nil {
		return 0, utils.ErrorMsg(c, utils.ErrInternalServer, utils.Trace(), utils.HandleError(err.Error()), utils.ErrFromMysqlDB)
	}
	return int(count), nil
}
