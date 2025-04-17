package repository

import (
	"context"
	"fmt"
	"main/features/slime_war/model/entity"
	_errors "main/features/slime_war/model/errors"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

// 카드 정보 모두 삭제
func AbnormalDeleteAllCards(ctx context.Context, tx *gorm.DB, AbnormalEntity *entity.WSAbnormalEntity) *entity.ErrorInfo {
	err := tx.Model(&mysql.UserBirdCards{}).Where("room_id = ?", AbnormalEntity.RoomID).Delete(&mysql.UserBirdCards{}).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("카드 삭제 실패: %v", err),
			Type: _errors.ErrDeleteCardFailed,
		}
	}
	return nil
}

// 방 삭제 처리
func AbnormalDeleteRoom(ctx context.Context, tx *gorm.DB, AbnormalEntity *entity.WSAbnormalEntity) *entity.ErrorInfo {
	err := tx.Model(&mysql.GameRooms{}).Where("id = ?", AbnormalEntity.RoomID).Delete(&mysql.GameRooms{}).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("방 삭제 실패: %v", err),
			Type: _errors.ErrDeleteRoomFailed,
		}
	}
	return nil
}

// 방 유저 정보 삭제
func AbnormalDeleteRoomUsers(ctx context.Context, tx *gorm.DB, AbnormalEntity *entity.WSAbnormalEntity) *entity.ErrorInfo {
	err := tx.Model(&mysql.GameRoomUsers{}).Where("room_id = ?", AbnormalEntity.RoomID).Delete(&mysql.GameRoomUsers{}).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("방 유저 삭제 실패: %v", err),
			Type: _errors.ErrDeleteRoomUserFailed,
		}
	}
	return nil
}

// 룸 유저 상태 변경 (play -> disconnected)
func AbnormalUpdateRoomUsers(ctx context.Context, tx *gorm.DB, AbnormalEntity *entity.WSAbnormalEntity) *entity.ErrorInfo {
	err := tx.Model(&mysql.GameRoomUsers{}).Where("room_id = ? and user_id = ?", AbnormalEntity.RoomID, AbnormalEntity.UserID).Update("player_state", "disconnected").Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("유저 상태 변경 실패: %v", err),
			Type: _errors.ErrUpdateUserStateFailed,
		}
	}
	return nil
}
