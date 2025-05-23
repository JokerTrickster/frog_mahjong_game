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

type SlimeWarRankBoardGameRepository struct {
	GormDB *gorm.DB
}

type SequenceRankBoardGameRepository struct {
	GormDB *gorm.DB
}

type SequenceResultBoardGameRepository struct {
	GormDB *gorm.DB
}

type GameOverBoardGameRepository struct {
	GormDB *gorm.DB
}

type FrogCardListBoardGameRepository struct {
	GormDB *gorm.DB
}
