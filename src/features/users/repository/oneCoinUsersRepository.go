package repository

import (
	"context"
	_interface "main/features/users/model/interface"
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewOneCoinUsersRepository(gormDB *gorm.DB) _interface.IOneCoinUsersRepository {
	return &OneCoinUsersRepository{GormDB: gormDB}
}
func (d *OneCoinUsersRepository) OneCoin(ctx context.Context) error {
	// coin이 30 미만인 유저들을 코인 +1 한다
	err := d.GormDB.WithContext(ctx).Model(&mysql.Users{}).Where("coin < 30").Update("coin", gorm.Expr("coin + 1")).Error
	if err != nil {
		return utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), utils.HandleError(err.Error()), utils.ErrFromMongoDB)
	}

	return nil
}
