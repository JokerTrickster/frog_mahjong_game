package response

type ResV2Result struct {
	Winner   uint            `json:"winner"`
	Missions []ResultMission `json:"missions"`
}

type ResultMission struct {
	MissionID uint  `json:"missionID"`
	Cards     []int `json:"cards"`
}
