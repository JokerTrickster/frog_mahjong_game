package repository

import (
	"context"
	_interface "main/features/game/model/interface"
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewListGameRepository(gormDB *gorm.DB) _interface.IListGameRepository {
	return &ListGameRepository{GormDB: gormDB}
}

func (d *ListGameRepository) FindGame(c context.Context) ([]*mysql.Games, error) {
	var games []*mysql.Games
	err := d.GormDB.Model(&games).Find(&games).Error
	if err != nil {
		return nil, utils.ErrorMsg(c, utils.ErrInternalServer, utils.Trace(), utils.HandleError(err.Error()), utils.ErrFromMysqlDB)
	}
	return games, nil
}
