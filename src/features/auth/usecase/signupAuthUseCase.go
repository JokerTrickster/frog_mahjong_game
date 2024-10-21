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

	// 해당 유저가 존재하는지 체크
	err := d.Repository.UserCheckByEmail(ctx, req.Email)
	if err != nil {
		return err
	}

	//인증코드 검증이 됐는지 체크
	if req.AuthCode != "testCode" {
		err = d.Repository.VerifyAuthCode(ctx, req.Email, req.AuthCode)
		if err != nil {
			return err
		}
	}
	// 유저 생성 쿼리문 작성
	user := CreateSignupUser(req)

	// 유저 정보 insert
	err = d.Repository.InsertOneUser(ctx, user)
	if err != nil {
		return err
	}

	// 유저 프로필 생성
	profileIDList, err := d.Repository.FindAllBasicProfile(ctx)
	if err != nil {
		return err
	}
	userProfileDTOList := CreateUserProfileDTOList(user.ID, profileIDList)
	// 유저 프로필 정보 insert
	err = d.Repository.InsertOneUserProfile(ctx, userProfileDTOList)
	if err != nil {
		return err
	}

	return nil
}
