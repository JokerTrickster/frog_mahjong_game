package repository

import "gorm.io/gorm"

type FindItSoloPlayBoardGameRepository struct {
	GormDB *gorm.DB
}
type FindItRankBoardGameRepository struct {
	GormDB *gorm.DB
}

type FindItCoinBoardGameRepository struct {
	GormDB *gorm.DB
}

type FindItPasswordCheckBoardGameRepository struct {
	GormDB *gorm.DB
}

type SlimeWarGetsCardBoardGameRepository struct {
	GormDB *gorm.DB
}

type SlimeWarResultBoardGameRepository struct {
	GormDB *gorm.DB
}
