package entity

type V2WSMissionEntity struct {
	RoomID    uint  `json:"roomID omitempty"`
	UserID    uint  `json:"userID"`
	Cards     []int `json:"cards"`
	MissionID int   `json:"missionID"`
}
