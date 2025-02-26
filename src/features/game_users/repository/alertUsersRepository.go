package repository

import (
	"context"
	_interface "main/features/game_users/model/interface"
	"main/features/game_users/model/request"
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewAlertUsersRepository(gormDB *gorm.DB) _interface.IAlertUsersRepository {
	return &AlertUsersRepository{GormDB: gormDB}
}

func (d *AlertUsersRepository) FindOneAndUpdateUsers(ctx context.Context, userID uint, req *request.ReqAlertGameUsers) error {
	err := d.GormDB.WithContext(ctx).Model(&mysql.GameUsers{}).Where("id = ?", userID).Update("alert_enabled", req.Alert).Error
	if err != nil {
		return utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), utils.HandleError(err.Error()), utils.ErrFromMongoDB)
	}
	return nil
}
