package repository

import (
	"context"
	_interface "main/features/game_users/model/interface"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewListProfilesUsersRepository(gormDB *gorm.DB) _interface.IListProfilesUsersRepository {
	return &ListProfilesUsersRepository{GormDB: gormDB}
}

func (d *ListProfilesUsersRepository) FindAllProfiles(ctx context.Context, userID uint) ([]*mysql.GameUserProfiles, error) {
	profiles := make([]*mysql.GameUserProfiles, 0)
	err := d.GormDB.WithContext(ctx).Where("user_id = ?", userID).Find(&profiles).Error
	if err != nil {
		return nil, err
	}
	return profiles, nil
}
