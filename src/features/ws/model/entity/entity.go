package entity

import (
	"main/utils/db/mysql"
	"net/http"

	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

var (
	WSUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		}}
	WSClients   = make(map[uint]map[*websocket.Conn]WSClient)
	WSBroadcast = make(chan WSMessage)
)

type WSClient struct {
	RoomID uint
	UserID uint
	Conn   *websocket.Conn
}

type WSMessage struct {
	Message string `json:"message"`
	Event   string `json:"event"`
	RoomID  uint   `json:"roomID"`
	UserID  uint   `json:"userID"`
}

/*
유저 ID

	이메일
	이름
	유저 상태 : ready or not ready
	방장인지 여부 :
	가지고 있는 패 정보들 :
	버린 패 정보들 :
	현재 보유하고 있는 코인 :
*/
type RoomInfo struct {
	Users    []*User   `json:"users"`    // 유저 정보
	GameInfo *GameInfo `json:"gameInfo"` // 게임 정보
}
type User struct {
	ID             uint    `json:"id"`
	Email          string  `json:"email"`
	Name           string  `json:"name"`
	PlayerState    string  `json:"playerState"`
	IsOwner        bool    `json:"isOwner"`
	Cards          []*Card `json:"cards"`
	DiscardedCards []*Card `json:"discardedCards"`
	Coin           int     `json:"coin"`
	TurnNumber     int     `json:"turnNumber"`
}
type GameInfo struct {
	PlayTurn int   `json:"playTurn"`
	Dora     *Card `json:"dora"`     // 도라
	AllReady bool  `json:"allReady"` // 게임 시작 여부
}

/*
카드 ID
이름 : oen, two, three, four, five, six, seven, eight, nine , chung, bal
색깔 : green, red, normal
상태 : 버려진 패 or 소유하고 있는 패 or 가운데 놓여져 있는 패
유저 ID
*/
type Card struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Color  string `json:"color"`
	State  string `json:"state"`
	UserID uint   `json:"userID"`
}

type RoomUsers struct {
	gorm.Model
	UserID         int         `json:"userID" gorm:"column:user_id"`
	RoomID         int         `json:"roomID" gorm:"column:room_id"`
	Score          int         `json:"score" gorm:"column:score"`
	OwnedCardCount int         `json:"ownedCardCount" gorm:"column:owned_card_count"`
	PlayerState    string      `json:"playerState" gorm:"column:player_state"`
	TurnNumber     int         `json:"turnNumber" gorm:"column:turn_number"`
	User           mysql.Users `gorm:"foreignKey:UserID"`
	Room           mysql.Rooms `gorm:"foreignKey:RoomID"`
}