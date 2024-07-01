package mysql

import "gorm.io/gorm"

type Tokens struct {
	gorm.Model
	UserID           uint   `json:"userID" gorm:"column:user_id"`
	AccessToken      string `json:"accessToken" gorm:"column:access_token"`
	RefreshToken     string `json:"refreshToken" gorm:"column:refresh_token"`
	RefreshExpiredAt int64  `json:"refreshExpiredAt" gorm:"column:refresh_expired_at"`
}

type Users struct {
	gorm.Model
	Name     string `json:"name" gorm:"column:name"`
	Email    string `json:"email" gorm:"uniqueIndex;column:email"`
	Password string `json:"password" gorm:"column:password"`
	Score    int    `json:"score" gorm:"column:score"`
	State    string `json:"state" gorm:"column:state"` //logout, wait, play
	RoomID   int    `json:"roomID" gorm:"column:room_id"`
	Provider string `json:"provider" gorm:"column:provider"`
}

type Rooms struct {
	gorm.Model
	CurrentCount int    `json:"currentCount" gorm:"column:current_count"`
	MaxCount     int    `json:"maxCount" gorm:"column:max_count"`
	MinCount     int    `json:"minCount" gorm:"column:min_count"`
	Name         string `json:"name" gorm:"column:name"`
	Password     string `json:"password" gorm:"column:password"`
	State        string `json:"state" gorm:"column:state"` //wait, play, end
	OwnerID      int    `json:"ownerID" gorm:"column:owner_id"`
}

type RoomUsers struct {
	gorm.Model
	UserID         int    `json:"userID" gorm:"column:user_id"`
	RoomID         int    `json:"roomID" gorm:"column:room_id"`
	Score          int    `json:"score" gorm:"column:score"`
	OwnedCardCount int    `json:"ownedCardCount" gorm:"column:owned_card_count"`
	PlayerState    string `json:"playerState" gorm:"column:player_state"` // wait(대기중), ready(준비 완료), play(플레이할 차례), play_wait(다음 차례 대기)
	TurnNumber     int    `json:"turnNumber" gorm:"column:turn_number"`   // 1, 2, 3, 4 ....
}

type Cards struct {
	gorm.Model
	RoomID int    `json:"roomID" gorm:"column:room_id"`
	UserID int    `json:"userID" gorm:"column:user_id"`
	Name   string `json:"name" gorm:"column:name"`   // 1, 2, 3, 4, 5, 6, 7, 8, 9, 중, 발
	Color  string `json:"color" gorm:"column:color"` // red, green, normal
	State  string `json:"state" gorm:"column:state"` // owned, discard, none
}

type Chat struct {
	gorm.Model
	UserID int    `json:"userID" gorm:"column:user_id"`
	Name   string `json:"name" gorm:"column:name"`
	Secret string `json:"secret" gorm:"column:secret"`
}
