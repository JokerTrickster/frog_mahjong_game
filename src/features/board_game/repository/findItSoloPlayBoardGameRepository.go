package repository

import (
	"context"
	_interface "main/features/board_game/model/interface"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewFindItSoloPlayBoardGameRepository(gormDB *gorm.DB) _interface.IFindItSoloPlayBoardGameRepository {
	return &FindItSoloPlayBoardGameRepository{GormDB: gormDB}
}

func (d *FindItSoloPlayBoardGameRepository) FindRandomImage(ctx context.Context, round int) ([]*mysql.FindItImages, error) {
	var images []*mysql.FindItImages
	err := d.GormDB.WithContext(ctx).
		Order("RAND()").
		Limit(round).
		Find(&images).Error
	if err != nil {
		return nil, err
	}
	return images, nil
}

func (d *FindItSoloPlayBoardGameRepository) FindCorrectByImageID(ctx context.Context, imageID uint) ([]*mysql.FindItImageCorrectPositions, error) {
	var corrects []*mysql.FindItImageCorrectPositions
	err := d.GormDB.WithContext(ctx).Where("image_id = ?", imageID).Find(&corrects).Error
	if err != nil {
		return nil, err
	}
	return corrects, nil
}
