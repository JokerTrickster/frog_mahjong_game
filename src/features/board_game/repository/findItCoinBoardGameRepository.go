package repository

import (
	"context"
	_interface "main/features/board_game/model/interface"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewFindItCoinBoardGameRepository(gormDB *gorm.DB) _interface.IFindItCoinBoardGameRepository {
	return &FindItCoinBoardGameRepository{GormDB: gormDB}
}
func (d *FindItCoinBoardGameRepository) UpdateUserCoin(ctx context.Context, userID int, coin int) error {
	//조회한다.
	gameUserDTO := mysql.GameUsers{}
	err := d.GormDB.WithContext(ctx).Where("id = ?", userID).First(&gameUserDTO).Error
	if err != nil {
		return err
	}
	tmpCoin := gameUserDTO.Coin
	tmpCoin += coin
	if tmpCoin <= 0 {
		gameUserDTO.Coin = 0
	} else {
		gameUserDTO.Coin = tmpCoin
	}
	//업데이트한다.
	err = d.GormDB.WithContext(ctx).Model(&gameUserDTO).
		UpdateColumn("coin", gameUserDTO.Coin).Error
	if err != nil {
		return err
	}
	return nil
}
