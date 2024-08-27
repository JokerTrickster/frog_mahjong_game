package request

type ReqResult struct {
	RoomID uint         `json:"roomID"`
	Cards  []ResultCard `json:"cards"`
}

type ResultCard struct {
	CardID uint `json:"cardID"`
}
