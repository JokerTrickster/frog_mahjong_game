package entity

type WSDoraEntity struct {
	RoomID uint   `json:"roomID omitempty"`
	Name   string `json:"name"`
	Color  string `json:"color"`
	State  string `json:"state"`
}
