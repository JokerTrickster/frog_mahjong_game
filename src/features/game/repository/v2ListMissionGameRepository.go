package repository

import (
	"context"
	_interface "main/features/game/model/interface"
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewListMissionGameRepository(gormDB *gorm.DB) _interface.IListMissionGameRepository {
	return &ListMissionGameRepository{GormDB: gormDB}
}

func (d *ListMissionGameRepository) FindAllMission(c context.Context) ([]*mysql.Missions, error) {
	var missions []*mysql.Missions
	if err := d.GormDB.Find(&missions).Error; err != nil {
		return nil, utils.ErrorMsg(c, utils.ErrInternalDB, utils.Trace(), err.Error(), utils.ErrFromMysqlDB)
	}
	return missions, nil
}
