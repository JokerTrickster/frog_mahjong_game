package repository

import (
	"context"
	_interface "main/features/game/model/interface"
	"main/features/game/model/request"
	"main/utils"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewOwnershipGameRepository(gormDB *gorm.DB) _interface.IOwnershipGameRepository {
	return &OwnershipGameRepository{GormDB: gormDB}
}

func (g *OwnershipGameRepository) UpdateCardState(c context.Context, req *request.ReqOwnership) error {
	// 카드 상태 업데이트
	// room_id, card_id, state로 찾고 카드 업데이트할 때 트랜잭션 처리해줘
	for _, card := range req.Cards {
		err := g.GormDB.Model(&mysql.Cards{}).Where("room_id = ? and id = ? and state = ?", card.RoomID, card.CardID, "none").Updates(&mysql.Cards{State: card.State, UserID: card.UserID}).Error
		if err != nil {
			return utils.ErrorMsg(c, utils.ErrInternalDB, utils.Trace(), err.Error(), utils.ErrFromMysqlDB)
		}
	}
	return nil
}

func (g *OwnershipGameRepository) UpdateRoomUserCardCount(c context.Context, req *request.ReqOwnership) error {
	// 유저id로 room_users에서 찾아서 card_count를 더한 후 업데이트 한다.
	for _, card := range req.Cards {
		err := g.GormDB.Model(&mysql.RoomUsers{}).Where("room_id = ? AND user_id = ?", card.RoomID, card.UserID).Update("owned_card_count", gorm.Expr("owned_card_count + 1")).Error
		if err != nil {
			return utils.ErrorMsg(c, utils.ErrInternalDB, utils.Trace(), err.Error(), utils.ErrFromMysqlDB)
		}
	}
	return nil
}
