package request

// 타이머, 방 인원 수
type ReqWSMatch struct {
	Tkn       string `query:"tkn" validate:"required"`
	Timer     int    `query:"timer" validate:"required"`
	Count     int    `query:"count" validate:"required"`
	SessionID string `query:"sessionID"`
}

type ReqWSMatchEvent struct {
	Timer int `json:"timer" validate:"required"`
	Count int `json:"count" validate:"required"`
}
