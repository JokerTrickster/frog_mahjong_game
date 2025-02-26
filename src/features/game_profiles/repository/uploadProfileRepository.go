package repository

import (
	"context"
	_interface "main/features/game_profiles/model/interface"
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewUploadProfilesRepository(gormDB *gorm.DB) _interface.IUploadProfilesRepository {
	return &UploadProfilesRepository{GormDB: gormDB}
}
func (g *UploadProfilesRepository) InsertOneProfile(ctx context.Context, profile *mysql.GameProfiles) error {
	err := g.GormDB.WithContext(ctx).Create(profile).Error
	if err != nil {
		return utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), utils.HandleError(err.Error(), profile), utils.ErrFromMysqlDB)
	}
	return nil
}
