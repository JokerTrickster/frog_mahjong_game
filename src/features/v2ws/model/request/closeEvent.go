package request

type ReqWSClose struct {
	RoomID uint `json:"roomID" validate:"required"`
}
