package entity

import (
	"main/utils/db/mysql"
	"net/http"

	"github.com/gorilla/websocket"
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

type MessageInfo struct {
	Users            []*User           `json:"users"`            // 유저 정보
	SlimeWarGameInfo *SlimeWarGameInfo `json:"slimeWarGameInfo"` // 게임 정보
	ErrorInfo        *ErrorInfo        `json:"errorInfo"`        // 에러 정보

}

type ErrorInfo struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Type string `json:"type"`
}
type User struct {
	ID             uint        `json:"id"`
	Email          string      `json:"email"`
	Name           string      `json:"name"`
	IsOwner        bool        `json:"isOwner"`        // 방장 여부
	ProfileID      int         `json:"profileID"`      // 프로필 ID
	Turn           int         `json:"turn"`           // 턴
	SlimePositions []*Position `json:"slimePositions"` // 슬라임 위치
	HeroCardCount  int         `json:"heroCardCount"`  // 영웅 카드 개수
}
type SlimeWarGameInfo struct {
	AllReady              bool      `json:"allReady"`              // 게임 시작 여부
	Timer                 int       `json:"timer"`                 // 타이머
	IsFull                bool      `json:"isFull"`                // 방이 꽉 찼는지 여부
	RoomID                uint      `json:"roomID"`                // 방 ID
	Password              string    `json:"password"`              // 방 비밀번호
	StartTime             int64     `json:"startTime"`             // 게임 시작 시간 (epoch time in milliseconds)
	Round                 int       `json:"round"`                 // 라운드
	RoundCount            int       `json:"roundCount"`            // 라운드 카운트
	KingPosition          *Position `json:"kingPosition"`          // 왕의 위치
	SlimeCount            int       `json:"slimeCount"`            // 슬라임 개수
	DroppedDummyIndices   []int     `json:"droppedDummyIndices"`   // 버려진 더미 인덱스 배열
	RemainingDummyIndices []int     `json:"remainingDummyIndices"` // 남은 더미 인덱스 배열
}
type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// PreloadUsers - 게임 방에 있는 유저 정보 + 관련 데이터 로드
type PreloadUsers struct {
	UserID uint             `json:"userID" gorm:"column:user_id"`  // 유저 ID
	RoomID uint             `json:"roomID" gorm:"column:room_id"`  // 방 ID
	User   *mysql.GameUsers `json:"user" gorm:"foreignKey:UserID"` // 유저 정보 (game_users)
	Room   *mysql.GameRooms `json:"room" gorm:"foreignKey:RoomID"` // 방 정보 (game_rooms)
}

func (c *WSClient) Close() {
	c.Closed = true
	c.Conn.Close()
}

func (c *WSClient) IsClosed() bool {
	return c.Closed
}
