package repository

import (
	"context"
	"fmt"
	"main/features/v2ws/model/entity"
	"main/utils/db/mysql"

	_errors "main/features/v2ws/model/errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func ImportSingleCardFindAllRoomUsers(ctx context.Context, tx *gorm.DB, roomID uint) ([]entity.RoomUsers, error) {
	var roomUsers []entity.RoomUsers
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Preload("User").Preload("UserItems").Preload("Room").Preload("RoomMission").Preload("Cards", func(db *gorm.DB) *gorm.DB {
		return db.Where("room_id = ?", roomID).Order("updated_at ASC")
	}).Preload("UserMissions", func(db *gorm.DB) *gorm.DB {
		return db.Where("room_id = ?", roomID)
	}).Where("room_id = ?", roomID).Find(&roomUsers).Error; err != nil {
		return nil, fmt.Errorf("room_users 조회 에러: %v", err.Error())
	}
	return roomUsers, nil
}
func ImportSingleCardUpdateAllCardState(c context.Context, tx *gorm.DB, roomID uint) error {
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Model(&mysql.UserBirdCards{}).Where("room_id = ? and state = ?", roomID, "picked").Update("state", "owned").Error
	if err != nil {
		return fmt.Errorf("카드 상태 업데이트 실패 %v", err.Error())
	}
	return nil
}

// 카드정보가 존재하면 상태를 업데이트하고 없으면 카드를 생성한다.
func ImportSingleCardCreateCard(c context.Context, tx *gorm.DB, userBirdCardDTO *mysql.UserBirdCards) error {
	// 카드 정보가 존재하는지 확인
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Model(&mysql.UserBirdCards{}).Where("room_id = ? and card_id = ?", userBirdCardDTO.RoomID, userBirdCardDTO.CardID).Updates(userBirdCardDTO).Error
	if err != nil {
		return fmt.Errorf("카드 조회 및 업데이트 실패 %v", err.Error())
	}

	return nil
}

func ImportSingleCardUpdateRoomUserCardCount(c context.Context, tx *gorm.DB, entity *entity.WSImportSingleCardEntity) error {
	// 유저id로 room_users에서 찾아서 card_count를 더한 후 업데이트 한다.
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Model(&mysql.RoomUsers{}).Where("room_id = ? AND user_id = ?", entity.RoomID, entity.UserID).Update("owned_card_count", gorm.Expr("owned_card_count + 1")).Error
	if err != nil {
		return fmt.Errorf("방 유저 카드 카운트 업데이트 실패 %v", err.Error())
	}
	return nil
}

func ImportSingleCardFindAllCard(c context.Context, tx *gorm.DB, roomID uint, userID uint) ([]*mysql.Cards, error) {
	cards := make([]*mysql.Cards, 0)
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Model(&mysql.Cards{}).Where("room_id = ? and user_id = ?", roomID, userID).Find(&cards).Error
	if err != nil {
		return nil, fmt.Errorf("카드를 찾을 수 없습니다. %v", err.Error())
	}
	return cards, nil
}

func ImportSingleCardUpdateOpenCards(ctx context.Context, tx *gorm.DB, roomID uint) error {
	// 오픈 카드가 비어 있다면 새로운 카드를 오픈한다.
	// 현재 오픈 카드가 몇개 있는지 카운트 한다.
	var count int64
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Model(&mysql.UserBirdCards{}).Where("room_id = ? and state = ?", roomID, "opened").Count(&count).Error
	if err != nil {
		return fmt.Errorf("오픈 카드 카운트 실패 %v", err.Error())
	}
	if count != 3 {
		openCardCount := 3 - count

		// 상태가 'none'인 카드 중에서 랜덤으로 openCardCount 수만큼 카드 ID를 가져온다.
		var cardIDs []int
		err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).WithContext(ctx).
			Model(&mysql.UserBirdCards{}).
			Where("room_id = ? AND state = ?", roomID, "none").
			Order("RAND()").
			Limit(int(openCardCount)).
			Pluck("card_id", &cardIDs).Error
		if err != nil {
			return fmt.Errorf("카드 조회 실패: %v", err.Error())
		}

		// 선택된 카드들의 상태를 opened로 변경한다.
		if len(cardIDs) > 0 {
			err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).WithContext(ctx).
				Model(&mysql.UserBirdCards{}).
				Where("room_id = ? AND card_id IN ?", roomID, cardIDs).
				Update("state", "opened").Error
			if err != nil {
				return fmt.Errorf("카드 상태 업데이트 실패: %v", err.Error())
			}
		}
	}

	return nil
}

func ImportSingleCardOwnerCardCount(ctx context.Context, roomID uint, userID uint) (int, *entity.ErrorInfo) {
	var roomUsers mysql.RoomUsers
	err := mysql.GormMysqlDB.Model(&mysql.RoomUsers{}).Where("room_id = ? and user_id = ?", roomID, userID).Find(&roomUsers).Error
	if err != nil {
		return 0, &entity.ErrorInfo{Code: _errors.ErrCodeInternal, Msg: err.Error(), Type: _errors.ErrInternalServer}
	}
	if roomUsers.OwnedCardCount > 3 {
		return 0, &entity.ErrorInfo{Code: _errors.ErrCodeBadRequest, Msg: "카드를 4장 이상 가지"}
	}
	return roomUsers.OwnedCardCount, nil
}

func ImportSingleCardFindOneCard(ctx context.Context, roomID uint, cardID uint) *entity.ErrorInfo {
	var card mysql.UserBirdCards
	result := mysql.GormMysqlDB.Model(&mysql.UserBirdCards{}).Where("room_id = ? AND card_id = ? AND (state = ? OR state = ?)", roomID, cardID, "opened", "none").Find(&card)
	if result.Error != nil {
		return &entity.ErrorInfo{Code: 500, Msg: result.Error.Error(), Type: _errors.ErrInternalServer}
	}
	if result.RowsAffected == 0 {
		return &entity.ErrorInfo{Code: _errors.ErrCodeBadRequest, Msg: "이미 카드를 가지고 갔습니다.", Type: _errors.ErrNotFoundCard}
	}
	return nil
}
