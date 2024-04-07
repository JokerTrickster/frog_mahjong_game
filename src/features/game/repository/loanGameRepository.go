package repository

import (
	_interface "main/features/game/model/interface"

	"gorm.io/gorm"
)

func NewLoanGameRepository(gormDB *gorm.DB) _interface.ILoanGameRepository {
	return &LoanGameRepository{GormDB: gormDB}
}
