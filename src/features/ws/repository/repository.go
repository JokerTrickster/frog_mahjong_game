package repository

import "gorm.io/gorm"

type JoinEventWebsocketRepository struct {
	GormDB *gorm.DB
}
