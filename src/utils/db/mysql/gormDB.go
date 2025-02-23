package mysql

import (
	"time"

	"gorm.io/gorm"
)

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
	Name         string `json:"name" gorm:"column:name"`
	Email        string `json:"email" gorm:"uniqueIndex;column:email"`
	Password     string `json:"password" gorm:"column:password"`
	Coin         int    `json:"coin" gorm:"column:coin"`
	State        string `json:"state" gorm:"column:state"` //logout, wait, play
	RoomID       int    `json:"roomID" gorm:"column:room_id"`
	Provider     string `json:"provider" gorm:"column:provider"`
	ProfileID    int    `json:"profileID" gorm:"column:profile_id"`
	AlertEnabled bool   `json:"alertEnabled" gorm:"column:alert_enabled"`
}

type Rooms struct {
	gorm.Model
	CurrentCount int       `json:"currentCount" gorm:"column:current_count"`
	MaxCount     int       `json:"maxCount" gorm:"column:max_count"`
	MinCount     int       `json:"minCount" gorm:"column:min_count"`
	Name         string    `json:"name" gorm:"column:name"`
	Password     string    `json:"password" gorm:"column:password"`
	State        string    `json:"state" gorm:"column:state"` //wait, play, end
	OwnerID      int       `json:"ownerID" gorm:"column:owner_id"`
	PlayTurn     int       `json:"playTurn" gorm:"column:play_turn"`
	Timer        int       `json:"timeOut" gorm:"column:timer"`
	StartTime    time.Time `json:"startTime" gorm:"column:start_time"`
	GameID       int       `json:"gameID" gorm:"column:game_id"` // 1: 개굴작 2: 윙스팬
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
	State  string `json:"state" gorm:"column:state"`    // owned, discard, none, opened
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
	Project  string `json:"project" gorm:"column:project"`
	IsActive bool   `json:"isActive" gorm:"column:is_active"`
}

type Reports struct {
	gorm.Model
	TargetUserID   int    `json:"targetUserID" gorm:"column:target_user_id"`
	ReporterUserID int    `json:"reporterUserID" gorm:"column:reporter_user_id"`
	CategoryID     int    `json:"categoryID" gorm:"column:category_id"`
	Reason         string `json:"reason" gorm:"column:reason"`
}

type Categories struct {
	gorm.Model
	Reason string `json:"reason" gorm:"column:reason"`
	Type   string `json:"type" gorm:"column:type"`
}

type Profiles struct {
	gorm.Model
	Name        string `json:"name" gorm:"column:name"`
	TotalCount  int    `json:"totalCount" gorm:"column:total_count"`
	Image       string `json:"image" gorm:"column:image"`
	Description string `json:"description" gorm:"column:description"`
}

type UserProfiles struct {
	gorm.Model
	UserID       int  `json:"userID" gorm:"column:user_id"`
	ProfileID    int  `json:"profileID" gorm:"column:profile_id"`
	IsAchieved   bool `json:"isAchieved" gorm:"column:is_achieved"`
	CurrentCount int  `json:"currentCount" gorm:"column:current_count"`
}

type Missions struct {
	gorm.Model
	Name        string `json:"name" gorm:"column:name"`
	Description string `json:"description" gorm:"column:description"`
	Image       string `json:"image" gorm:"column:image"`
}

type RoomMissions struct {
	gorm.Model
	RoomID    int `json:"roomID" gorm:"column:room_id"`
	MissionID int `json:"missionID" gorm:"column:mission_id"`
}

type BirdCards struct {
	gorm.Model
	Name          string `json:"name" gorm:"column:name"`
	Image         string `json:"image" gorm:"column:image"`
	Description   string `json:"description" gorm:"column:description"`
	Size          int    `json:"size" gorm:"column:size"`
	Habitat       string `json:"habitat" gorm:"column:habitat"`              // water, forest, field
	BeakDirection string `json:"beakDirection" gorm:"column:beak_direction"` //left, right, center
	Nest          string `json:"nest" gorm:"column:nest"`                    //그릇형 bowl, 구멍둥지 cavity, 자유형 wild, 땅둥지 ground, 평평형 platform
}

type UserMissions struct {
	gorm.Model
	UserID    int `json:"userID" gorm:"column:user_id"`
	MissionID int `json:"missionID" gorm:"column:mission_id"`
	RoomID    int `json:"roomID" gorm:"column:room_id"`
}

type UserMissionCards struct {
	gorm.Model
	UserMissionID int `json:"userMissionID" gorm:"column:user_mission_id"`
	CardID        int `json:"cardID" gorm:"column:card_id"`
}

type UserBirdCards struct {
	gorm.Model
	UserID int    `json:"userID" gorm:"column:user_id"`
	CardID int    `json:"cardID" gorm:"column:card_id"`
	RoomID int    `json:"roomID" gorm:"column:room_id"`
	State  string `json:"state" gorm:"column:state"`
}

type UserTokens struct {
	gorm.Model
	UserID uint   `json:"userID" gorm:"column:user_id"`
	Token  string `json:"token" gorm:"column:token"`
}

type Items struct {
	gorm.Model
	Name        string `json:"name" gorm:"column:name"`
	Description string `json:"description" gorm:"column:description"`
	MaxUses     int    `json:"maxUses" gorm:"column:max_uses"`
}

