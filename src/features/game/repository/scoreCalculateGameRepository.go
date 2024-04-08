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

func NewScoreCalculateGameRepository(gormDB *gorm.DB) _interface.IScoreCalculateGameRepository {
	return &ScoreCalculateGameRepository{GormDB: gormDB}
}

func (d *ScoreCalculateGameRepository) CheckCardCount(c context.Context, userID uint, req *request.ReqScoreCalculate) error {
	// 카드가 6장 소유했는지 체크
	var roomUser mysql.RoomUsers
	err := d.GormDB.Model(&roomUser).Where("room_id = ? AND user_id = ?", req.RoomID, userID).First(&roomUser).Error
	if err != nil {
		return utils.ErrorMsg(c, utils.ErrBadParameter, utils.Trace(), _errors.ErrNotEnoughCard.Error(), utils.ErrFromClient)
	}
	return nil
}

func (d *ScoreCalculateGameRepository) GetDoraCard(c context.Context, req *request.ReqScoreCalculate) (mysql.Cards, error) {
	// dora 카드를 가져온다.
	var doraCard mysql.Cards
	err := d.GormDB.Model(&doraCard).Where("room_id = ? AND state = ?", req.RoomID, "dora").First(&doraCard).Error
	if err != nil {
		return mysql.Cards{}, utils.ErrorMsg(c, utils.ErrBadParameter, utils.Trace(), err.Error(), utils.ErrFromClient)
	}
	return doraCard, nil
}
