package repository

import (
	"context"
	_interface "main/features/auth/model/interface"
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewReissueAuthRepository(gormDB *gorm.DB) _interface.IReissueAuthRepository {
	return &ReissueAuthRepository{GormDB: gormDB}
}

func (d *ReissueAuthRepository) SaveToken(ctx context.Context, token mysql.Tokens) error {
	err := d.GormDB.Create(&token).Error
	if err != nil {
		return utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), err.Error(), utils.ErrFromMysqlDB)
	}
	return nil
}

func (d *ReissueAuthRepository) DeleteToken(ctx context.Context, uID uint) error {
	err := d.GormDB.Where("user_id = ?", uID).Delete(&mysql.Tokens{}).Error
	if err != nil {
		return utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), err.Error(), utils.ErrFromMysqlDB)
	}
	return nil
}
