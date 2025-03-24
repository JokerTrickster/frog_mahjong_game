package find_it

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
