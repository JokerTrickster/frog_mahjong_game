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

// SuccessFindAllRoomUsers retrieves all room users with necessary preloads
func SuccessFindAllRoomUsers(ctx context.Context, tx *gorm.DB, roomID uint) ([]entity.RoomUsers, *entity.ErrorInfo) {
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
			Msg:  "방 유저를 찾을 수 없습니다",
			Type: _errors.ErrRoomUsersNotFound,
		}
	}
	return roomUsers, nil
}

// SuccessFindOneDora retrieves the "dora" card for a specific room
func SuccessFindOneDora(ctx context.Context, tx *gorm.DB, roomID uint) (*mysql.FrogUserCards, *entity.ErrorInfo) {
	var dora mysql.FrogUserCards
	err := tx.Model(&mysql.FrogUserCards{}).
		Where("room_id = ?", roomID).
		Where("state = ?", "dora").
		First(&dora).Error
	if err != nil {
		return nil, &entity.ErrorInfo{
			Code: _errors.ErrCodeNotFound, // 404
			Msg:  "도라 카드를 찾을 수 없습니다",
			Type: _errors.ErrNotFoundCard,
		}
	}
	return &dora, nil
}

// SuccessFindAllCards retrieves all cards for a specific room and user
func SuccessFindAllCards(ctx context.Context, tx *gorm.DB, successEntity *entity.WSSuccessEntity) ([]*mysql.FrogUserCards, *entity.ErrorInfo) {
	var cards []*mysql.FrogUserCards
	err := tx.Model(&mysql.FrogUserCards{}).
		Where("room_id = ?", successEntity.RoomID).
		Where("user_id = ?", successEntity.UserID).
		Where("card_id IN ?", successEntity.Cards).
		Find(&cards).Error
	if err != nil {
		return nil, &entity.ErrorInfo{
			Code: _errors.ErrCodeNotFound, // 404
			Msg:  "카드를 찾을 수 없습니다",
			Type: _errors.ErrNotFoundCard,
		}
	}
	return cards, nil
}

// SuccessDeleteAllCards deletes all cards for a specific room
func SuccessDeleteAllCards(ctx context.Context, tx *gorm.DB, successEntity *entity.WSSuccessEntity) *entity.ErrorInfo {
	err := tx.Model(&mysql.FrogUserCards{}).
		Where("room_id = ?", successEntity.RoomID).
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

// SuccessUpdateRoomUsers updates the player state of room users
func SuccessUpdateRoomUsers(ctx context.Context, tx *gorm.DB, successEntity *entity.WSSuccessEntity) *entity.ErrorInfo {
	err := tx.Model(&mysql.FrogRoomUsers{}).
		Where("room_id = ?", successEntity.RoomID).
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

// SuccessLoanDiffCoin subtracts coins from a specific user
func SuccessLoanDiffCoin(ctx context.Context, tx *gorm.DB, successEntity *entity.WSSuccessEntity) *entity.ErrorInfo {
	coinExpr := fmt.Sprintf("coin - %d", successEntity.Score)
	err := tx.Model(&mysql.Users{}).
		Where("id = ?", successEntity.LoanInfo.TargetUserID).
		Update("coin", gorm.Expr(coinExpr)).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal, // 500
			Msg:  "유저 코인 차감 실패",
			Type: _errors.ErrUpdateFailed,
		}
	}
	return nil
}

// SuccessLoanAddCoin adds coins to a specific user
func SuccessLoanAddCoin(ctx context.Context, tx *gorm.DB, successEntity *entity.WSSuccessEntity) *entity.ErrorInfo {
	coinExpr := fmt.Sprintf("coin + %d", successEntity.Score)
	err := tx.Model(&mysql.Users{}).
		Where("id = ?", successEntity.UserID).
		Update("coin", gorm.Expr(coinExpr)).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal, // 500
			Msg:  "유저 코인 추가 실패",
			Type: _errors.ErrUpdateFailed,
		}
	}
	return nil
}
