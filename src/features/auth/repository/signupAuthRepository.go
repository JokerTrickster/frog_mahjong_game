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
func (g *SignupAuthRepository) UserCheckByEmail(ctx context.Context, email string) error {
	var user mysql.Users
	result := g.GormDB.WithContext(ctx).Where("email = ?", email).First(&user)
	if result.RowsAffected == 0 {
		return nil
	} else {
		return fmt.Errorf("%s email is already existed", email)
	}
}
func (g *SignupAuthRepository) InsertOneUser(ctx context.Context, user mysql.Users) error {
	result := g.GormDB.WithContext(ctx).Create(&user)
	if result.RowsAffected == 0 {
		return fmt.Errorf("failed user insert one")
	}
	if result.Error != nil {
		return result.Error
	}
	return nil
}
