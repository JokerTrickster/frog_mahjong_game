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

func CreateMatchUserDTO(uID uint, roomID uint) *mysql.SlimeWarUsers {
	result := &mysql.SlimeWarUsers{
		UserID:    int(uID),
		RoomID:    int(roomID),
		HeroCount: 4,
		Turn:      0,
		ColorType: 0,
	}
	return result
}

func CreateMatchGameRoomSettingDTO(roomID uint) *mysql.SlimeWarGameRoomSettings {
	result := &mysql.SlimeWarGameRoomSettings{
		RoomID:              int(roomID),
		Timer:               60,
		RemainingCardCount:  48,
		KingIndex:           50,
		CurrentRound:        1,
		RemainingSlimeCount: 52,
	}
	return result
}
