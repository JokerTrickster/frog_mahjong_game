package repository

import (
	"context"
	_interface "main/features/game/model/interface"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewFindItImageInfoGameRepository(gormDB *gorm.DB) _interface.IFindItImageInfoGameRepository {
	return &FindItImageInfoGameRepository{GormDB: gormDB}
}

func (d *FindItImageInfoGameRepository) SaveImageInfo(c context.Context, imageDTO *mysql.FindItImages) (int, error) {
	err := d.GormDB.Create(imageDTO).Error
	if err != nil {
		return 0, err
	}
	return int(imageDTO.ID), nil
}
func (d *FindItImageInfoGameRepository) SaveImageCorrectInfo(c context.Context, imageCorrectDTOs []*mysql.FindItImageCorrectPositions) error {
	for _, imageCorrectDTO := range imageCorrectDTOs {
		err := d.GormDB.Create(imageCorrectDTO).Error
		if err != nil {
			return err
		}
	}
	return nil
}

