package repository

import (
	"context"
	"main/features/frog/model/entity"
	"main/utils/db/mysql"

	_errors "main/features/frog/model/errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// LoanFindAllRoomUsers retrieves all room users with necessary preloads
func LoanFindAllRoomUsers(ctx context.Context, tx *gorm.DB, roomID uint) ([]entity.RoomUsers, *entity.ErrorInfo) {
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

// LoanCardFindOneDora retrieves the "dora" card for a specific room
func LoanCardFindOneDora(ctx context.Context, tx *gorm.DB, roomID uint) (*mysql.FrogUserCards, *entity.ErrorInfo) {
	var dora mysql.FrogUserCards
	err := tx.Model(&mysql.FrogUserCards{}).
		Where("room_id = ?", roomID).
		Where("state = ?", "dora").
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

// LoanCheckLoan checks if a loan is possible for a specific card
func LoanCheckLoan(ctx context.Context, tx *gorm.DB, loanEntity *entity.WSLoanEntity) *entity.ErrorInfo {
	var card mysql.FrogUserCards
	err := tx.Model(&mysql.FrogUserCards{}).
		Where("room_id = ?", loanEntity.RoomID).
		Where("state = ?", "discard").
		Where("user_id = ?", loanEntity.TargetUserID).
		Where("card_id = ?", loanEntity.CardID).
		Order("updated_at desc").
		First(&card).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeForbidden, // 403
			Msg:  "대여할 수 없는 카드입니다.",
			Type: _errors.ErrItemNotAvailable,
		}
	}
	return nil
}

// LoanCardLoan processes a card loan
func LoanCardLoan(ctx context.Context, tx *gorm.DB, loanEntity *entity.WSLoanEntity) *entity.ErrorInfo {
	err := tx.Model(&mysql.FrogUserCards{}).
		Where("room_id = ?", loanEntity.RoomID).
		Where("user_id = ?", loanEntity.TargetUserID).
		Where("card_id = ?", loanEntity.CardID).
		Where("state = ?", "discard").
		Updates(map[string]interface{}{
			"user_id": loanEntity.UserID,
			"state":   "owned",
		}).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal, // 500
			Msg:  "카드 대여 실패",
			Type: _errors.ErrUpdateFailed,
		}
	}
	return nil
}

// LoanUpdateRoomUserCardCount updates the room user's card count and state
func LoanUpdateRoomUserCardCount(ctx context.Context, tx *gorm.DB, loanEntity *entity.WSLoanEntity) *entity.ErrorInfo {
	err := tx.Model(&mysql.FrogRoomUsers{}).
		Where("room_id = ?", loanEntity.RoomID).
		Where("user_id = ?", loanEntity.UserID).
		Updates(map[string]interface{}{
			"owned_card_count": gorm.Expr("owned_card_count + 1"),
			"player_state":     "loan",
		}).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal, // 500
			Msg:  "방 유저 카드 카운트 업데이트 실패",
			Type: _errors.ErrUpdateFailed,
		}
	}
	return nil
}
