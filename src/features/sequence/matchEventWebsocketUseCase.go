package sequence

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
		GameID:       mysql.SEQUENCE,
	}
	return result
}

func CreateMatchUserDTO(uID uint, roomID uint) *mysql.SequenceUsers {
	result := &mysql.SequenceUsers{
		UserID:    int(uID),
		RoomID:    int(roomID),
		Turn:      0,
		ColorType: 0,
	}
	return result
}

func CreateMatchGameRoomSettingDTO(roomID uint) *mysql.SequenceGameRoomSettings {
	result := &mysql.SequenceGameRoomSettings{
		RoomID:       int(roomID),
		Timer:        60,
		CurrentRound: 1,
	}
	return result
}
