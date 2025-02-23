package repository

import (
	"context"
	"fmt"
	_interface "main/features/game_auth/model/interface"
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewNameCheckAuthRepository(gormDB *gorm.DB) _interface.INameCheckAuthRepository {
	return &NameCheckAuthRepository{GormDB: gormDB}
}
func (d *NameCheckAuthRepository) CheckName(ctx context.Context, name string) error {
	var user mysql.GameUsers
	err := d.GormDB.WithContext(ctx).Model(&user).Where("name = ?", name).First(&user).Error
	fmt.Println(err)
	fmt.Println(user)
	if err == nil {
		return utils.ErrorMsg(ctx, utils.ErrNameAlreadyExist, utils.Trace(), utils.HandleError(name), utils.ErrFromClient)
	}
	return nil
}
