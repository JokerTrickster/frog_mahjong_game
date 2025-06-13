package repository

import (
	"context"
	_interface "main/features/game_profiles/model/interface"
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewUpdateProfilesRepository(gormDB *gorm.DB) _interface.IUpdateProfilesRepository {
	return &UpdateProfilesRepository{GormDB: gormDB}
}

func (d *UpdateProfilesRepository) UpdateOneProfile(ctx context.Context, userID int, profileID int) error {
	err := d.GormDB.WithContext(ctx).Model(&mysql.GameUsers{}).Where("user_id = ?", userID).Update("profile_id", profileID).Error
	if err != nil {
		return utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), utils.HandleError(err.Error(), profileID), utils.ErrFromMysqlDB)
	}
	return nil
}
