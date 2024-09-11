package request

type ReqWSJoin struct {
	Tkn      string `query:"tkn" validate:"required"`
	RoomID   uint   `query:"roomID" validate:"required"`
	Password string `query:"password"`
}
