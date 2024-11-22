package repository

import (
	_interface "main/features/game/model/interface"

	"gorm.io/gorm"
)

func NewReportImageUploadGameRepository(gormDB *gorm.DB) _interface.IReportImageUploadGameRepository {
	return &ReportImageUploadGameRepository{GormDB: gormDB}
}
