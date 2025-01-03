package ws

import (
	"main/utils/db/mysql"
)

func CreateRoomUserDTO(uID uint, roomID int) (mysql.FrogRoomUsers, error) {
	result := mysql.FrogRoomUsers{
		UserID:         int(uID),
		RoomID:         roomID,
		Score:          0,
		OwnedCardCount: 0,
		PlayerState:    "play",
	}
	return result, nil
}
