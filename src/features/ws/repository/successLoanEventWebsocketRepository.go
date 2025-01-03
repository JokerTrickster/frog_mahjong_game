package repository

import (
	"context"
	"fmt"
	"main/features/ws/model/entity"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func SuccessFindAllRoomUsers(ctx context.Context, tx *gorm.DB, roomID uint) ([]entity.RoomUsers, error) {
	var roomUsers []entity.RoomUsers
	if err := tx.Preload("User").Preload("Room").Preload("Cards", func(db *gorm.DB) *gorm.DB {
		return db.Where("room_id = ?", roomID).Order("updated_at ASC")
	}).Where("room_id = ?", roomID).Find(&roomUsers).Error; err != nil {
		return nil, fmt.Errorf("room_users 조회 에러: %v", err.Error())
	}
	return roomUsers, nil
}
func SuccessFindOneDora(c context.Context, tx *gorm.DB, roomID uint) (*mysql.FrogUserCards, error) {
	dora := mysql.FrogUserCards{}
	err := tx.Model(&mysql.FrogUserCards{}).Where("room_id = ? and state = ?", roomID, "dora").First(&dora).Error
	if err != nil {
		return nil, fmt.Errorf("도라 카드를 찾을 수 없습니다. %v", err.Error())
	}
	return &dora, nil
}

// 카드 정보 체크 (소유하고 있는지 체크)
func SuccessFindAllCards(c context.Context, tx *gorm.DB, SuccessEntity *entity.WSSuccessEntity) ([]*mysql.FrogUserCards, error) {
	cards := make([]*mysql.FrogUserCards, 0)
	err := tx.Model(&mysql.FrogUserCards{}).Where("room_id = ? and user_id = ? and card_id IN ?", SuccessEntity.RoomID, SuccessEntity.UserID, SuccessEntity.Cards).Find(&cards).Error
	if err != nil {
		return nil, fmt.Errorf("카드를 찾을 수 없습니다. %v", err.Error())
	}
	return cards, nil
}

// 카드 정보 모두 삭제
func SuccessDeleteAllCards(ctx context.Context, tx *gorm.DB, SuccessEntity *entity.WSSuccessEntity) error {
	err := tx.Model(&mysql.FrogUserCards{}).Where("room_id = ?", SuccessEntity.RoomID).Delete(&mysql.FrogUserCards{}).Error
	if err != nil {
		return fmt.Errorf("카드 삭제 실패 %v", err.Error())
	}
	return nil
}

// 유저 상태 변경 (play -> wait)
func SuccessUpdateRoomUsers(c context.Context, tx *gorm.DB, SuccessEntity *entity.WSSuccessEntity) error {
	err := tx.Model(&mysql.FrogRoomUsers{}).Where("room_id = ?", SuccessEntity.RoomID).Update("player_state", "wait").Error
	if err != nil {
		return fmt.Errorf("방 유저 상태 변경 실패 %v", err.Error())
	}
	return nil
}

// 론인 경우 해당 유저에 코인 차감한다.
func SuccessLoanDiffCoin(c context.Context, tx *gorm.DB, SuccessEntity *entity.WSSuccessEntity) error {
	coinStr := fmt.Sprintf("coin - %d", SuccessEntity.Score)
	err := tx.Model(&mysql.Users{}).Where("id = ?", SuccessEntity.LoanInfo.TargetUserID).Update("coin", gorm.Expr(coinStr)).Error
	if err != nil {
		return fmt.Errorf("유저 코인 차감 실패 %v", err.Error())
	}
	return nil
}

// 론인 경우 해당 유저에 코인 추가한다.
func SuccessLoanAddCoin(c context.Context, tx *gorm.DB, SuccessEntity *entity.WSSuccessEntity) error {
	coinStr := fmt.Sprintf("coin + %d", SuccessEntity.Score)
	err := tx.Model(&mysql.Users{}).Where("id = ?", SuccessEntity.UserID).Update("coin", gorm.Expr(coinStr)).Error
	if err != nil {
		return fmt.Errorf("유저 코인 추가 실패 %v", err.Error())
	}
	return nil
}
