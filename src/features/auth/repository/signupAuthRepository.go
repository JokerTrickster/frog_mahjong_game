package repository

import (
	"context"
	"fmt"
	_interface "main/features/auth/model/interface"
	"main/utils/db/mysql"

	"gorm.io/gorm"
)

func NewSignupAuthRepository(gormDB *gorm.DB) _interface.ISignupAuthRepository {
	return &SignupAuthRepository{GormDB: gormDB}
}

func (g *SignupAuthRepository) FindOneUserAuth(ctx context.Context, name string) error {
	var userAuthDTO mysql.GormUserAuthDTO
	result := g.GormDB.WithContext(ctx).Where("name = ?", name).Find(&userAuthDTO)
	fmt.Println(result.RowsAffected)
	if result.RowsAffected == 0 {
		return nil
	} else {
		return fmt.Errorf("%s name is already existed", name)
	}
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (g *SignupAuthRepository) InsertOneUserDTO(ctx context.Context, userDTO mysql.GormUserDTO) (string, error) {
	result := g.GormDB.WithContext(ctx).Create(&userDTO)
	if result.RowsAffected == 0 {
		return "", fmt.Errorf("failed userDTO insert one")
	}
	if result.Error != nil {
		return "", result.Error
	}
	return userDTO.GormModel.ID, nil
}
func (g *SignupAuthRepository) InsertOneUserAuthDTO(ctx context.Context, userAuthDTO mysql.GormUserAuthDTO) error {
	result := g.GormDB.WithContext(ctx).Create(&userAuthDTO)
	if result.RowsAffected == 0 {
		return fmt.Errorf("failed userAuthDTO insert one")
	}
	if result.Error != nil {
		return result.Error
	}
	return nil
}
