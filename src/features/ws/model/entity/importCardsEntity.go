package entity

import "main/utils/db/mysql"

type WSImportCardsEntity struct {
	RoomID uint           `json:"roomID omitempty"`
	UserID uint           `json:"userID"`
	Cards  []*mysql.Cards `json:"cards"`
}
