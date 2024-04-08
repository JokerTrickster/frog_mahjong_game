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

func NewLoanGameRepository(gormDB *gorm.DB) _interface.ILoanGameRepository {
	return &LoanGameRepository{GormDB: gormDB}
}

func (d *LoanGameRepository) CheckLoan(c context.Context, req *request.ReqLoan) error {
	var card mysql.Cards
	err := d.GormDB.Model(&card).Where("room_id = ? AND state = ? and user_id = ? and id = ?", req.RoomID, "discard", req.LoanUserID, req.LoanCardID).Order("updated_at desc").First(&card).Error
	if err != nil {
		return utils.ErrorMsg(c, utils.ErrNotLoanCard, utils.Trace(), _errors.ErrNotLoanCard.Error(), utils.ErrFromClient)
	}
	return nil
}

// loan 하기 (상대방이 버린 카드를 가져온다)
func (d *LoanGameRepository) Loan(c context.Context, req *request.ReqLoan) error {

	err := d.GormDB.Model(&mysql.Cards{}).Where("room_id = ? AND user_id = ? and id = ? and state = ?", req.RoomID, req.LoanUserID, req.LoanCardID, "discard").Update("state", "loan").Error
	if err != nil {
		return utils.ErrorMsg(c, utils.ErrInternalDB, utils.Trace(), err.Error(), utils.ErrFromMysqlDB)
	}

	return nil
}

// loan한 유저 룸 유저 정보 변경 (카드 수 증가 , 상태 변경 loan)
func (d *LoanGameRepository) UpdateRoomUserCardCount(c context.Context, userID uint, roomID uint) error {
	var roomUser mysql.RoomUsers
	err := d.GormDB.Model(&roomUser).Where("room_id = ? AND user_id = ?", roomID, userID).Updates(map[string]interface{}{"owned_card_count": gorm.Expr("owned_card_count + 1"), "player_state": "loan"}).Error
	if err != nil {
		return utils.ErrorMsg(c, utils.ErrInternalDB, utils.Trace(), err.Error(), utils.ErrFromMysqlDB)
	}
	return nil
}
