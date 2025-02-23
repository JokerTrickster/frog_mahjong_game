package usecase

import (
	"context"
	"main/features/game_auth/model/entity"
	_interface "main/features/game_auth/model/interface"

	"time"
)

type ValidateSignupAuthUseCase struct {
	Repository     _interface.IValidateSignupAuthRepository
	ContextTimeout time.Duration
}

func NewValidateSignupAuthUseCase(repo _interface.IValidateSignupAuthRepository, timeout time.Duration) _interface.IValidateSignupAuthUseCase {
	return &ValidateSignupAuthUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *ValidateSignupAuthUseCase) ValidateSignup(c context.Context, e entity.ValidateSignupAuthEntity) error {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()

	// 1. 코드 검증
	userAuthDTO := CreateUserAuth(e.Email, e.Code)
	err := d.Repository.CheckAuthCode(ctx, userAuthDTO)
	if err != nil {
		return err
	}
	// 2. 검증 활성화 한다. 
	err = d.Repository.UpdateAuthCode(ctx, userAuthDTO)
	if err != nil {
		return err
	}

	return nil
}
