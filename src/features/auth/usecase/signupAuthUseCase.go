package usecase

import (
	"context"
	"fmt"
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

func (s *SignupAuthUseCase) Signup(c context.Context, req *request.ReqSignup) error {
	ctx, cancel := context.WithTimeout(c, s.ContextTimeout)
	defer cancel()
	fmt.Println(ctx)
	//1. 해당 이름이 존재하는지 체크

	//2. 회원가입 정보 DTO를 만든다.

	//3. 유저 정보를 저장한다.

	return nil
}
