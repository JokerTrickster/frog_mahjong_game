package repository

import (
	"context"
	_interface "main/features/users/model/interface"
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewFullCoinUsersRepository(gormDB *gorm.DB) _interface.IFullCoinUsersRepository {
	return &FullCoinUsersRepository{GormDB: gormDB}
}
func (d *FullCoinUsersRepository) FullCoin(ctx context.Context) error {
	// coin이 30 미만인 유저들을 모두 30으로 만든다.
	err := d.GormDB.WithContext(ctx).Model(&mysql.Users{}).Where("coin < 30").Update("coin", 30).Error
	if err != nil {
		return utils.ErrorMsg(ctx, utils.ErrInternalDB, utils.Trace(), utils.HandleError(err.Error()), utils.ErrFromMongoDB)
	}

	return nil
}
