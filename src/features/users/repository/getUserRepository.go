package repository

import (
	_interface "main/features/users/model/interface"

	"gorm.io/gorm"
)

func NewGetUsersRepository(gormDB *gorm.DB) _interface.IGetUsersRepository {
	return &GetUsersRepository{GormDB: gormDB}
}

