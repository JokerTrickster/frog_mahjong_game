package repository

import (
	"context"
	"fmt"
	"main/features/slime_war/model/entity"
	_errors "main/features/slime_war/model/errors"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func AbnormalDeleteAllCards(ctx context.Context, tx *gorm.DB, AbnormalEntity *entity.WSAbnormalEntity) *entity.ErrorInfo {
	err := tx.Model(&mysql.SlimeWarRoomCards{}).Where("room_id = ?", AbnormalEntity.RoomID).Delete(&mysql.SlimeWarRoomCards{}).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("카드 삭제 실패: %v", err),
			Type: _errors.ErrDeleteCardFailed,
		}
	}
	return nil
}

func AbnormalDeleteAllMaps(ctx context.Context, tx *gorm.DB, AbnormalEntity *entity.WSAbnormalEntity) *entity.ErrorInfo {
	err := tx.Model(&mysql.SlimeWarRoomMaps{}).Where("room_id = ?", AbnormalEntity.RoomID).Delete(&mysql.SlimeWarRoomMaps{}).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("맵 삭제 실패: %v", err),
			Type: _errors.ErrDeleteMapFailed,
		}
	}
	return nil
}

func AbnormalDeleteRoomUsers(ctx context.Context, tx *gorm.DB, AbnormalEntity *entity.WSAbnormalEntity) *entity.ErrorInfo {
	err := tx.Model(&mysql.SlimeWarUsers{}).Where("room_id = ?", AbnormalEntity.RoomID).Delete(&mysql.SlimeWarUsers{}).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("방 유저 삭제 실패: %v", err),
			Type: _errors.ErrDeleteRoomUserFailed,
		}
	}
	return nil
}

func AbnormalDeleteGameRoomSetting(ctx context.Context, tx *gorm.DB, AbnormalEntity *entity.WSAbnormalEntity) *entity.ErrorInfo {
	err := tx.Model(&mysql.SlimeWarGameRoomSettings{}).Where("room_id = ?", AbnormalEntity.RoomID).Delete(&mysql.SlimeWarGameRoomSettings{}).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("게임 셋팅 삭제 실패: %v", err),
			Type: _errors.ErrDeleteGameRoomSettingFailed,
		}
	}
	return nil
}

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
