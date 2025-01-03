package v2ws

import (
	"main/utils/db/mysql"
	"time"
)

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
		StartTime:    time.Now(),
		GameID: 2,
	}
	return result
}

func CreateMatchRoomUserDTO(uID uint, roomID int) mysql.RoomUsers {
	result := mysql.RoomUsers{
		UserID:         int(uID),
		RoomID:         roomID,
		Score:          0,
		OwnedCardCount: 0,
		PlayerState:    "play",
		TurnNumber:     0,
	}
	return result
}

func CreateMatchUserItemDTO(uID uint, roomID uint, item mysql.Items) mysql.UserItems {
	result := mysql.UserItems{
		UserID:        int(uID),
		RoomID:        int(roomID),
		ItemID:        int(item.ID),
		RemainingUses: item.MaxUses,
	}
	return result
}
