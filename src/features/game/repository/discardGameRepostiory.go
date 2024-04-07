package repository

import (
	"context"
	_errors "main/features/game/model/errors"
	_interface "main/features/game/model/interface"
	"main/features/game/model/request"
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewDiscardGameRepository(gormDB *gorm.DB) _interface.IDiscardGameRepository {
	return &DiscardGameRepository{GormDB: gormDB}
}

func (d *DiscardGameRepository) PlayerCheckTurn(c context.Context, req *request.ReqDiscard) (mysql.RoomUsers, error) {
	// room_id, user_id로 찾고 player_state가 play_turn인지 체크
	var roomUser mysql.RoomUsers
	err := d.GormDB.Model(&roomUser).Where("room_id = ? AND user_id = ? AND player_state = ?", req.RoomID, req.UserID, "play").First(&roomUser).Error
	if err != nil {
		return mysql.RoomUsers{}, utils.ErrorMsg(c, utils.ErrBadRequest, utils.Trace(), _errors.ErrBadRequest.Error(), utils.ErrFromMysqlDB)
	}
	return roomUser, nil
}
func (d *DiscardGameRepository) UpdateCardStateDiscard(c context.Context, req *request.ReqDiscard) error {
	// 카드 상태 업데이트
	// room_id, card_id, state로 찾고 카드 업데이트할 때 트랜잭션 처리해줘
	err := d.GormDB.Model(&mysql.Cards{}).Where("room_id = ? AND id = ? AND state = ? and user_id = ?", req.RoomID, req.CardID, "owned", req.UserID).Update("state", "discard").Error
	if err != nil {
		return utils.ErrorMsg(c, utils.ErrInternalDB, utils.Trace(), err.Error(), utils.ErrFromMysqlDB)
	}

	return nil
}

func (d *DiscardGameRepository) UpdateRoomUser(c context.Context, updateRoomUser mysql.RoomUsers) error {
	// user_id, room_id로 찾고  owned_card_count -1 & player_state를 play_wait로 업데이트 해줘
	err := d.GormDB.Model(&updateRoomUser).Where("room_id = ? AND user_id = ?", updateRoomUser.RoomID, updateRoomUser.UserID).Updates(updateRoomUser).Error
	if err != nil {
		return utils.ErrorMsg(c, utils.ErrInternalDB, utils.Trace(), err.Error(), utils.ErrFromMysqlDB)
	}
	return nil
}
