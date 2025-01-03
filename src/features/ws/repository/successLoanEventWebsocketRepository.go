package repository

import (
	"context"
	"fmt"
	"main/features/ws/model/entity"
	"main/utils/db/mysql"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SuccessFindAllRoomUsers(ctx context.Context, tx *gorm.DB, roomID uint) ([]entity.RoomUsers, error) {
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
func SuccessFindOneDora(ctx context.Context, tx *gorm.DB, roomID uint) (*mysql.FrogUserCards, error) {
	var dora mysql.FrogUserCards
	err := tx.Model(&mysql.FrogUserCards{}).
		Where("room_id = ?", roomID).
		Where("state = ?", "dora").
		First(&dora).Error
	if err != nil {
		return nil, fmt.Errorf("도라 카드를 찾을 수 없습니다: %v", err)
	}
	return &dora, nil
}

func SuccessFindAllCards(ctx context.Context, tx *gorm.DB, successEntity *entity.WSSuccessEntity) ([]*mysql.FrogUserCards, error) {
	var cards []*mysql.FrogUserCards
	err := tx.Model(&mysql.FrogUserCards{}).
		Where("room_id = ?", successEntity.RoomID).
		Where("user_id = ?", successEntity.UserID).
		Where("card_id IN ?", successEntity.Cards).
		Find(&cards).Error
	if err != nil {
		return nil, fmt.Errorf("카드를 찾을 수 없습니다: %v", err)
	}
	return cards, nil
}

func SuccessDeleteAllCards(ctx context.Context, tx *gorm.DB, successEntity *entity.WSSuccessEntity) error {
	err := tx.Model(&mysql.FrogUserCards{}).
		Where("room_id = ?", successEntity.RoomID).
		Delete(&mysql.FrogUserCards{}).Error
	if err != nil {
		return fmt.Errorf("카드 삭제 실패: %v", err)
	}
	return nil
}

func SuccessUpdateRoomUsers(ctx context.Context, tx *gorm.DB, successEntity *entity.WSSuccessEntity) error {
	err := tx.Model(&mysql.FrogRoomUsers{}).
		Where("room_id = ?", successEntity.RoomID).
		Update("player_state", "wait").Error
	if err != nil {
		return fmt.Errorf("방 유저 상태 변경 실패: %v", err)
	}
	return nil
}

func SuccessLoanDiffCoin(ctx context.Context, tx *gorm.DB, successEntity *entity.WSSuccessEntity) error {
	coinExpr := fmt.Sprintf("coin - %d", successEntity.Score)
	err := tx.Model(&mysql.Users{}).
		Where("id = ?", successEntity.LoanInfo.TargetUserID).
		Update("coin", gorm.Expr(coinExpr)).Error
	if err != nil {
		return fmt.Errorf("유저 코인 차감 실패: %v", err)
	}
	return nil
}

func SuccessLoanAddCoin(ctx context.Context, tx *gorm.DB, successEntity *entity.WSSuccessEntity) error {
	coinExpr := fmt.Sprintf("coin + %d", successEntity.Score)
	err := tx.Model(&mysql.Users{}).
		Where("id = ?", successEntity.UserID).
		Update("coin", gorm.Expr(coinExpr)).Error
	if err != nil {
		return fmt.Errorf("유저 코인 추가 실패: %v", err)
	}
	return nil
}
