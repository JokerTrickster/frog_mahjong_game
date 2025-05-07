package repository

import (
	"context"
	_interface "main/features/board_game/model/interface"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewGameOverBoardGameRepository(gormDB *gorm.DB) _interface.IGameOverBoardGameRepository {
	return &GameOverBoardGameRepository{GormDB: gormDB}
}
func (d *GameOverBoardGameRepository) GameOverInsertGameResult(ctx context.Context, gameResultDTO *mysql.GameResults) error {
	err := d.GormDB.WithContext(ctx).Create(gameResultDTO).Error
	if err != nil {
		return err
	}
	return nil
}
