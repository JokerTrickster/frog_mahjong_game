package request

type ReqWSJoinPlay struct {
	Tkn       string `query:"tkn" validate:"required"`
	Password  string `query:"password" validate:"required"`
	SessionID string `query:"sessionID"`
}

type ReqWSJoinPlayEvent struct {
	Password string `json:"password"`
}
