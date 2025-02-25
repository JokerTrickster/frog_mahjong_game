package usecase

import (
	"context"
	"encoding/json"
	"log"
	"main/features/game_auth/model/entity"
	_interface "main/features/game_auth/model/interface"
	"main/features/game_auth/model/response"
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

func (d *GoogleOauthCallbackAuthUseCase) GoogleOauthCallback(c context.Context, code string) (response.ResGameGoogleOauthCallback, error) {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()
	data, err := getAppGoogleUserInfo(ctx, code)
	if err != nil {
		return response.ResGameGoogleOauthCallback{}, err
	}
	var googleUser entity.GoogleUser
	// JSON 파싱
	if err := json.Unmarshal(data, &googleUser); err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	sqlEntity := &entity.GoogleOauthCallbackSQLQuery{
		Email: googleUser.Email,
	}
	var user *mysql.GameUsers
	//유저 체크 후 있으면 로그인 처리
	user, err = d.Repository.FindOneAndUpdateUser(ctx, sqlEntity)
	if err != nil {
		return response.ResGameGoogleOauthCallback{}, err
	}
	//유저가 없으면 새로 생성한다.
	if user == nil {
		//유저 생성
		userCreateSQLEntity := CreateUserSQL(googleUser.Email)
		user, err = d.Repository.CreateUser(ctx, userCreateSQLEntity)
		if err != nil {
			return response.ResGameGoogleOauthCallback{}, err
		}
		// 기본 프로필 정보를 가져온다
		profileIDList, err := d.Repository.FindAllBasicProfile(ctx)
		if err != nil {
			return response.ResGameGoogleOauthCallback{}, err
		}
		userProfileDTOList := CreateUserProfileDTOList(user.ID, profileIDList)
		// 유저 프로필 정보 insert
		err = d.Repository.InsertOneUserProfile(ctx, userProfileDTOList)
		if err != nil {
			return response.ResGameGoogleOauthCallback{}, err
		}
	}

	// 기존 토큰이 있는지 체크
	prevTokens, err := d.Repository.CheckToken(ctx, user.ID)
	if err != nil {
		return response.ResGameGoogleOauthCallback{}, err
	}
	res := response.ResGameGoogleOauthCallback{
		IsDuplicateLogin: false,
	}
	if prevTokens != nil {
		res.IsDuplicateLogin = true
	}

	//토큰 생성
	// token create
	accessToken, _, refreshToken, refreshTknExpiredAt, err := utils.GenerateToken(user.Email, user.ID)
	if err != nil {
		return response.ResGameGoogleOauthCallback{}, err
	}

	// 기존 토큰 제거
	err = d.Repository.DeleteToken(ctx, user.ID)
	if err != nil {
		return response.ResGameGoogleOauthCallback{}, err
	}
	// token db save
	err = d.Repository.SaveToken(ctx, user.ID, accessToken, refreshToken, refreshTknExpiredAt)
	if err != nil {
		return response.ResGameGoogleOauthCallback{}, err
	}

	//response create
	res.AccessToken = accessToken
	res.RefreshToken = refreshToken
	res.UserID = user.ID

	return res, nil
}
