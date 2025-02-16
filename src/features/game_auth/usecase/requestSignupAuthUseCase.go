package usecase

import (
	"context"
	"main/features/game_auth/model/entity"
	_interface "main/features/game_auth/model/interface"
	"main/utils/aws"
	"main/utils/db/mysql"

	"time"
)

type RequestSignupAuthUseCase struct {
	Repository     _interface.IRequestSignupAuthRepository
	ContextTimeout time.Duration
}

func NewRequestSignupAuthUseCase(repo _interface.IRequestSignupAuthRepository, timeout time.Duration) _interface.IRequestSignupAuthUseCase {
	return &RequestSignupAuthUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *RequestSignupAuthUseCase) RequestSignup(c context.Context, e entity.RequestSignupAuthEntity) (string, error) {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()

	// 1. 이메일 사용자 조회
	err := d.Repository.FindOneUserByEmail(ctx, e.Email)
	if err != nil {
		return "", err
	}

	// 2. 인증 코드 재설정 토큰 생성
	authCode := GeneratePasswordAuthCode()

	// 3. 토큰 저장
	userAuthDTO := mysql.UserAuths{
		Email:    e.Email,
		AuthCode: authCode,
		Type:     "signup",
	}
	// 4. 기존 이메일로 있는 인증 코드가 있다면 삭제한다.
	err = d.Repository.DeleteAuthCodeByEmail(ctx, e.Email)
	if err != nil {
		return "", err
	}

	// 5. 새로운 인증 코드를 삽입한다.
	err = d.Repository.InsertAuthCode(ctx, userAuthDTO)
	if err != nil {
		return "", err
	}
	//TODO 4. 추후 이메일 전송
	aws.EmailSendSignup(e.Email, authCode)

	return authCode, nil
}
