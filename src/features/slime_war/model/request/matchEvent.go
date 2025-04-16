package request

// 타이머, 방 인원 수
type ReqWSMatch struct {
	Tkn       string `query:"tkn" validate:"required"`
	SessionID string `query:"sessionID"`
}
