package repository

import (
	"context"
	"fmt"
	"main/features/ws/model/entity"
	_errors "main/features/ws/model/errors"
	"main/utils/db/mysql"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// FailedLoanFindAllRoomUsers retrieves all room users with necessary preloads
func FailedLoanFindAllRoomUsers(ctx context.Context, tx *gorm.DB, roomID uint) ([]entity.RoomUsers, *entity.ErrorInfo) {
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

// FailedLoanCheckCard verifies if the user owns the specified card
func FailedLoanCheckCard(ctx context.Context, tx *gorm.DB, loanEntity *entity.WSLoanEntity) *entity.ErrorInfo {
	var card mysql.FrogUserCards
	err := tx.Model(&mysql.FrogUserCards{}).
		Where("room_id = ?", loanEntity.RoomID).
		Where("state = ?", "owned").
		Where("user_id = ?", loanEntity.UserID).
		Where("card_id = ?", loanEntity.CardID).
		Order("updated_at desc").
		First(&card).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeNotFound, // 404
			Msg:  "카드 소유 여부 확인 실패",
			Type: _errors.ErrNotFoundCard,
		}
	}
	return nil
}

// FailedLoanRollbackCard rolls back the card to the target user
func FailedLoanRollbackCard(ctx context.Context, tx *gorm.DB, loanEntity *entity.WSLoanEntity) *entity.ErrorInfo {
	err := tx.Model(&mysql.FrogUserCards{}).
		Where("room_id = ?", loanEntity.RoomID).
		Where("user_id = ?", loanEntity.UserID).
		Where("card_id = ?", loanEntity.CardID).
		Where("state = ?", "owned").
		Updates(map[string]interface{}{
			"user_id": loanEntity.TargetUserID,
			"state":   "discard",
		}).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal, // 500
			Msg:  "카드 롤백 실패",
			Type: _errors.ErrUpdateFailed,
		}
	}
	return nil
}

// FailedLoanPenalty deducts a penalty coin from the user
func FailedLoanPenalty(ctx context.Context, tx *gorm.DB, loanEntity *entity.WSLoanEntity, penaltyCoin int) *entity.ErrorInfo {
	penaltyStr := fmt.Sprintf("coin - %d", penaltyCoin)
	err := tx.Model(&mysql.Users{}).
		Where("id = ?", loanEntity.UserID).
		Updates(map[string]interface{}{
			"coin": gorm.Expr(penaltyStr),
		}).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal, // 500
			Msg:  "패널티 부여 실패",
			Type: _errors.ErrUpdateFailed,
		}
	}
	return nil
}

// FailedLoanAddCoin adds coins to all players except the loan user
func FailedLoanAddCoin(ctx context.Context, tx *gorm.DB, loanEntity *entity.WSLoanEntity) *entity.ErrorInfo {
	err := tx.Model(&mysql.Users{}).
		Where("id != ?", loanEntity.UserID).
		Where("room_id = ?", loanEntity.RoomID).
		Updates(map[string]interface{}{
			"coin": gorm.Expr("coin + 2"),
		}).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal, // 500
			Msg:  "플레이어 코인 추가 실패",
			Type: _errors.ErrUpdateFailed,
		}
	}
	return nil
}
