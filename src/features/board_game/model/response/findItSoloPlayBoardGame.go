package response

type ResFindItSoloPlayBoardGame struct {
	GameInfoList []SoloPlayGameInfo `json:"gameInfoList"`
}

type SoloPlayGameInfo struct {
	Round            int        `json:"round"`
	NormalUrl        string     `json:"normalUrl"`
	AbnormalUrl      string     `json:"abnormalUrl"`
	CorrectPositions []Position `json:"correctPositions"`
	ImageID          int        `json:"imageID"`
}

type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
