package repository

import (
	"context"
	"fmt"
	"main/features/ws/model/entity"
	"main/utils/db/mysql"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func LoanFindAllRoomUsers(ctx context.Context, tx *gorm.DB, roomID uint) ([]entity.RoomUsers, error) {
	var roomUsers []entity.RoomUsers
	if err := tx.Table("frog_room_users").Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("room_id = ?", roomID).
		Preload("User").
		Preload("Room").
		Preload("Cards", func(db *gorm.DB) *gorm.DB {
			return db.Where("room_id = ?", roomID).Order("updated_at ASC")
		}).Where("room_id = ?", roomID).Find(&roomUsers).Error; err != nil {
		return nil, fmt.Errorf("room_users 조회 실패: %v", err.Error())
	}
	return roomUsers, nil
}
func LoanCardFindOneDora(ctx context.Context, tx *gorm.DB, roomID uint) (*mysql.FrogUserCards, error) {
	var dora mysql.FrogUserCards
	err := tx.Model(&mysql.FrogUserCards{}).
		Where("room_id = ?", roomID).
		Where("state = ?", "dora").
		First(&dora).Error
	if err != nil {
		return nil, fmt.Errorf("도라 카드를 찾을 수 없습니다: %v", err.Error())
	}
	return &dora, nil
}

func LoanCheckLoan(ctx context.Context, tx *gorm.DB, loanEntity *entity.WSLoanEntity) error {
	var card mysql.FrogUserCards
	err := tx.Model(&mysql.FrogUserCards{}).
		Where("room_id = ?", loanEntity.RoomID).
		Where("state = ?", "discard").
		Where("user_id = ?", loanEntity.TargetUserID).
		Where("card_id = ?", loanEntity.CardID).
		Order("updated_at desc").
		First(&card).Error
	if err != nil {
		return fmt.Errorf("대여할 수 없는 카드입니다: %v", err.Error())
	}
	return nil
}

func LoanCardLoan(ctx context.Context, tx *gorm.DB, loanEntity *entity.WSLoanEntity) error {
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
		return fmt.Errorf("카드 대여 실패: %v", err.Error())
	}
	return nil
}

func LoanUpdateRoomUserCardCount(ctx context.Context, tx *gorm.DB, loanEntity *entity.WSLoanEntity) error {
	err := tx.Model(&mysql.FrogRoomUsers{}).
		Where("room_id = ?", loanEntity.RoomID).
		Where("user_id = ?", loanEntity.UserID).
		Updates(map[string]interface{}{
			"owned_card_count": gorm.Expr("owned_card_count + 1"),
			"player_state":     "loan",
		}).Error
	if err != nil {
		return fmt.Errorf("방 유저 카드 카운트 업데이트 실패: %v", err.Error())
	}
	return nil
}
