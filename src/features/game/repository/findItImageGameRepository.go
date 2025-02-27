package repository

import (
	_interface "main/features/game/model/interface"

	"gorm.io/gorm"
)

func NewFindItImageGameRepository(gormDB *gorm.DB) _interface.IFindItImageGameRepository {
	return &FindItImageGameRepository{GormDB: gormDB}
}
