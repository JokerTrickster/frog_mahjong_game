package repository

import "gorm.io/gorm"

type CreateRoomsRepository struct {
	GormDB *gorm.DB
}
type V02CreateRoomsRepository struct {
	GormDB *gorm.DB
}

type JoinPlayRoomsRepository struct {
	GormDB *gorm.DB
}
type V02JoinRoomsRepository struct {
	GormDB *gorm.DB
}
type OutRoomsRepository struct {
	GormDB *gorm.DB
}

type ReadyRoomsRepository struct {
	GormDB *gorm.DB
}

type ListRoomsRepository struct {
	GormDB *gorm.DB
}

type UserListRoomsRepository struct {
	GormDB *gorm.DB
}

type MetaRoomsRepository struct {
	GormDB *gorm.DB
}

type CheckSessionRoomsRepository struct {
	GormDB *gorm.DB
}