type UserItems struct {
	gorm.Model
	UserID        int `json:"userID" gorm:"column:user_id"`
	ItemID        int `json:"itemID" gorm:"column:item_id"`
	RoomID        int `json:"roomID" gorm:"column:room_id"`
	RemainingUses int `json:"remainingUses" gorm:"column:remaining_uses"`
}

type FrogRoomUsers struct {
	gorm.Model
	UserID         int    `json:"userID" gorm:"column:user_id"`
	RoomID         int    `json:"roomID" gorm:"column:room_id"`
	Score          int    `json:"score" gorm:"column:score"`
	OwnedCardCount int    `json:"ownedCardCount" gorm:"column:owned_card_count"`
	PlayerState    string `json:"playerState" gorm:"column:player_state"`
	TurnNumber     int    `json:"turnNumber" gorm:"column:turn_number"`
}

type FrogCards struct {
	gorm.Model
	Name  string `json:"name" gorm:"column:name"`
	Color string `json:"color" gorm:"column:color"`
	Image string `json:"image" gorm:"column:image"`
}

type FrogUserCards struct {
	gorm.Model
	UserID int    `json:"userID" gorm:"column:user_id"`
	CardID int    `json:"cardID" gorm:"column:card_id"`
	RoomID int    `json:"roomID" gorm:"column:room_id"`
	State  string `json:"state" gorm:"column:state"`
}

// 틀린 그림 찾기 게임 테이블
type GameUsers struct {
	gorm.Model
	Name         string `json:"name" gorm:"column:name"`
	Email        string `json:"email" gorm:"column:email"`
	Password     string `json:"password" gorm:"column:password"`
	ProfileID    int    `json:"profileID" gorm:"column:profile_id"`
	Coin         int    `json:"coin" gorm:"column:coin"`
	State        string `json:"state" gorm:"column:state"`
	RoomID       int    `json:"roomID" gorm:"column:room_id"`
	Provider     string `json:"provider" gorm:"column:provider"`
	AlertEnabled bool   `json:"alertEnabled" gorm:"column:alert_enabled"`
}

type GameRooms struct {
	gorm.Model
	CurrentCount int       `json:"currentCount" gorm:"column:current_count"`
	MaxCount     int       `json:"maxCount" gorm:"column:max_count"`
	MinCount     int       `json:"minCount" gorm:"column:min_count"`
	Name         string    `json:"name" gorm:"column:name"`
	Password     string    `json:"password" gorm:"column:password"`
	State        string    `json:"state" gorm:"column:state"`
	OwnerID      int       `json:"ownerID" gorm:"column:owner_id"`
	GameID       int       `json:"gameID" gorm:"column:game_id"`
	StartTime    time.Time `json:"startTime" gorm:"column:start_time"`
}
type GameRoomUsers struct {
	gorm.Model
	UserID      int    `json:"userID" gorm:"column:user_id"`
	RoomID      int    `json:"roomID" gorm:"column:room_id"`
	PlayerState string `json:"playerState" gorm:"column:player_state"`
}
type FindItImages struct {
	gorm.Model
	Level            int    `json:"level" gorm:"column:level"`
	NormalImageUrl   string `json:"normalImageUrl" gorm:"column:normal_image_url"`     // ✅ 일반 이미지 URL
	AbnormalImageUrl string `json:"abnormalImageUrl" gorm:"column:abnormal_image_url"` // ✅ 비정상 이미지 URL
}
type FindItRoomSettings struct {
	gorm.Model
	RoomID             int `json:"roomID" gorm:"column:room_id"`
	Timer              int `json:"timer" gorm:"column:timer"`
	Lifes              int `json:"lifes" gorm:"column:lifes"`
	ItemHintCount      int `json:"itemHintCount" gorm:"column:item_hint_count"`
	Round              int `json:"round" gorm:"column:round"`
	ItemTimerStopCount int `json:"itemTimerStopCount" gorm:"column:item_timer_stop_count"`
}

type FindItCorrectPositions struct {
	gorm.Model
	RoomID            int `json:"roomID" gorm:"column:room_id"`
	UserID            int `json:"userID" gorm:"column:user_id"`
	Round             int `json:"round" gorm:"column:round"`
	ImageID           int `json:"imageID" gorm:"column:image_id"`                      // ✅ 정답을 맞춘 이미지 ID
	CorrectPositionID int `json:"correctPositionID" gorm:"column:correct_position_id"` // ✅ 맞춘 정답의 ID
}

type FindItRoundImages struct {
	gorm.Model
	RoomID     int `json:"roomID" gorm:"column:room_id"`
	Round      int `json:"round" gorm:"column:round"`
	ImageSetId int `json:"imageSetId" gorm:"column:image_set_id"`
}

type FindItImageCorrectPositions struct {
	gorm.Model
	ImageID   int     `json:"imageID" gorm:"column:image_id"`
	XPosition float64 `json:"xPosition" gorm:"column:x_position"`
	YPosition float64 `json:"yPosition" gorm:"column:y_position"`
}

type FindItUserCorrectPositions struct {
	gorm.Model
	UserID            int `json:"userID" gorm:"column:user_id"`
	RoomID            int `json:"roomID" gorm:"column:room_id"`
	Round             int `json:"round" gorm:"column:round"`
	ImageID           int `json:"imageID" gorm:"column:image_id"`                      // ✅ 정답을 맞춘 이미지 ID
	CorrectPositionID int `json:"correctPositionID" gorm:"column:correct_position_id"` // ✅ 맞춘 정답의 ID
}
