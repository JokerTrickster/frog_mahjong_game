package repository

import (
	"context"
	"fmt"
	"main/features/ws/model/entity"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func RequestWinFindAllRoomUsers(ctx context.Context, tx *gorm.DB, roomID uint) ([]entity.RoomUsers, error) {
	var roomUsers []entity.RoomUsers
	if err := tx.Preload("User").Preload("Room").Preload("Cards", func(db *gorm.DB) *gorm.DB {
		return db.Where("room_id = ?", roomID).Order("updated_at ASC")
	}).Where("room_id = ?", roomID).Find(&roomUsers).Error; err != nil {
		return nil, fmt.Errorf("room_users 조회 에러: %v", err.Error())
	}
	return roomUsers, nil
}
func RequestWinFindOneDora(c context.Context, tx *gorm.DB, roomID uint) (*mysql.Cards, error) {
	dora := mysql.Cards{}
	err := tx.Model(&mysql.Cards{}).Where("room_id = ? and state = ?", roomID, "dora").First(&dora).Error
	if err != nil {
		return nil, fmt.Errorf("도라 카드를 찾을 수 없습니다. %v", err.Error())
	}
	return &dora, nil
}

// 카드 정보 체크 (소유하고 있는지 체크)
func RequestWinFindAllCards(c context.Context, tx *gorm.DB, requestWinEntity *entity.WSRequestWinEntity) ([]*mysql.Cards, error) {
	cards := make([]*mysql.Cards, 0)
	err := tx.Model(&mysql.Cards{}).Where("room_id = ? and user_id = ? and card_id IN ?", requestWinEntity.RoomID, requestWinEntity.UserID, requestWinEntity.Cards).Find(&cards).Error
	if err != nil {
		return nil, fmt.Errorf("카드를 찾을 수 없습니다. %v", err.Error())
	}
	return cards, nil
}

// 카드 정보 모두 삭제
func RequestWinDeleteAllCards(ctx context.Context, tx *gorm.DB, requestWinEntity *entity.WSRequestWinEntity) error {
	err := tx.Model(&mysql.Cards{}).Where("room_id = ?", requestWinEntity.RoomID).Delete(&mysql.Cards{}).Error
	if err != nil {
		return fmt.Errorf("카드 삭제 실패 %v", err.Error())
	}
	return nil
}

// 유저 상태 변경 (play -> wait)
func RequestWinUpdateRoomUsers(c context.Context, tx *gorm.DB, requestWinEntity *entity.WSRequestWinEntity) error {
	err := tx.Model(&mysql.RoomUsers{}).Where("room_id = ?", requestWinEntity.RoomID).Update("player_state", "wait").Error
	if err != nil {
		return fmt.Errorf("방 유저 상태 변경 실패 %v", err.Error())
	}
	return nil
}

// 론인 경우 해당 유저에 코인 차감한다.
func RequestWinLoanDiffCoin(c context.Context, tx *gorm.DB, requestWinEntity *entity.WSRequestWinEntity) error {
	coinStr := fmt.Sprintf("coin - %d", requestWinEntity.Score)
	err := tx.Model(&mysql.Users{}).Where("id = ?", requestWinEntity.LoanInfo.TargetUserID).Update("coin", gorm.Expr(coinStr)).Error
	if err != nil {
		return fmt.Errorf("유저 코인 차감 실패 %v", err.Error())
	}
	return nil
}

// 론인 경우 해당 유저에 코인 추가한다.
func RequestWinLoanAddCoin(c context.Context, tx *gorm.DB, requestWinEntity *entity.WSRequestWinEntity) error {
	coinStr := fmt.Sprintf("coin + %d", requestWinEntity.Score)
	err := tx.Model(&mysql.Users{}).Where("id = ?", requestWinEntity.UserID).Update("coin", gorm.Expr(coinStr)).Error
	if err != nil {
		return fmt.Errorf("유저 코인 추가 실패 %v", err.Error())
	}
	return nil
}

// 론이 아닌 경우 모든 플레이어에게 점수 차감
func RequestWinDiffCoin(c context.Context, tx *gorm.DB, requestWinEntity *entity.WSRequestWinEntity, coin int) error {
	coinStr := fmt.Sprintf("coin - %d", coin)
	err := tx.Model(&mysql.Users{}).Where("room_id = ? and id != ?", requestWinEntity.RoomID, requestWinEntity.UserID).Update("coin", gorm.Expr(coinStr)).Error
	if err != nil {
		return fmt.Errorf("유저 코인 차감 실패 %v", err.Error())
	}
	return nil
}

// 론이 아닌 경우 해당 유저에 코인 추가한다.
func RequestWinAddCoin(c context.Context, tx *gorm.DB, requestWinEntity *entity.WSRequestWinEntity) error {
	coinStr := fmt.Sprintf("coin + %d", requestWinEntity.Score)
	err := tx.Model(&mysql.Users{}).Where("id = ?", requestWinEntity.UserID).Update("coin", gorm.Expr(coinStr)).Error
	if err != nil {
		return fmt.Errorf("유저 코인 추가 실패 %v", err.Error())
	}
	return nil
}
