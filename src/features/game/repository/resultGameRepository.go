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
func (d *ResultGameRepository) FindCards(c context.Context, entitySQL *entity.ResultEntitySQL) ([]*mysql.FrogUserCards, error) {
	var frogUserCards []*mysql.FrogUserCards
	err := d.GormDB.Model(&frogUserCards).Where("room_id = ? and card_id IN ?", entitySQL.RoomID, entitySQL.Cards).Find(&frogUserCards).Error
	if err != nil {
		return nil, utils.ErrorMsg(c, utils.ErrInternalDB, utils.Trace(), utils.HandleError(err.Error(),entitySQL), utils.ErrFromMysqlDB)
	}
	return frogUserCards, nil
}

func (d *ResultGameRepository) GetDoraCard(c context.Context, req *request.ReqResult) (*mysql.FrogUserCards, error) {
	// dora 카드를 가져온다.
	var doraCard *mysql.FrogUserCards
	err := d.GormDB.Model(&doraCard).Where("room_id = ? AND state = ?", req.RoomID, "dora").First(&doraCard).Error
	if err != nil {
		return &mysql.FrogUserCards{}, utils.ErrorMsg(c, utils.ErrBadParameter, utils.Trace(),utils.HandleError( err.Error(),req), utils.ErrFromClient)
	}
	return doraCard, nil
}

func (d *ResultGameRepository) FindOneFrogCard(c context.Context, cardID uint) (*mysql.FrogCards, error) {
	// 카드 ID 에 해당되는 개구리 카드 정보를 가져온다.
	var frogCard mysql.FrogCards
	err := d.GormDB.Model(&frogCard).Where("id = ?", cardID).First(&frogCard).Error
	if err != nil {
		return nil, utils.ErrorMsg(c, utils.ErrInternalDB, utils.Trace(), utils.HandleError(err.Error(),cardID), utils.ErrFromMysqlDB)
	}
	return &frogCard, nil
}