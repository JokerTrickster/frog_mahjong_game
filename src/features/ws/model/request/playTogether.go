package request

// 타이머, 방 인원 수
type ReqWSPlayTogether struct {
	Tkn string `query:"tkn" validate:"required"`
}
type ReqWSPlayTogetherEvent struct {
	Timer int `json:"timer" validate:"required"`
	Count int `json:"count" validate:"required"`
}
