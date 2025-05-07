package repository

import (
	"context"
	"fmt"
	"main/features/board_game/model/entity"
	_interface "main/features/board_game/model/interface"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewSequenceRankBoardGameRepository(gormDB *gorm.DB) _interface.ISequenceRankBoardGameRepository {
	return &SequenceRankBoardGameRepository{GormDB: gormDB}
}

func (d *SequenceRankBoardGameRepository) FindTop3User(ctx context.Context) ([]*entity.SequenceRankEntity, error) {
	var results []*entity.SequenceRankEntity

	err := d.GormDB.WithContext(ctx).
		Table("game_results").
		Select("user_id, COUNT(*) as score").
		Where("result = 1 and game_type = ?", mysql.SEQUENCE).
		Group("user_id").
		Order("score DESC").
		Limit(3).
		Scan(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}
func (d *SequenceRankBoardGameRepository) FindOneUser(ctx context.Context, userId int) (*mysql.GameUsers, error) {
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
