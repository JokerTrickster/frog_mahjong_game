package repository

import (
	"context"
	_interface "main/features/board_game/model/interface"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewFindItPasswordCheckBoardGameRepository(gormDB *gorm.DB) _interface.IFindItPasswordCheckBoardGameRepository {
	return &FindItPasswordCheckBoardGameRepository{GormDB: gormDB}
}

func (d *FindItPasswordCheckBoardGameRepository) FindPasswordCheck(ctx context.Context, password string) (bool, error) {
	gameRoom := &mysql.GameRooms{}
	err := d.GormDB.WithContext(ctx).Model(&mysql.GameRooms{}).Where("password = ?", password).First(&gameRoom).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
