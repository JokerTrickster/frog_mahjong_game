package repository

import "gorm.io/gorm"

type StartGameRepository struct {
	GormDB *gorm.DB
}

type DoraGameRepository struct {
	GormDB *gorm.DB
}

type OwnershipGameRepository struct {
	GormDB *gorm.DB
}

type DiscardGameRepository struct {
	GormDB *gorm.DB
}

type NextTurnGameRepository struct {
	GormDB *gorm.DB
}

type LoanGameRepository struct {
	GormDB *gorm.DB
}

type ScoreCalculateGameRepository struct {
	GormDB *gorm.DB
}

type WinRequestGameRepository struct {
	GormDB *gorm.DB
}

type ResultGameRepository struct {
	GormDB *gorm.DB
}

type ReportGameRepository struct {
	GormDB *gorm.DB
}

type MetaGameRepository struct {
	GormDB *gorm.DB
}

type DeckCardGameRepository struct {
	GormDB *gorm.DB
}

type ListMissionGameRepository struct {
	GormDB *gorm.DB
}

type CreateMissionGameRepository struct {
	GormDB *gorm.DB
}

type ListCardGameRepository struct {
	GormDB *gorm.DB
}

// v2
type V2DeckCardGameRepository struct {
	GormDB *gorm.DB
}

type V2ReportGameRepository struct {
	GormDB *gorm.DB
}
type V2ResultGameRepository struct {
	GormDB *gorm.DB
}
