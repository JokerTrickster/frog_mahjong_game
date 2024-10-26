package repository

import (
	"context"
	_interface "main/features/game/model/interface"
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewCreateMissionGameRepository(gormDB *gorm.DB) _interface.ICreateMissionGameRepository {
	return &CreateMissionGameRepository{GormDB: gormDB}
}

func (d *CreateMissionGameRepository) SaveMission(c context.Context, missionDTO *mysql.Missions) error {
	err := d.GormDB.Create(missionDTO).Error
	if err != nil {
		return utils.ErrorMsg(c, utils.ErrInternalDB, utils.Trace(), err.Error(), utils.ErrFromMysqlDB)
	}
	return nil
}
