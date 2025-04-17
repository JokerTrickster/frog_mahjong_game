package request

// 타이머, 방 인원 수
type ReqWSPlayTogether struct {
	Tkn string `query:"tkn" validate:"required"`
}
