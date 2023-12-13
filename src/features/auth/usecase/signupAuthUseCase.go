package usecase

import (
	"context"
	_interface "main/features/auth/model/interface"
	"main/features/auth/model/request"
	"time"
)

type SignupAuthUseCase struct {
	Repository     _interface.ISignupAuthRepository
	ContextTimeout time.Duration
}

func NewSignupAuthUseCase(repo _interface.ISignupAuthRepository, timeout time.Duration) _interface.ISignupAuthUseCase {
	return &SignupAuthUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *SignupAuthUseCase) Signup(c context.Context, req *request.ReqSignup) error {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()

	//1. 해당 이름이 존재하는지 체크
	err := d.Repository.FindOneUserAuth(ctx, req.Name)
	if err != nil {
		return err
	}

	//2. 유저 DTO를 만든다.
	gUserDTO, err := CreateSignupUserDTO(req.Name, req.Email)
	if err != nil {
		return err
	}
	//3. 유저 정보를 저장한다.
	userID, err := d.Repository.InsertOneUserDTO(ctx, gUserDTO)
	if err != nil {
		return err
	}

	//4. 유저 인증 DTO를 만든다.
	gUserAuthDTO, err := CreateSignupUserAuthDTO(userID, req)
	if err != nil {
		return err
	}

	//5. 유저 인증 정보를 저장한다.
	err = d.Repository.InsertOneUserAuthDTO(ctx, gUserAuthDTO)
	if err != nil {
		return err
	}

	return nil
}
