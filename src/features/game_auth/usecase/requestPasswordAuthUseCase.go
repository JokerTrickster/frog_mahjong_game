package usecase

import (
	"context"
	"main/features/game_auth/model/entity"
	_interface "main/features/game_auth/model/interface"
	"main/utils/aws"
	"main/utils/db/mysql"

	"time"
)

type RequestPasswordAuthUseCase struct {
	Repository     _interface.IRequestPasswordAuthRepository
	ContextTimeout time.Duration
}

func NewRequestPasswordAuthUseCase(repo _interface.IRequestPasswordAuthRepository, timeout time.Duration) _interface.IRequestPasswordAuthUseCase {
	return &RequestPasswordAuthUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *RequestPasswordAuthUseCase) RequestPassword(c context.Context, e entity.RequestPasswordAuthEntity) (string, error) {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()

	// 1. 이메일 사용자 조회
	err := d.Repository.FindOneUserByEmail(ctx, e.Email)
	if err != nil {
		return "", err
	}

	// 2. 비밀번호 재설정 토큰 생성
	authCode := GeneratePasswordAuthCode()

	// 3. 토큰 저장
	userAuthDTO := mysql.UserAuths{
		Email:    e.Email,
		AuthCode: authCode,
		Type:     "password",
		Project:  "board_game",
		IsActive: false,
	}
	err = d.Repository.InsertAuthCode(ctx, userAuthDTO)
	if err != nil {
		return "", err
	}
	//TODO 4. 추후 이메일 전송
	aws.EmailSendPassword(e.Email, authCode)

	return authCode, nil
}
