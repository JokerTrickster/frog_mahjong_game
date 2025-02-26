package repository

import (
	"context"
	"main/features/game_users/model/entity"
	_interface "main/features/game_users/model/interface"
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewUpdateUsersRepository(gormDB *gorm.DB) _interface.IUpdateUsersRepository {
	return &UpdateUsersRepository{GormDB: gormDB}
}
func (d *UpdateUsersRepository) FindOneAndUpdateUsers(ctx context.Context, entitySQL *entity.UpdateUsersEntitySQL) error {
	err := d.GormDB.WithContext(ctx).Model(&mysql.GameUsers{}).Where("id = ?", entitySQL.UserID).Updates(entitySQL).Error
	if err != nil {
		return utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), utils.HandleError(err.Error()), utils.ErrFromMongoDB)
	}
	return nil
}
