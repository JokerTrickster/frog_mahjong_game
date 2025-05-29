package repository

import (
	"context"
	"fmt"
	"main/features/sequence/model/entity"
	_errors "main/features/sequence/model/errors"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func AbnormalDeleteAllCards(ctx context.Context, tx *gorm.DB, AbnormalEntity *entity.WSAbnormalEntity) *entity.ErrorInfo {
	err := tx.Model(&mysql.SequenceRoomCards{}).Where("room_id = ?", AbnormalEntity.RoomID).Delete(&mysql.SequenceRoomCards{}).Error
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
	err := tx.Model(&mysql.SequenceRoomMaps{}).Where("room_id = ?", AbnormalEntity.RoomID).Delete(&mysql.SequenceRoomMaps{}).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("맵 삭제 실패: %v", err),
			Type: _errors.ErrDeleteMapFailed,
		}
	}
	return nil
}

func AbnormalDeleteGameRoomSetting(ctx context.Context, tx *gorm.DB, AbnormalEntity *entity.WSAbnormalEntity) *entity.ErrorInfo {
	err := tx.Model(&mysql.SequenceGameRoomSettings{}).Where("room_id = ?", AbnormalEntity.RoomID).Delete(&mysql.SequenceGameRoomSettings{}).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("게임 셋팅 삭제 실패: %v", err),
			Type: _errors.ErrDeleteGameRoomSettingFailed,
		}
	}
	return nil
}

func AbnormalDeleteRoomUsers(ctx context.Context, tx *gorm.DB, AbnormalEntity *entity.WSAbnormalEntity) *entity.ErrorInfo {
	err := tx.Model(&mysql.SequenceUsers{}).Where("room_id = ?", AbnormalEntity.RoomID).Delete(&mysql.SequenceUsers{}).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("방 유저 삭제 실패: %v", err),
			Type: _errors.ErrDeleteRoomUserFailed,
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

func AbnormalUpdateRoomUsers(ctx context.Context, tx *gorm.DB, AbnormalEntity *entity.WSAbnormalEntity) *entity.ErrorInfo {
	err := tx.Model(&mysql.SequenceUsers{}).Where("room_id = ? and user_id = ?", AbnormalEntity.RoomID, AbnormalEntity.UserID).Update("player_state", "disconnected").Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("방 유저 상태 변경 실패: %v", err),
			Type: _errors.ErrUpdateUserStateFailed,
		}
	}
	return nil
}
