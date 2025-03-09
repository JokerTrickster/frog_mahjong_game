package repository

import "gorm.io/gorm"

type FindItSoloPlayBoardGameRepository struct {
	GormDB *gorm.DB
}
type FindItRankBoardGameRepository struct {
	GormDB *gorm.DB
}
