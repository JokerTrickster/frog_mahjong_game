package request

// 타이머, 방 인원 수
type ReqWSNextRound struct {
	Round  int  `json:"round"`
	UserID uint `json:"userID"`
	OpponentCanMove bool `json:"opponentCanMove"`
}
