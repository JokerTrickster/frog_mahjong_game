package repository

import (
	"context"
	"fmt"
	"main/features/frog/model/entity"
	_errors "main/features/frog/model/errors"
	"main/utils/db/mysql"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// RequestWinFindAllRoomUsers retrieves all room users with necessary preloads
func RequestWinFindAllRoomUsers(ctx context.Context, tx *gorm.DB, roomID uint) ([]entity.RoomUsers, *entity.ErrorInfo) {
	var roomUsers []entity.RoomUsers
	if err := tx.Table("frog_room_users").
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("room_id = ?", roomID).
		Preload("User").
		Preload("Room").
		Preload("Cards", func(db *gorm.DB) *gorm.DB {
			return db.Where("room_id = ?", roomID).Order("updated_at ASC")
		}).
		Find(&roomUsers).Error; err != nil {
		return nil, &entity.ErrorInfo{
			Code: _errors.ErrCodeNotFound, // 404
			Msg:  "room_users 조회 실패",
			Type: _errors.ErrRoomUsersNotFound,
		}
	}
	return roomUsers, nil
}

// RequestWinFindOneDora retrieves the "dora" card for a specific room
func RequestWinFindOneDora(c context.Context, tx *gorm.DB, roomID uint) (*mysql.FrogUserCards, *entity.ErrorInfo) {
	dora := mysql.FrogUserCards{}
	err := tx.Model(&mysql.FrogUserCards{}).
		Where("room_id = ? and state = ?", roomID, "dora").
		First(&dora).Error
	if err != nil {
		return nil, &entity.ErrorInfo{
			Code: _errors.ErrCodeNotFound, // 404
			Msg:  "도라 카드를 찾을 수 없습니다.",
			Type: _errors.ErrNotFoundCard,
		}
	}
	return &dora, nil
}

// RequestWinFindAllCards checks card ownership
func RequestWinFindAllCards(c context.Context, tx *gorm.DB, requestWinEntity *entity.WSRequestWinEntity) ([]*mysql.FrogUserCards, *entity.ErrorInfo) {
	cards := make([]*mysql.FrogUserCards, 0)
	err := tx.Model(&mysql.FrogUserCards{}).
		Where("room_id = ? and user_id = ? and card_id IN ?", requestWinEntity.RoomID, requestWinEntity.UserID, requestWinEntity.Cards).
		Find(&cards).Error
	if err != nil {
		return nil, &entity.ErrorInfo{
			Code: _errors.ErrCodeNotFound, // 404
			Msg:  "카드를 찾을 수 없습니다.",
			Type: _errors.ErrNotFoundCard,
		}
	}
	return cards, nil
}

// RequestWinDeleteAllCards deletes all cards in the room
func RequestWinDeleteAllCards(ctx context.Context, tx *gorm.DB, requestWinEntity *entity.WSRequestWinEntity) *entity.ErrorInfo {
	err := tx.Model(&mysql.FrogUserCards{}).
		Where("room_id = ?", requestWinEntity.RoomID).
		Delete(&mysql.FrogUserCards{}).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal, // 500
			Msg:  "카드 삭제 실패",
			Type: _errors.ErrDeleteFailed,
		}
	}
	return nil
}

// RequestWinUpdateRoomUsers updates the player state to "wait"
func RequestWinUpdateRoomUsers(c context.Context, tx *gorm.DB, requestWinEntity *entity.WSRequestWinEntity) *entity.ErrorInfo {
	err := tx.Model(&mysql.FrogRoomUsers{}).
		Where("room_id = ?", requestWinEntity.RoomID).
		Update("player_state", "wait").Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal, // 500
			Msg:  "방 유저 상태 변경 실패",
			Type: _errors.ErrUpdateUserStateFailed,
		}
	}
	return nil
}

// RequestWinLoanDiffCoin subtracts coins from the target user in loan case
func RequestWinLoanDiffCoin(c context.Context, tx *gorm.DB, requestWinEntity *entity.WSRequestWinEntity) *entity.ErrorInfo {
	coinStr := fmt.Sprintf("coin - %d", requestWinEntity.Score)
	err := tx.Model(&mysql.GameUsers{}).
		Where("id = ?", requestWinEntity.LoanInfo.TargetUserID).
		Update("coin", gorm.Expr(coinStr)).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal, // 500
			Msg:  "유저 코인 차감 실패",
			Type: _errors.ErrUpdateFailed,
		}
	}
	return nil
}

// RequestWinLoanAddCoin adds coins to the winning user in loan case
func RequestWinLoanAddCoin(c context.Context, tx *gorm.DB, requestWinEntity *entity.WSRequestWinEntity) *entity.ErrorInfo {
	coinStr := fmt.Sprintf("coin + %d", requestWinEntity.Score)
	err := tx.Model(&mysql.GameUsers{}).
		Where("id = ?", requestWinEntity.UserID).
		Update("coin", gorm.Expr(coinStr)).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal, // 500
			Msg:  "유저 코인 추가 실패",
			Type: _errors.ErrUpdateFailed,
		}
	}
	return nil
}

// RequestWinDiffCoin subtracts coins from all players except the winner
func RequestWinDiffCoin(c context.Context, tx *gorm.DB, requestWinEntity *entity.WSRequestWinEntity, coin int) *entity.ErrorInfo {
	coinStr := fmt.Sprintf("coin - %d", coin)
	err := tx.Model(&mysql.GameUsers{}).
		Where("room_id = ? and id != ?", requestWinEntity.RoomID, requestWinEntity.UserID).
		Update("coin", gorm.Expr(coinStr)).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal, // 500
			Msg:  "유저 코인 차감 실패",
			Type: _errors.ErrUpdateFailed,
		}
	}
	return nil
}

// RequestWinAddCoin adds coins to the winner in non-loan case
func RequestWinAddCoin(c context.Context, tx *gorm.DB, requestWinEntity *entity.WSRequestWinEntity) *entity.ErrorInfo {
	coinStr := fmt.Sprintf("coin + %d", requestWinEntity.Score)
	err := tx.Model(&mysql.GameUsers{}).
		Where("id = ?", requestWinEntity.UserID).
		Update("coin", gorm.Expr(coinStr)).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal, // 500
			Msg:  "유저 코인 추가 실패",
			Type: _errors.ErrUpdateFailed,
		}
	}
	return nil
}
