package ws

import "main/utils/db/mysql"

func CreateRoomUserDTO(uID uint, roomID int, playerState string) (mysql.RoomUsers, error) {
	result := mysql.RoomUsers{
		UserID:         int(uID),
		RoomID:         roomID,
		Score:          0,
		OwnedCardCount: 0,
		PlayerState:    playerState,
	}
	return result, nil
}
