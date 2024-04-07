package repository

import (
	_interface "main/features/game/model/interface"

	"gorm.io/gorm"
)

func NewNextTurnGameRepository(gormDB *gorm.DB) _interface.INextTurnGameRepository {
	return &NextTurnGameRepository{GormDB: gormDB}
}
