package request

type ReqStart struct {
	RoomID uint   `json:"roomID"`
	State  string `json:"state"`
}
