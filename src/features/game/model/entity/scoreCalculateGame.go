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
	Count  int    `json:"count"`
	Color  string `json:"color"`
}
