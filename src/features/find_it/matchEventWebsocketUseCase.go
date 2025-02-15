package find_it

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
		GameID:       1,
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
