package usecase

import (
	"context"
	"main/features/auth/model/entity"
	_interface "main/features/auth/model/interface"
	"main/utils/db/mysql"

	"time"
)

type ValidatePasswordAuthUseCase struct {
	Repository     _interface.IValidatePasswordAuthRepository
	ContextTimeout time.Duration
}

func NewValidatePasswordAuthUseCase(repo _interface.IValidatePasswordAuthRepository, timeout time.Duration) _interface.IValidatePasswordAuthUseCase {
	return &ValidatePasswordAuthUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *ValidatePasswordAuthUseCase) ValidatePassword(c context.Context, e entity.ValidatePasswordAuthEntity) error {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()

	// 1. 코드 검증
	err := d.Repository.CheckAuthCode(ctx, e.Email, e.Code)
	if err != nil {
		return err
	}

	// 2. 비밀번호 변경
	user := mysql.Users{
		Email:    e.Email,
		Password: e.Password,
	}
	err = d.Repository.UpdatePassword(ctx, user)
	if err != nil {
		return err
	}
	
	// 3. 코드 삭제
	err = d.Repository.DeleteAuthCode(ctx, e.Email)
	if err != nil {
		return err
	}

	return nil
}
