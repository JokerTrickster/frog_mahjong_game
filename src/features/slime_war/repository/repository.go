package repository

import (
	"context"
	"fmt"
	"main/features/slime_war/model/entity"
	_errors "main/features/slime_war/model/errors"
	"main/utils/db/mysql"
	_redis "main/utils/db/redis"

	"gorm.io/gorm"
)

type JoinEventWebsocketRepository struct {
	GormDB *gorm.DB
}

func PreloadUsers(ctx context.Context, tx *gorm.DB, roomID uint) ([]entity.PreloadUsers, *entity.ErrorInfo) {
	var preloadUsers []entity.PreloadUsers

	if err := tx.Table("game_room_users"). // ✅ 올바른 테이블 명시
						Preload("User").        // ✅ 유저 정보 (game_users)
						Preload("Room").        // ✅ 게임 방 정보 (game_rooms)
						Preload("RoomSetting"). // ✅ 방 설정 정보 (slime_war_room_settings)
						Preload("UserCorrectPositions", func(db *gorm.DB) *gorm.DB {
			return db.Where("room_id = ?", roomID) // ✅ 유저가 맞춘 정답 정보 가져오기
		}).
		Preload("RoundImages", func(db *gorm.DB) *gorm.DB {
			return db.Where("room_id = ?", roomID) // ✅ 현재 방의 라운드 이미지만 가져오기
		}).
		Where("room_id = ?", roomID).
		Find(&preloadUsers).Error; err != nil {
		return nil, &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("room_users 조회 실패: %v", err.Error()),
			Type: _errors.ErrRoomUsersNotFound,
		}
	}

	return preloadUsers, nil
}

func FindAllCorrectPositions(ctx context.Context, correctIDList []int) ([]mysql.FindItImageCorrectPositions, *entity.ErrorInfo) {
	var correctPositions []mysql.FindItImageCorrectPositions
	if err := mysql.GormMysqlDB.WithContext(ctx).Where("id in ?", correctIDList).Find(&correctPositions).Error; err != nil {
		return nil, &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("FindAllCorrectPositions: %v", err.Error()),
			Type: _errors.ErrInternalServer,
		}
	}
	return correctPositions, nil
}

func FindOneRoundImage(c context.Context, imageID int) (*mysql.FindItImages, *entity.ErrorInfo) {
	var image mysql.FindItImages
	if err := mysql.GormMysqlDB.WithContext(c).Where("id = ?", imageID).First(&image).Error; err != nil {
		return nil, &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("FindOneRoundImage: %v", err.Error()),
			Type: _errors.ErrInternalServer,
		}
	}
	return &image, nil
}

func FindAllOpenCards(c context.Context, roomID int) ([]int, *entity.ErrorInfo) {
	var cards []int
	if err := mysql.GormMysqlDB.WithContext(c).Model(&mysql.UserBirdCards{}).Where("room_id = ? and state = ?", roomID, "opened").Pluck("card_id", &cards).Error; err != nil {
		return nil, &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("FindAllOpenCards: %v", err.Error()),
			Type: _errors.ErrInternalServer,
		}
	}
	return cards, nil
}

func ReconnectedUpdateRoomUser(c context.Context, roomID uint, userID uint) *entity.ErrorInfo {
	err := mysql.GormMysqlDB.Model(&mysql.RoomUsers{}).Where("room_id = ? and user_id = ?", roomID, userID).Update("player_state", "play").Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("ReconnectedUpdateRoomUser: %v", err.Error()),
			Type: _errors.ErrInternalServer,
		}
	}
	return nil
}

func RedisSessionDelete(ctx context.Context, sessionID string) *entity.ErrorInfo {
	redisKey := fmt.Sprintf("abnormal_session:%s", sessionID)
	// Redis에서 키 삭제
	err := _redis.Client.Del(ctx, redisKey).Err()
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("세션 삭제 실패: %v", err.Error()),
			Type: _errors.ErrInternalServer,
		}
	}

	return nil
}

func DeleteAllUserBirdCards(c context.Context, tx *gorm.DB, userID uint) *entity.ErrorInfo {
	err := tx.Where("user_id = ?", userID).Delete(&mysql.UserBirdCards{}).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("DeleteAllUserCards: %v", err.Error()),
			Type: _errors.ErrInternalServer,
		}
	}
	return nil
}

func DeleteAllRoomUsers(c context.Context, tx *gorm.DB, userID uint) *entity.ErrorInfo {
	err := tx.Where("user_id = ?", userID).Delete(&mysql.GameRoomUsers{}).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("DeleteAllRoomUsers: %v", err.Error()),
			Type: _errors.ErrInternalServer,
		}
	}
	return nil
}
func DeleteAllRooms(c context.Context, tx *gorm.DB, userID uint) *entity.ErrorInfo {
	err := tx.Where("owner_id = ?", userID).Delete(&mysql.GameRooms{}).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("DeleteAllRooms: %v", err.Error()),
			Type: _errors.ErrInternalServer,
		}
	}
	return nil
}

func DeleteAllUserMissions(c context.Context, tx *gorm.DB, userID uint) *entity.ErrorInfo {
	err := tx.Where("user_id = ?", userID).Delete(&mysql.UserMissions{}).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("DeleteAllUserMissions: %v", err.Error()),
			Type: _errors.ErrInternalServer,
		}
	}
	return nil
}

func DeleteAllUserItems(c context.Context, tx *gorm.DB, userID uint) *entity.ErrorInfo {
	err := tx.Where("user_id = ?", userID).Delete(&mysql.UserItems{}).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("DeleteAllUserItems: %v", err.Error()),
			Type: _errors.ErrInternalServer,
		}
	}
	return nil
}

func DeleteAllGameRooms(c context.Context, tx *gorm.DB, userID uint) *entity.ErrorInfo {
	err := tx.Where("owner_id = ?", userID).Delete(&mysql.GameRooms{}).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("DeleteAllGameRooms: %v", err.Error()),
			Type: _errors.ErrInternalServer,
		}
	}
	return nil
}

func DeleteAllGameRoomUsers(c context.Context, tx *gorm.DB, userID uint) *entity.ErrorInfo {
	err := tx.Where("user_id = ?", userID).Delete(&mysql.GameRoomUsers{}).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal,
			Msg:  fmt.Sprintf("DeleteAllGameRooms: %v", err.Error()),
			Type: _errors.ErrInternalServer,
		}
	}
	return nil
}
