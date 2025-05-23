package request

type ReqFrogCardListBoardGame struct {
	RoomID int `query:"room_id" validate:"required"`
}
