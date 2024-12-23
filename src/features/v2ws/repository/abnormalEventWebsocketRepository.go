package repository

import (
	"context"
	"fmt"
	"main/features/v2ws/model/entity"
	_errors "main/features/v2ws/model/errors"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func AbnormalFindAllRoomUsers(ctx context.Context, tx *gorm.DB, roomID uint) ([]entity.RoomUsers, *entity.ErrorInfo) {
	var roomUsers []entity.RoomUsers
	if err := tx.Preload("User").Preload("Room").Preload("RoomUsers").Preload("Cards", func(db *gorm.DB) *gorm.DB {
		return db.Where("room_id = ?", roomID).Order("updated_at ASC")
	}).Where("room_id = ?", roomID).Find(&roomUsers).Error; err != nil {
		return nil, &entity.ErrorInfo{Code: _errors.ErrCodeInternal, Msg: fmt.Sprintf("room_users 조회 에러: %v", err), Type: _errors.ErrInternalServer}
	}
	return roomUsers, nil
}

// 카드 정보 모두 삭제
func AbnormalDeleteAllCards(ctx context.Context, tx *gorm.DB, AbnormalEntity *entity.WSAbnormalEntity) *entity.ErrorInfo {
	err := tx.Model(&mysql.Cards{}).Where("room_id = ?", AbnormalEntity.RoomID).Delete(&mysql.Cards{}).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("카드 삭제 실패: %v", err),
			Type:_errors. ErrDeleteCardFailed,
		}
	}
	return nil
}

// 방 삭제 처리
func AbnormalDeleteRoom(ctx context.Context, tx *gorm.DB, AbnormalEntity *entity.WSAbnormalEntity) *entity.ErrorInfo {
	err := tx.Model(&mysql.Rooms{}).Where("id = ?", AbnormalEntity.RoomID).Delete(&mysql.Rooms{}).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("방 삭제 실패: %v", err),
			Type: _errors.ErrDeleteRoomFailed,
		}
	}
	return nil
}

// 룸 유저 상태 변경 (play -> disconnected)
func AbnormalUpdateRoomUsers(ctx context.Context, tx *gorm.DB, AbnormalEntity *entity.WSAbnormalEntity) *entity.ErrorInfo {
	err := tx.Model(&mysql.RoomUsers{}).Where("room_id = ? and user_id = ?", AbnormalEntity.RoomID, AbnormalEntity.UserID).Update("player_state", "disconnected").Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("유저 상태 변경 실패: %v", err),
			Type: _errors.ErrUpdateUserStateFailed,
		}
	}
	return nil
}
