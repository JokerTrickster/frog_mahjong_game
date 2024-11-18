package request

type ReqV2WSMissionEvent struct {
	MissionIDs []int `json:"missionIDs"`
	Cards      []int `json:"cards"`
}
