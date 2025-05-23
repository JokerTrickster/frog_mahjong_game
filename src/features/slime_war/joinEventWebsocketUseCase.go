package slime_war

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

func CreateJoinPlayUserDTO(uID uint, roomID uint) *mysql.SlimeWarUsers {
	result := &mysql.SlimeWarUsers{
		UserID:    int(uID),
		RoomID:    int(roomID),
		HeroCount: 4,
		Turn:      0,
		ColorType: 0,
	}
	return result
}

func CreateJoinPlayGameRoomSettingDTO(roomID uint) *mysql.SlimeWarGameRoomSettings {
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
