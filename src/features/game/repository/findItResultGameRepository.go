package repository

import (
	"context"
	_interface "main/features/game/model/interface"
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewFindItResultGameRepository(gormDB *gorm.DB) _interface.IFindItResultGameRepository {
	return &FindItResultGameRepository{GormDB: gormDB}
}

func (d *FindItResultGameRepository) FindOneRoomSetting(ctx context.Context, roomID int) (*mysql.FindItRoomSettings, error) {
	var roomSetting mysql.FindItRoomSettings
	err := d.GormDB.Model(&roomSetting).Where("room_id = ?", roomID).First(&roomSetting).Error
	if err != nil {
		return nil, utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), utils.HandleError(err.Error(), roomID), utils.ErrFromMysqlDB)
	}
	return &roomSetting, nil
}

func (d *FindItResultGameRepository) FindFindItUserCorrectPositions(ctx context.Context, roomID int) ([]*mysql.FindItUserCorrectPositions, error) {
	var userCorrectPositions []*mysql.FindItUserCorrectPositions
	err := d.GormDB.Model(&userCorrectPositions).Where("room_id = ?", roomID).Find(&userCorrectPositions).Error
	if err != nil {
		return nil, utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), utils.HandleError(err.Error(), roomID), utils.ErrFromMysqlDB)
	}
	return userCorrectPositions, nil
}
