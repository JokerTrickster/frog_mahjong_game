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
