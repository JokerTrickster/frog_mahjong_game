package repository

import (
	"context"
	_interface "main/features/board_game/model/interface"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewSlimeWarGetsCardBoardGameRepository(gormDB *gorm.DB) _interface.ISlimeWarGetsCardBoardGameRepository {
	return &SlimeWarGetsCardBoardGameRepository{GormDB: gormDB}
}

func (d *SlimeWarGetsCardBoardGameRepository) FindCardList(ctx context.Context) ([]*mysql.SlimeWarCards, error) {
	var cardList []*mysql.SlimeWarCards
	err := d.GormDB.WithContext(ctx).Find(&cardList).Error
	if err != nil {
		return nil, err
	}
	return cardList, nil
}
