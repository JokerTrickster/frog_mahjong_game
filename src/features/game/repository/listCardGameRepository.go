package repository

import (
	"context"
	_interface "main/features/game/model/interface"
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewListCardGameRepository(gormDB *gorm.DB) _interface.IListCardGameRepository {
	return &ListCardGameRepository{GormDB: gormDB}
}

func (d *ListCardGameRepository) FindFrogCard(c context.Context) ([]*mysql.FrogCards, error) {
	var cards []*mysql.FrogCards
	err := d.GormDB.Model(&cards).Find(&cards).Error
	if err != nil {
		return nil, utils.ErrorMsg(c, utils.ErrInternalServer, utils.Trace(), utils.HandleError(err.Error()), utils.ErrFromMysqlDB)
	}
	return cards, nil
}

func (d *ListCardGameRepository) CountFrogCard(c context.Context) (int, error) {
	var count int64
	err := d.GormDB.Model(&mysql.FrogCards{}).Count(&count).Error
	if err != nil {
		return 0, utils.ErrorMsg(c, utils.ErrInternalServer, utils.Trace(), utils.HandleError(err.Error()), utils.ErrFromMysqlDB)
	}
	return int(count), nil
}
