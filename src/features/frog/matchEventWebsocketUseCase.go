package ws

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
		GameID:       mysql.FROG,
		StartTime:    time.Now(),
	}
	return result
}

func CreateMatchRoomUserDTO(uID uint, roomID int) *mysql.FrogRoomUsers {
	result := &mysql.FrogRoomUsers{
		UserID:         int(uID),
		RoomID:         roomID,
		Score:          0,
		OwnedCardCount: 0,
		PlayerState:    "play",
		TurnNumber:     0,
	}
	return result
}
