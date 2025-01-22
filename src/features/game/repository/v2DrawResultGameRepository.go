package repository

import (
	"context"
	_interface "main/features/game/model/interface"
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewV2DrawResultGameRepository(gormDB *gorm.DB) _interface.IV2DrawResultGameRepository {
	return &V2DrawResultGameRepository{GormDB: gormDB}
}

func (d *V2DrawResultGameRepository) FindAllUserMission(c context.Context, userID, roomID int) ([]mysql.UserMissions, error) {
	var userMissions []mysql.UserMissions
	err := d.GormDB.Model(&userMissions).Where("user_id = ? and room_id = ?", userID, roomID).Find(&userMissions).Error
	if err != nil {
		return nil, utils.ErrorMsg(c, utils.ErrInternalDB, utils.Trace(), utils.HandleError(err.Error(), roomID), utils.ErrFromMysqlDB)
	}
	return userMissions, nil
}
func (d *V2DrawResultGameRepository) FindAllRoomUsers(c context.Context, roomID int) ([]mysql.RoomUsers, error) {
	var roomUsers []mysql.RoomUsers
	err := d.GormDB.Model(&roomUsers).Where("room_id = ?", roomID).Find(&roomUsers).Error
	if err != nil {
		return nil, utils.ErrorMsg(c, utils.ErrInternalDB, utils.Trace(), utils.HandleError(err.Error(), roomID), utils.ErrFromMysqlDB)
	}
	return roomUsers, nil
}
