package mysql

import "gorm.io/gorm"

// 전체, 한식, 중식, 일식, 양식, 분식, 패스트푸드, 카페, 술집, 기타
type Times struct {
	gorm.Model
	Timer       uint   `json:"timer" gorm:"column:timer"`
	Description string `json:"description" gorm:"column:description"`
}
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
	Coin     int    `json:"coin" gorm:"column:coin"`
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
	PlayTurn     int    `json:"playTurn" gorm:"column:play_turn"`
	TimeOut      int    `json:"timeOut" gorm:"column:time_out"`
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
	CardID int    `json:"cardID" gorm:"column:card_id"` // 1 ~  44
	Name   string `json:"name" gorm:"column:name"`      // one, two, three, four, five .... nine, chung, bal
	Color  string `json:"color" gorm:"column:color"`    // red, green, normal
	State  string `json:"state" gorm:"column:state"`    // owned, discard, none
}

type Chats struct {
	gorm.Model
	UserID  int    `json:"userID" gorm:"column:user_id"`
	RoomID  int    `json:"roomID" gorm:"column:room_id"`
	Name    string `json:"name" gorm:"column:name"`
	Message string `json:"message" gorm:"column:message"`
}

type UserAuths struct {
	gorm.Model
	Email    string `json:"email" gorm:"column:email"`
	AuthCode string `json:"authCode" gorm:"column:auth_code"`
	Type     string `json:"type" gorm:"column:type"`
}
