package usecase

import (
	"context"
	"fmt"
	_interface "main/features/auth/model/interface"
	"time"
)

type GoogleOauthAuthUseCase struct {
	Repository             _interface.IGoogleOauthAuthRepository
	ContextTimeGoogleOauth time.Duration
}

func NewGoogleOauthAuthUseCase(repo _interface.IGoogleOauthAuthRepository, timeGoogleOauth time.Duration) _interface.IGoogleOauthAuthUseCase {
	return &GoogleOauthAuthUseCase{Repository: repo, ContextTimeGoogleOauth: timeGoogleOauth}
}

func (d *GoogleOauthAuthUseCase) GoogleOauth(c context.Context) error {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeGoogleOauth)
	defer cancel()
	fmt.Println(ctx)
	return nil
}
