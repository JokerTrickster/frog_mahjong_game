package repository

import (
	"context"
	_interface "main/features/profiles/model/interface"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewListProfilesRepository(gormDB *gorm.DB) _interface.IListProfilesRepository {
	return &ListProfilesRepository{GormDB: gormDB}
}

func (d *ListProfilesRepository) FindAllProfiles(ctx context.Context) ([]*mysql.Profiles, error) {
	profiles := make([]*mysql.Profiles, 0)
	err := d.GormDB.WithContext(ctx).Find(&profiles).Error
	if err != nil {
		return nil, err
	}
	return profiles, nil
}
