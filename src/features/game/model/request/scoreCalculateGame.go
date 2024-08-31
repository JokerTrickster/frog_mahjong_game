package request

type ReqScoreCalculate struct {
	RoomID uint        `json:"roomID"`
	Cards  []ScoreCard `json:"cards"`
}

type ScoreCard struct {
	CardID uint `json:"cardID"`
}
