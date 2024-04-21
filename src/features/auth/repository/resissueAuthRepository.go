package repository

import (
	_interface "main/features/auth/model/interface"

	"gorm.io/gorm"
)

func NewReissueAuthRepository(gormDB *gorm.DB) _interface.IReissueAuthRepository {
	return &ReissueAuthRepository{GormDB: gormDB}
}
