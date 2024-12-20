package repository

import (
	"context"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

type JoinEventWebsocketRepository struct {
	GormDB *gorm.DB
}

func FindAllOpenCards(c context.Context, roomID int) ([]int, error) {
	var cards []int
	if err := mysql.GormMysqlDB.WithContext(c).Model(&mysql.UserBirdCards{}).Where("room_id = ? and state = ?", roomID, "opened").Pluck("card_id", &cards).Error; err != nil {
		return nil, err
	}
	return cards, nil
}

func ReconnectedUpdateRoomUser(c context.Context, roomID uint, userID uint) error {
	err := mysql.GormMysqlDB.Model(&mysql.RoomUsers{}).Where("room_id = ? and user_id = ?", roomID, userID).Update("player_state", "play").Error
	if err != nil {
		return err
	}
	return nil
}
