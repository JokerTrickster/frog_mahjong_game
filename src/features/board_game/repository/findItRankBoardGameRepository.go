package repository

import (
	"context"
	"fmt"
	"main/features/board_game/model/entity"
	_interface "main/features/board_game/model/interface"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewFindItRankBoardGameRepository(gormDB *gorm.DB) _interface.IFindItRankBoardGameRepository {
	return &FindItRankBoardGameRepository{GormDB: gormDB}
}

func (d *FindItRankBoardGameRepository) FindTop3UserCorrect(ctx context.Context) ([]*entity.FindItRankEntity, error) {
	var results []*entity.FindItRankEntity

	// GORM의 SubQuery를 사용하여 삭제되지 않은 유저의 id만 선택합니다.
	subQuery := d.GormDB.Model(&mysql.GameUsers{}).Select("id").Where("deleted_at IS NULL")

	err := d.GormDB.WithContext(ctx).
		Table("find_it_user_correct_positions").
		Select("user_id, COUNT(*) as correct_count").
		Where("user_id IN (?)", subQuery).
		Group("user_id").
		Order("correct_count DESC").
		Limit(3).
		Scan(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}
func (d *FindItRankBoardGameRepository) FindOneUser(ctx context.Context, userId int) (*mysql.GameUsers, error) {
	var result *mysql.GameUsers
	err := d.GormDB.WithContext(ctx).
		Table("game_users").
		Where("id = ?", userId).
		First(&result).Error
	fmt.Println(result)
	if err != nil {
		return nil, err
	}
	return result, nil

}
