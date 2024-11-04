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

func (d *ListCardGameRepository) FindAllBirdCard(c context.Context) ([]*mysql.BirdCards, error) {
	var cards []*mysql.BirdCards
	err := d.GormDB.Model(&cards).Find(&cards).Error
	if err != nil {
		return nil, utils.ErrorMsg(c, utils.ErrInternalServer, utils.Trace(), utils.HandleError(err.Error()), utils.ErrFromMysqlDB)
	}
	return cards, nil
}
