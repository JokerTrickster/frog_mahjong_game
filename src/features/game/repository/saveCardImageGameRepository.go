package repository

import (
	_interface "main/features/game/model/interface"

	"gorm.io/gorm"
)

func NewSaveCardImageGameRepository(gormDB *gorm.DB) _interface.ISaveCardImageGameRepository {
	return &SaveCardImageGameRepository{GormDB: gormDB}
}
