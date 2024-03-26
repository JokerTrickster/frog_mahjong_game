package request

type ReqReady struct {
	RoomID      uint   `json:"roomID"`
	PlayerState string `json:"playerState"`
}
