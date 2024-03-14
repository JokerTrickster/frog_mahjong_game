package mysql

import "gorm.io/gorm"

type Users struct {
	gorm.Model
	Name     string `json:"name" gorm:"column:name"`
	Email    string `json:"email" gorm:"uniqueIndex;column:email"`
	Password string `json:"password" gorm:"column:password"`
	Score    int    `json:"score" gorm:"column:score"`
	State    string `json:"state" gorm:"column:state"`
	RoomID   int    `json:"roomID" gorm:"column:room_id"`
}

type Rooms struct {
	gorm.Model
	CurrentCount int    `json:"currentCount" gorm:"column:current_count"`
	MaxCount     int    `json:"maxCount" gorm:"column:max_count"`
	MinCount     int    `json:"minCount" gorm:"column:min_count"`
	Name         string `json:"name" gorm:"column:name"`
	Password     string `json:"password" gorm:"column:password"`
	State        string `json:"state" gorm:"column:state"`
	Owner        string `json:"owner" gorm:"column:owner"`
}

type RoomUsers struct {
	gorm.Model
	UserID      int    `json:"userID" gorm:"column:user_id"`
	RoomID      int    `json:"roomID" gorm:"column:room_id"`
	Score       int    `json:"score" gorm:"column:score"`
	CardCount   int    `json:"cardCount" gorm:"column:card_count"`
	PlayerState string `json:"playerState" gorm:"column:player_state"`
}

type Cards struct {
	gorm.Model
	RoomID int    `json:"roomID" gorm:"column:room_id"`
	UserID int    `json:"userID" gorm:"column:user_id"`
	Name   string `json:"name" gorm:"column:name"`
	Color  string `json:"color" gorm:"column:color"`
	State  string `json:"state" gorm:"column:state"`
}
