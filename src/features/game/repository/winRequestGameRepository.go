package repository

import (
	_interface "main/features/game/model/interface"

	"gorm.io/gorm"
)

func NewWinRequestGameRepository(gormDB *gorm.DB) _interface.IWinRequestGameRepository {
	return &WinRequestGameRepository{GormDB: gormDB}
}
