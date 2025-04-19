package slime_war

import (
	"main/utils/db/mysql"
	"time"
)

func CreateMatchRoomDTO(uID uint) *mysql.GameRooms {
	result := &mysql.GameRooms{
		CurrentCount: 0,
		MaxCount:     2,
		MinCount:     2,
		State:        "wait",
		OwnerID:      int(uID),
		Name:         "speed match",
		StartTime:    time.Now(),
		GameID:       mysql.SLIME_WAR,
	}
	return result
}
