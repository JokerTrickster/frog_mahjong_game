package entity

type ScoreCalculateEntitySQL struct {
	UserID uint
	RoomID uint
	Cards  []uint
}

type ScoreCalculateEntity struct {
	Cards []ScoreCalculateCard `json:"cards"`
}

type ScoreCalculateCard struct {
	CardID uint   `json:"cardID"`
	Name   string `json:"name"`
	Color  string `json:"color"`
}
