package request

type ReqReady struct {
	RoomID      uint   `json:"roomID"`
	PlayerState string `json:"playerState" oneOf:"ready wait"`
}
