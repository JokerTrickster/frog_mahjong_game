package repository

import (
	"context"
	_errors "main/features/game/model/errors"
	_interface "main/features/game/model/interface"
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewV2DeckCardGameRepository(gormDB *gorm.DB) _interface.IV2DeckCardGameRepository {
	return &V2DeckCardGameRepository{GormDB: gormDB}
}

func (d *V2DeckCardGameRepository) CheckRoomUser(c context.Context, userID int, roomID int) error {
	// room_id, user_id로 찾고 player_state가 play_turn인지 체크
	var roomUser mysql.RoomUsers
	err := d.GormDB.Model(&roomUser).Where("room_id = ? AND user_id = ?", roomID, userID).First(&roomUser).Error
	if err != nil {
		return utils.ErrorMsg(c, utils.ErrBadRequest, utils.Trace(), utils.HandleError(_errors.ErrBadRequest.Error(),userID,roomID), utils.ErrFromMysqlDB)
	}
	return nil

}
