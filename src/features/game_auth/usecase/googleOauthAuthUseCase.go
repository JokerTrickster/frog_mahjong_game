package usecase

import (
	"context"
	_interface "main/features/game_auth/model/interface"
	"time"
)

type GoogleOauthAuthUseCase struct {
	Repository     _interface.IGoogleOauthAuthRepository
	ContextTimeout time.Duration
}

func NewGoogleOauthAuthUseCase(repo _interface.IGoogleOauthAuthRepository, timeout time.Duration) _interface.IGoogleOauthAuthUseCase {
	return &GoogleOauthAuthUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *GoogleOauthAuthUseCase) GoogleOauth(c context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()
	state := GenerateStateOauthCookie(ctx)
	return state, nil
}
