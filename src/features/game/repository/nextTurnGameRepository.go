package repository

import (
	"context"
	_interface "main/features/game/model/interface"
	"main/features/game/model/request"
	"main/utils"

	"gorm.io/gorm"
)

func NewNextTurnGameRepository(gormDB *gorm.DB) _interface.INextTurnGameRepository {
	return &NextTurnGameRepository{GormDB: gormDB}
}

func (d *NextTurnGameRepository) UpdatePlayerNextTurn(c context.Context, req *request.ReqNextTurn) error {
	// 해당 턴 넘버를 가진 room user가 play_wait인지 확인 후 플레이 상태를 play로 변경
	err := d.GormDB.Table("room_users").Where("room_id = ? AND user_id = ? and turn_number = ? and player_state = ?", req.RoomID, req.UserID, req.TurnNumber, "play_wait").Update("player_state", "play").Error
	if err != nil {
		return utils.ErrorMsg(c, utils.ErrInternalDB, utils.Trace(), err.Error(), utils.ErrFromMysqlDB)
	}
	return nil
}
