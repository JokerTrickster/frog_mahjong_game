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

func NewResultGameRepository(gormDB *gorm.DB) _interface.IResultGameRepository {
	return &ResultGameRepository{GormDB: gormDB}
}

// 카드 ID 에 해당되는 카드 정보들을 가져온다.
func (d *ResultGameRepository) FindCards(c context.Context, entitySQL *entity.ResultEntitySQL) ([]mysql.Cards, error) {
	var cards []mysql.Cards
	err := d.GormDB.Model(&cards).Where("room_id = ? and card_id IN ?", entitySQL.RoomID, entitySQL.Cards).Find(&cards).Error
	if err != nil {
		return nil, utils.ErrorMsg(c, utils.ErrInternalDB, utils.Trace(), err.Error(), utils.ErrFromMysqlDB)
	}
	return cards, nil
}

func (d *ResultGameRepository) GetDoraCard(c context.Context, req *request.ReqResult) (mysql.Cards, error) {
	// dora 카드를 가져온다.
	var doraCard mysql.Cards
	err := d.GormDB.Model(&doraCard).Where("room_id = ? AND state = ?", req.RoomID, "dora").First(&doraCard).Error
	if err != nil {
		return mysql.Cards{}, utils.ErrorMsg(c, utils.ErrBadParameter, utils.Trace(), err.Error(), utils.ErrFromClient)
	}
	return doraCard, nil
}
