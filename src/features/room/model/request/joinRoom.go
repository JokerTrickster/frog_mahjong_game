package request

type ReqJoin struct {
	RoomID   uint   `json:"roomID" validate:"required"`
	Password string `json:"password,omitempty"`
}
