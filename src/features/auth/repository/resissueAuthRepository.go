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
		return utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), utils.HandleError(err.Error(), token), utils.ErrFromMysqlDB)
	}
	return nil
}

func (d *ReissueAuthRepository) DeleteToken(ctx context.Context, uID uint) error {
	err := d.GormDB.Where("user_id = ?", uID).Delete(&mysql.Tokens{}).Error
	if err != nil {
		return utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), utils.HandleError(err.Error(), uID), utils.ErrFromMysqlDB)
	}
	return nil
}

func (d *ReissueAuthRepository) CheckToken(ctx context.Context, uID uint, refreshToken string) error {
	var token mysql.Tokens
	err := d.GormDB.Where("user_id = ? and refresh_token = ?", uID, refreshToken).First(&token).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.ErrorMsg(ctx, utils.ErrBadToken, utils.Trace(), utils.HandleError(err.Error(), uID), utils.ErrFromClient)
		}
		return utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), utils.HandleError(err.Error(), uID), utils.ErrFromMysqlDB)
	}
	return nil
}
