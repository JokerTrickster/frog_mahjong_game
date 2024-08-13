package repository

import (
	"context"
	"fmt"
	"main/features/ws/model/entity"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func LoanFindAllRoomUsers(ctx context.Context, tx *gorm.DB, roomID uint) ([]entity.RoomUsers, error) {
	var roomUsers []entity.RoomUsers
	if err := tx.Preload("User").Preload("Room").Preload("Cards", func(db *gorm.DB) *gorm.DB {
		return db.Where("room_id = ?", roomID).Order("updated_at ASC")
	}).Where("room_id = ?", roomID).Find(&roomUsers).Error; err != nil {
		return nil, fmt.Errorf("room_users 조회 에러: %v", err.Error())
	}
	return roomUsers, nil
}
func LoanCardFindOneDora(c context.Context, tx *gorm.DB, roomID uint) (*mysql.Cards, error) {
	dora := mysql.Cards{}
	err := tx.Model(&mysql.Cards{}).Where("room_id = ? and state = ?", roomID, "dora").First(&dora).Error
	if err != nil {
		return nil, fmt.Errorf("도라 카드를 찾을 수 없습니다. %v", err.Error())
	}
	return &dora, nil
}
func LoanCheckLoan(c context.Context, tx *gorm.DB, loanEntity *entity.WSLoanEntity) error {
	var card mysql.Cards
	err := tx.Model(&card).Where("room_id = ? AND state = ? and user_id = ? and card_id = ?", loanEntity.RoomID, "discard", loanEntity.TargetUserID, loanEntity.CardID).Order("updated_at desc").First(&card).Error
	if err != nil {
		return fmt.Errorf("대여할 수 없는 카드입니다. %v", err.Error())
	}
	return nil
}

// loan 하기 (상대방이 버린 카드를 가져온다)
func LoanCardLoan(c context.Context, tx *gorm.DB, loanEntity *entity.WSLoanEntity) error {

	err := tx.Model(&mysql.Cards{}).Where("room_id = ? AND user_id = ? and card_id = ? and state = ?", loanEntity.RoomID, loanEntity.TargetUserID, loanEntity.CardID, "discard").Updates(map[string]interface{}{"user_id": loanEntity.UserID, "state": "owned"}).Error
	if err != nil {
		return fmt.Errorf("카드 대여 실패 %v", err.Error())
	}

	return nil
}

// loan한 유저 룸 유저 정보 변경 (카드 수 증가 , 상태 변경 loan)
func LoanUpdateRoomUserCardCount(c context.Context, tx *gorm.DB, loanEntity *entity.WSLoanEntity) error {
	var roomUser mysql.RoomUsers
	err := tx.Model(&roomUser).Where("room_id = ? AND user_id = ?", loanEntity.RoomID, loanEntity.UserID).Updates(map[string]interface{}{"owned_card_count": gorm.Expr("owned_card_count + 1"), "player_state": "loan"}).Error
	if err != nil {
		return fmt.Errorf("방 유저 카드 카운트 업데이트 실패 %v", err.Error())
	}
	return nil
}
