package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"main/features/auth/model/entity"
	_interface "main/features/auth/model/interface"
	"main/features/auth/model/response"
	"main/utils"
	"main/utils/db/mysql"
	"time"
)

type GoogleOauthCallbackAuthUseCase struct {
	Repository     _interface.IGoogleOauthCallbackAuthRepository
	ContextTimeout time.Duration
}

func NewGoogleOauthCallbackAuthUseCase(repo _interface.IGoogleOauthCallbackAuthRepository, timeout time.Duration) _interface.IGoogleOauthCallbackAuthUseCase {
	return &GoogleOauthCallbackAuthUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *GoogleOauthCallbackAuthUseCase) GoogleOauthCallback(c context.Context, code string) (response.GoogleOauthCallbackRes, error) {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()

	fmt.Println("code", code)
	data, err := getGoogleUserInfo(ctx, code)
	if err != nil {
		return response.GoogleOauthCallbackRes{}, err
	}
	var googleUser entity.GoogleUser
	fmt.Println("?")
	// JSON 파싱
	if err := json.Unmarshal(data, &googleUser); err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	sqlEntity := &entity.GoogleOauthCallbackSQLQuery{}
	var user *mysql.Users
	//유저 체크 후 있으면 로그인 처리
	user, err = d.Repository.FindOneAndUpdateUser(ctx, sqlEntity)
	if err != nil {
		return response.GoogleOauthCallbackRes{}, err
	}
	//유저가 없으면 새로 생성한다.
	if user == nil {
		//유저 생성
		userCreateSQLEntity := CreateUserSQL(googleUser.Email)
		user, err = d.Repository.CreateUser(ctx, userCreateSQLEntity)
		if err != nil {
			return response.GoogleOauthCallbackRes{}, err
		}
	}

	//토큰 생성
	// token create
	accessToken, _, refreshToken, refreshTknExpiredAt, err := utils.GenerateToken(user.Email, user.ID)
	if err != nil {
		return response.GoogleOauthCallbackRes{}, err
	}

	// 기존 토큰 제거
	err = d.Repository.DeleteToken(ctx, user.ID)
	if err != nil {
		return response.GoogleOauthCallbackRes{}, err
	}
	// token db save
	err = d.Repository.SaveToken(ctx, user.ID, accessToken, refreshToken, refreshTknExpiredAt)
	if err != nil {
		return response.GoogleOauthCallbackRes{}, err
	}

	//response create
	res := response.GoogleOauthCallbackRes{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return res, nil
}
