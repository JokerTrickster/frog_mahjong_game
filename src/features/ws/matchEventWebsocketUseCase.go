package ws

import "main/utils/db/mysql"

func CreateMatchRoomDTO(uID uint, count int, timer int) mysql.Rooms {
	result := mysql.Rooms{
		CurrentCount: 0,
		MaxCount:     count,
		MinCount:     count,
		State:        "wait",
		OwnerID:      int(uID),
		Timer:        timer,
		PlayTurn:     0,
		Name:         "speed match",
		GameID:       1,
	}
	return result
}

func CreateMatchRoomUserDTO(uID uint, roomID int, playerState string) mysql.RoomUsers {
	result := mysql.RoomUsers{
		UserID:         int(uID),
		RoomID:         roomID,
		Score:          0,
		OwnedCardCount: 0,
		PlayerState:    playerState,
		TurnNumber:     0,
	}
	return result
}
