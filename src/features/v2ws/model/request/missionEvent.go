package request

type ReqV2WSMissionEvent struct {
	MissionID int   `json:"missionID"`
	Cards     []int `json:"cards"`
}
