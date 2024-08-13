package repository

import (
	"context"
	"main/features/game/model/entity"
	_interface "main/features/game/model/interface"
	"main/features/game/model/request"
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewScoreCalculateGameRepository(gormDB *gorm.DB) _interface.IScoreCalculateGameRepository {
	return &ScoreCalculateGameRepository{GormDB: gormDB}
}

// 카드 6장 소유했는지 체크하고 카드 정보를 가져온다.
func (d *ScoreCalculateGameRepository) FindOwnedCards(c context.Context, entitySQL *entity.ScoreCalculateEntitySQL) ([]mysql.Cards, error) {
	var cards []mysql.Cards
	err := d.GormDB.Model(&cards).Where("room_id = ? and user_id = ? and state = ?", entitySQL.RoomID, entitySQL.UserID, "owned").Find(&cards).Error
	if err != nil {
		return nil, utils.ErrorMsg(c, utils.ErrInternalDB, utils.Trace(), err.Error(), utils.ErrFromMysqlDB)
	}
	return cards, nil

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
