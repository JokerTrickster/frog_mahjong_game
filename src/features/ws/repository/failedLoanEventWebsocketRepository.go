package repository

import (
	"context"
	"fmt"
	"main/features/ws/model/entity"
	"main/utils/db/mysql"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func FailedLoanFindAllRoomUsers(ctx context.Context, tx *gorm.DB, roomID uint) ([]entity.RoomUsers, error) {
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

// 소유하고 있는 카드인지 체크
func FailedLoanCheckCard(ctx context.Context, tx *gorm.DB, loanEntity *entity.WSLoanEntity) error {
	var card mysql.FrogUserCards
	err := tx.Model(&mysql.FrogUserCards{}).
		Where("room_id = ?", loanEntity.RoomID).
		Where("state = ?", "owned").
		Where("user_id = ?", loanEntity.UserID).
		Where("card_id = ?", loanEntity.CardID).
		Order("updated_at desc").
		First(&card).Error
	if err != nil {
		return fmt.Errorf("카드 소유 여부 확인 실패: %v", err)
	}
	return nil
}

// 카드 정보를 롤백한다.
func FailedLoanRollbackCard(ctx context.Context, tx *gorm.DB, loanEntity *entity.WSLoanEntity) error {
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
		return fmt.Errorf("카드 롤백 실패: %v", err)
	}
	return nil
}

// 패널티를 부여한다. (코인 차감)
func FailedLoanPenalty(ctx context.Context, tx *gorm.DB, loanEntity *entity.WSLoanEntity, penaltyCoin int) error {
	penaltyStr := fmt.Sprintf("coin - %d", penaltyCoin)
	err := tx.Model(&mysql.Users{}).
		Where("id = ?", loanEntity.UserID).
		Updates(map[string]interface{}{
			"coin": gorm.Expr(penaltyStr),
		}).Error
	if err != nil {
		return fmt.Errorf("패널티 부여 실패: %v", err)
	}
	return nil
}

// 모든 플레이어에게 코인 추가
func FailedLoanAddCoin(ctx context.Context, tx *gorm.DB, loanEntity *entity.WSLoanEntity) error {
	err := tx.Model(&mysql.Users{}).
		Where("id != ?", loanEntity.UserID).
		Where("room_id = ?", loanEntity.RoomID).
		Updates(map[string]interface{}{
			"coin": gorm.Expr("coin + 2"),
		}).Error
	if err != nil {
		return fmt.Errorf("플레이어 코인 추가 실패: %v", err)
	}
	return nil
}
