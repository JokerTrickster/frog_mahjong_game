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
	ID               uint       `json:"id"`
	Email            string     `json:"email"`
	Name             string     `json:"name"`
	IsOwner          bool       `json:"isOwner"`          // 방장 여부
	ProfileID        int        `json:"profileID"`        // 프로필 ID
	CorrectPositions []Position `json:"correctPositions"` // 맞은 위치 수 (x,y)
}
type GameInfo struct {
	AllReady       bool       `json:"allReady"`       // 게임 시작 여부
	Timer          int        `json:"timer"`          // 타이머
	IsFull         bool       `json:"isFull"`         // 방이 꽉 찼는지 여부
	RoomID         uint       `json:"roomID"`         // 방 ID
	Password       string     `json:"password"`       // 방 비밀번호
	StartTime      int64      `json:"startTime"`      // 게임 시작 시간 (epoch time in milliseconds)
	ItemTimerCount int        `json:"itemTimerCount"` // 아이템 타이머 카운트
	ItemHintCount  int        `json:"itemHintCount"`  // 아이템 힌트 카운트
	Round          int        `json:"round"`          // 라운드
	ImageInfo      *ImageInfo `json:"imageInfo"`      // 이미지 정보
	Life           int        `json:"life"`           // 생명
	WrongPosition  *Position   `json:"wrongPosition"`  // 틀린 위치 (x,y)
	CorrectCount   int        `json:"correctCount"`   // 맞은 개수
	HintPosition   *Position   `json:"hintPosition"`   // 힌트 위치 (x,y)
	TimerUsed      bool       `json:"timerUsed"`      // 타이머 사용 여부
}
type ImageInfo struct {
	ID               int    `json:"id"`
	NormalImageUrl   string `json:"normalImageUrl"`   // 일반 이미지 URL
	AbnormalImageUrl string `json:"abnormalImageUrl"` // 비정상 이미지 URL
}
type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// PreloadUsers - 게임 방에 있는 유저 정보 + 관련 데이터 로드
type PreloadUsers struct {
	UserID               uint                                `json:"userID" gorm:"column:user_id"`                                    // 유저 ID
	RoomID               uint                                `json:"roomID" gorm:"column:room_id"`                                    // 방 ID
	User                 *mysql.GameUsers                    `json:"user" gorm:"foreignKey:UserID"`                                   // 유저 정보 (game_users)
	Room                 *mysql.GameRooms                    `json:"room" gorm:"foreignKey:RoomID"`                                   // 방 정보 (game_rooms)
	RoomSetting          *mysql.FindItRoomSettings           `json:"roomSetting" gorm:"foreignKey:RoomID;references:RoomID"`          // 방 설정 정보 (find_it_room_settings)
	UserCorrectPositions []*mysql.FindItUserCorrectPositions `json:"userCorrectPositions" gorm:"foreignKey:UserID;references:UserID"` // 유저가 맞춘 정답 정보 (find_it_user_correct_positions)
	RoundImages          []*mysql.FindItRoundImages          `json:"roundImages" gorm:"foreignKey:RoomID;references:RoomID"`          // 해당 방의 라운드별 이미지 정보 (find_it_round_images)
}

func (c *WSClient) Close() {
	c.Closed = true
	c.Conn.Close()
}

func (c *WSClient) IsClosed() bool {
	return c.Closed
}
