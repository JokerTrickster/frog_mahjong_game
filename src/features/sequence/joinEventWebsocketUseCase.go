package sequence

import "main/utils/db/mysql"

func CreateJoinUserItemDTO(uID uint, roomID uint, item mysql.Items) mysql.UserItems {
	result := mysql.UserItems{
		UserID:        int(uID),
		RoomID:        int(roomID),
		ItemID:        int(item.ID),
		RemainingUses: item.MaxUses,
	}
	return result
}

func CreateJoinPlayUserDTO(uID uint, roomID uint) *mysql.SequenceUsers {
	result := &mysql.SequenceUsers{
		UserID:    int(uID),
		RoomID:    int(roomID),
		Turn:      0,
		ColorType: 0,
	}
	return result
}

func CreateJoinPlayGameRoomSettingDTO(roomID uint) *mysql.SequenceGameRoomSettings {
	result := &mysql.SequenceGameRoomSettings{
		RoomID:       int(roomID),
		Timer:        60,
		CurrentRound: 1,
	}
	return result
}
