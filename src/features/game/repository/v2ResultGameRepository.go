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

func NewV2ResultGameRepository(gormDB *gorm.DB) _interface.IV2ResultGameRepository {
	return &V2ResultGameRepository{GormDB: gormDB}
}

// 카드 ID 에 해당되는 카드 정보들을 가져온다.
func (d *V2ResultGameRepository) FindCards(c context.Context, entitySQL *entity.ResultEntitySQL) ([]mysql.Cards, error) {
	var cards []mysql.Cards
	err := d.GormDB.Model(&cards).Where("room_id = ? and card_id IN ?", entitySQL.RoomID, entitySQL.Cards).Find(&cards).Error
	if err != nil {
		return nil, utils.ErrorMsg(c, utils.ErrInternalDB, utils.Trace(), utils.HandleError(err.Error(), entitySQL), utils.ErrFromMysqlDB)
	}
	return cards, nil
}

func (d *V2ResultGameRepository) GetDoraCard(c context.Context, req *request.ReqResult) (mysql.Cards, error) {
	// dora 카드를 가져온다.
	var doraCard mysql.Cards
	err := d.GormDB.Model(&doraCard).Where("room_id = ? AND state = ?", req.RoomID, "dora").First(&doraCard).Error
	if err != nil {
		return mysql.Cards{}, utils.ErrorMsg(c, utils.ErrBadParameter, utils.Trace(), utils.HandleError(err.Error(), req), utils.ErrFromClient)
	}
	return doraCard, nil
}
func (d *V2ResultGameRepository) GetUserMissions(c context.Context, req *request.ReqV2Result) ([]*mysql.UserMissions, error) {
	var userMissions []*mysql.UserMissions
	err := d.GormDB.Model(&userMissions).Where("room_id = ? and user_id = ?", req.RoomID, req.UserID).Find(&userMissions).Error
	if err != nil {
		return nil, utils.ErrorMsg(c, utils.ErrInternalDB, utils.Trace(), utils.HandleError(err.Error(), req), utils.ErrFromMysqlDB)
	}

	return userMissions, nil
}

func (d *V2ResultGameRepository) GetUserMissionCards(c context.Context, missionID uint) ([]*mysql.UserMissionCards, error) {
	var userMissionCards []*mysql.UserMissionCards
	err := d.GormDB.Model(&userMissionCards).Where("user_mission_id = ?", missionID).Find(&userMissionCards).Error
	if err != nil {
		return nil, utils.ErrorMsg(c, utils.ErrInternalDB, utils.Trace(), utils.HandleError(err.Error(), missionID), utils.ErrFromMysqlDB)
	}

	return userMissionCards, nil
}
