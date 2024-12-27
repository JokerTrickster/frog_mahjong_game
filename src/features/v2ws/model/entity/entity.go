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
	WSClients    = make(map[string]*WSClient) // sessionID -> WSClient
	WSBroadcast  = make(chan WSMessage)       // 브로드캐스트 메시지
	RoomSessions = make(map[uint][]string)    // roomID -> sessionID 리스트
)

type WSClient struct {
	SessionID string // 고유 세션 ID
	RoomID    uint
	UserID    uint
	Conn      *websocket.Conn
	Closed    bool // 연결이 닫혔는지 여부를 추적하는 필드
}

type WSMessage struct {
	Message   string `json:"message"`
	Event     string `json:"event"`
	RoomID    uint   `json:"roomID"`
	UserID    uint   `json:"userID"`
	ChatID    uint   `json:"chatID"`
	SessionID string `json:"sessionID"`
	Name      string `json:"name"`
}

type ChatInfo struct {
	Name      string     `json:"name"`
	UserID    uint       `json:"userID"`
	Message   string     `json:"message"`
	ErrorInfo *ErrorInfo `json:"errorInfo"` // 에러 정보
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
	Users     []*User    `json:"users"`     // 유저 정보
	GameInfo  *GameInfo  `json:"gameInfo"`  // 게임 정보
	ErrorInfo *ErrorInfo `json:"errorInfo"` // 에러 정보

}

type ErrorInfo struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Type string `json:"type"`
}
type User struct {
	ID                  uint    `json:"id"`
	Email               string  `json:"email"`
	Name                string  `json:"name"`
	PlayerState         string  `json:"playerState"`         // 카드를 모두 선택했다면 done, 아니면 선택중이라면 picking
	IsOwner             bool    `json:"isOwner"`             // 방장 여부
	Cards               []*Card `json:"cards"`               // 소유하고 있는 카드
	DiscardedCards      []*Card `json:"discardedCards"`      //버린 카드
	PickedCards         []*Card `json:"pickedCards"`         // 선택한 카드
	Items               []*Item `json:"items"`               // 아이템 남은 횟수 (아이템ID, 남은 횟수)
	Coin                int     `json:"coin"`                // 보유하고 있는 코인
	ProfileID           int     `json:"profileID"`           // 프로필 ID
	MissionSuccessCount int     `json:"missionSuccessCount"` // 미션 성공 횟수
}
type GameInfo struct {
	PlayTurn   int    `json:"playTurn"`
	MissionIDs []int  `json:"missionIDs"` // 미션 ID 리스트
	AllReady   bool   `json:"allReady"`   // 게임 시작 여부
	Timer      int    `json:"timer"`      // 타이머
	IsFull     bool   `json:"isFull"`     // 방이 꽉 찼는지 여부
	AllPicked  bool   `json:"allPicked"`  // 모든 유저가 카드를 선택했는지 여부
	RoomID     uint   `json:"roomID"`     // 방 ID
	Password   string `json:"password"`   // 방 비밀번호
	Winner     uint   `json:"winner"`     // 승리자 ID
	OpenCards  []int  `json:"openCards"`  // 공개된 카드
	StartTime  int64  `json:"startTime"`  // 게임 시작 시간 (epoch time in milliseconds)
}

/*
이름 : oen, two, three, four, five, six, seven, eight, nine , chung, bal
색깔 : green, red, normal
상태 : 버려진 패 or 소유하고 있는 패 or 가운데 놓여져 있는 패
유저 ID
*/
type Card struct {
	CardID uint `json:"cardID"`
	UserID uint `json:"userID"`
}

type Item struct {
	ItemID        uint `json:"itemID"`
	RemainingUses int  `json:"remainingUses"`
}
type RoomUsers struct {
	gorm.Model
	UserID         int                   `json:"userID" gorm:"column:user_id"`
	RoomID         int                   `json:"roomID" gorm:"column:room_id"`
	Score          int                   `json:"score" gorm:"column:score"`
	OwnedCardCount int                   `json:"ownedCardCount" gorm:"column:owned_card_count"`
	PlayerState    string                `json:"playerState" gorm:"column:player_state"`
	TurnNumber     int                   `json:"turnNumber" gorm:"column:turn_number"`
	User           mysql.Users           `gorm:"foreignKey:UserID"`
	Room           mysql.Rooms           `gorm:"foreignKey:RoomID"`
	RoomMission    []mysql.RoomMissions  `gorm:"foreignKey:RoomID;references:RoomID"`
	Cards          []mysql.UserBirdCards `gorm:"foreignKey:UserID,RoomID;references:UserID,RoomID"`
	UserMissions   []mysql.UserMissions  `gorm:"foreignKey:UserID,RoomID;references:UserID,RoomID"`
	UserItems      []mysql.UserItems     `gorm:"foreignKey:UserID,RoomID;references:UserID,RoomID"`
	RoomUsers      mysql.RoomUsers       `gorm:"foreignKey:UserID,RoomID;references:UserID,RoomID"`
}

func (c *WSClient) Close() {
	c.Closed = true
	c.Conn.Close()
}

func (c *WSClient) IsClosed() bool {
	return c.Closed
}
