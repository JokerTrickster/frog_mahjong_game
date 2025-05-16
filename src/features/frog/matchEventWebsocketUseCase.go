package ws

import (
	"main/utils/db/mysql"
	"time"
)

func CreateMatchRoomDTO(uID uint, count int, timer int) *mysql.Rooms {
	result := &mysql.Rooms{
		CurrentCount: 0,
		MaxCount:     count,
		MinCount:     count,
		State:        "wait",
		OwnerID:      int(uID),
		Timer:        timer,
		PlayTurn:     0,
		Name:         "speed match",
		GameID:       1,
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
