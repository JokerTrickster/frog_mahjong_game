package repository

import (
	_interface "main/features/game/model/interface"

	"gorm.io/gorm"
)

func NewScoreCalculateGameRepository(gormDB *gorm.DB) _interface.IScoreCalculateGameRepository {
	return &ScoreCalculateGameRepository{GormDB: gormDB}
}
