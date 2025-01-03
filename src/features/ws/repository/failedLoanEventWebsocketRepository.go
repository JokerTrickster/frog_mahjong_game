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
func FailedLoanCardFindOneDora(c context.Context, tx *gorm.DB, roomID uint) (*mysql.FrogUserCards, error) {
	dora := mysql.FrogUserCards{}
	err := tx.Model(&mysql.FrogUserCards{}).Where("room_id = ? and state = ?", roomID, "dora").First(&dora).Error
	if err != nil {
		return nil, fmt.Errorf("도라 카드를 찾을 수 없습니다. %v", err.Error())
	}
	return &dora, nil
}

// 소유하고 있는 카드인지 체크
func FailedLoanCheckCard(c context.Context, tx *gorm.DB, loanEntity *entity.WSLoanEntity) error {
	var card mysql.FrogUserCards
	err := tx.Model(&card).Where("room_id = ? AND state = ? and user_id = ? and card_id = ?", loanEntity.RoomID, "owned", loanEntity.UserID, loanEntity.CardID).Order("updated_at desc").First(&card).Error
	if err != nil {
		return fmt.Errorf("카드 소유 여부 확인 실패 %v", err.Error())
	}
	return nil
}

// 카드 정보를 롤백한다.
func FailedLoanRollbackCard(c context.Context, tx *gorm.DB, loanEntity *entity.WSLoanEntity) error {

	err := tx.Model(&mysql.FrogUserCards{}).Where("room_id = ? AND user_id = ? and card_id = ? and state = ?", loanEntity.RoomID, loanEntity.UserID, loanEntity.CardID, "owned").Updates(map[string]interface{}{"user_id": loanEntity.TargetUserID, "state": "discard"}).Error
	if err != nil {
		return fmt.Errorf("카드 대여 실패 %v", err.Error())
	}

	return nil
}

// 패널티를 부여한다. (코인 차감)
func FailedLoanPenalty(c context.Context, tx *gorm.DB, loanEntity *entity.WSLoanEntity, penaltyCoin int) error {
	penaltyStr := fmt.Sprintf("coin - %d", penaltyCoin)
	var user mysql.Users
	err := tx.Model(&user).Where("id = ?", loanEntity.UserID).Updates(map[string]interface{}{"coin": gorm.Expr(penaltyStr)}).Error
	if err != nil {
		return fmt.Errorf("방 유저 카드 카운트 업데이트 실패 %v", err.Error())
	}
	return nil
}

// 모든 플레이어에게 코인 추가
func FailedLoanAddCoin(c context.Context, tx *gorm.DB, loanEntity *entity.WSLoanEntity) error {
	var user mysql.Users
	err := tx.Model(&user).Where("id != ? AND room_id = ?", loanEntity.UserID, loanEntity.RoomID).
		Updates(map[string]interface{}{"coin": gorm.Expr("coin + 2")}).Error
	if err != nil {
		return fmt.Errorf("방 유저 코인 업데이트 실패 %v", err.Error())
	}
	return nil
}
