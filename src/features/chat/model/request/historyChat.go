package request

type ReqHistory struct {
	RoomID  uint `query:"roomID" validate:"required"`
	Page     int `query:"page" validate:"omitempty,gte=0"`
	PageSize int `query:"pageSize" validate:"omitempty,gt=0"`
}
