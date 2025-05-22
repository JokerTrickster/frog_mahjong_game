package repository

import (
	"context"
	"main/features/frog/model/entity"
	"main/utils/db/mysql"

	_errors "main/features/frog/model/errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// StartFindAllRoomUsers retrieves all room users with necessary preloads
func StartFindAllRoomUsers(ctx context.Context, tx *gorm.DB, roomID uint) ([]entity.RoomUsers, *entity.ErrorInfo) {
	var roomUsers []entity.RoomUsers
	if err := tx.Table("frog_room_users").
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("room_id = ?", roomID).
		Preload("User").
		Preload("Room").
		Preload("Cards", func(db *gorm.DB) *gorm.DB {
			return db.Where("room_id = ?", roomID).Order("updated_at ASC")
		}).
		Find(&roomUsers).Error; err != nil {
		return nil, &entity.ErrorInfo{
			Code: _errors.ErrCodeNotFound, // 404
			Msg:  "방 유저를 찾을 수 없습니다",
			Type: _errors.ErrRoomUsersNotFound,
		}
	}
	return roomUsers, nil
}

// StartDeleteCards deletes all cards for a specific room
func StartDeleteCards(ctx context.Context, tx *gorm.DB, roomID uint) *entity.ErrorInfo {
	err := tx.WithContext(ctx).Where("room_id = ?", roomID).Delete(&mysql.FrogUserCards{})
	if err.Error != nil {
		if err.Error == gorm.ErrRecordNotFound {
			return nil // No cards to delete, not an error
		}
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal, // 500
			Msg:  "카드 삭제 실패",
			Type: _errors.ErrDeleteFailed,
		}
	}
	return nil
}

// StartCheckOwner checks if the user is the room owner
func StartCheckOwner(ctx context.Context, tx *gorm.DB, uID uint, roomID uint) *entity.ErrorInfo {
	room := mysql.GameRooms{}
	err := tx.WithContext(ctx).Where("id = ?", roomID).First(&room).Error
	if err != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeNotFound, // 404
			Msg:  "방 정보를 찾을 수 없습니다",
			Type: _errors.ErrRoomNotFound,
		}
	}
	if room.OwnerID != int(uID) {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeForbidden, // 403
			Msg:  "방장만 게임을 시작할 수 있습니다",
			Type: _errors.ErrUnauthorizedAction,
		}
	}
	return nil
}

// StartFindRoomUsers retrieves room users for a specific room
func StartFindRoomUsers(ctx context.Context, tx *gorm.DB, roomID uint) ([]mysql.FrogRoomUsers, *entity.ErrorInfo) {
	roomUsers := make([]mysql.FrogRoomUsers, 0)
	err := tx.WithContext(ctx).Where("room_id = ?", roomID).Find(&roomUsers).Error
	if err != nil {
		return nil, &entity.ErrorInfo{
			Code: _errors.ErrCodeNotFound, // 404
			Msg:  "방 유저 정보를 찾을 수 없습니다",
			Type: _errors.ErrRoomUsersNotFound,
		}
	}
	return roomUsers, nil
}

// StartUpdateRoomUser updates room users' data
func StartUpdateRoomUser(ctx context.Context, tx *gorm.DB, updateRoomUsers []mysql.FrogRoomUsers) *entity.ErrorInfo {
	for _, user := range updateRoomUsers {
		err := tx.WithContext(ctx).Model(&mysql.FrogRoomUsers{}).
			Where("room_id = ? AND user_id = ?", user.RoomID, user.UserID).
			Updates(user)
		if err.Error != nil {
			return &entity.ErrorInfo{
				Code: _errors.ErrCodeInternal, // 500
				Msg:  "방 유저 정보 업데이트 실패",
				Type: _errors.ErrUpdateUserStateFailed,
			}
		}
	}
	return nil
}

// StartUpdateRoom updates the room's state
func StartUpdateRoom(ctx context.Context, tx *gorm.DB, roomID uint, state string) *entity.ErrorInfo {
	err := tx.WithContext(ctx).Model(&mysql.GameRooms{}).Where("id = ? and state = ?", roomID, "wait").Update("state", state)
	if err.Error != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal, // 500
			Msg:  "방 상태 업데이트 실패",
			Type: _errors.ErrUpdateFailed,
		}
	}
	return nil
}

// StartCreateCards creates new card data
func StartCreateCards(ctx context.Context, tx *gorm.DB, cards []mysql.FrogUserCards) *entity.ErrorInfo {
	err := tx.WithContext(ctx).Create(&cards)
	if err.Error != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal, // 500
			Msg:  "카드 정보 생성 실패",
			Type: _errors.ErrCreateFailed,
		}
	}
	return nil
}

// StartFindCards retrieves all cards
func StartFindCards(ctx context.Context, tx *gorm.DB) ([]mysql.FrogCards, *entity.ErrorInfo) {
	var cards []mysql.FrogCards
	if err := tx.Find(&cards).Error; err != nil {
		return nil, &entity.ErrorInfo{
			Code: _errors.ErrCodeNotFound, // 404
			Msg:  "카드 조회 에러",
			Type: _errors.ErrNotFoundCard,
		}
	}
	return cards, nil
}

func StartCreateFrogGameRoomSettings(ctx context.Context, tx *gorm.DB, frogGameRoomSettingsDTO *mysql.FrogGameRoomSettings) *entity.ErrorInfo {
	err := tx.WithContext(ctx).Create(frogGameRoomSettingsDTO)
	if err.Error != nil {
		return &entity.ErrorInfo{
			Code: _errors.ErrCodeInternal, // 500
			Msg:  "게임 룸 설정 생성 실패",
			Type: _errors.ErrCreateFailed,
		}
	}
	return nil
}
