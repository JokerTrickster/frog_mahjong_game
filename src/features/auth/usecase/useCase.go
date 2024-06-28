package usecase

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"io/ioutil"
	"main/features/auth/model/request"
	"main/utils"
	"main/utils/db/mysql"
	"net/http"
)

func CreateUserSQL(email string) *mysql.Users {
	return &mysql.Users{
		Name:     email,
		Email:    email,
		State:    "wait",
		Score:    30,
		RoomID:   1,
		Provider: "google",
	}
}

func CreateTokenDTO(uID uint, accessToken string, accessTknExpiredAt int64, refreshToken string, refreshTknExpiredAt int64) mysql.Tokens {
	return mysql.Tokens{
		UserID:           uID,
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		RefreshExpiredAt: refreshTknExpiredAt,
	}
}

func CreateSignupUser(req *request.ReqSignup) mysql.Users {
	return mysql.Users{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Score:    30,
		RoomID:   1,
		State:    "logout",
		Provider: "email",
	}
}

func VerifyAccessAndRefresh(req *request.ReqReissue) error {
	if err := utils.VerifyToken(req.AccessToken); err != nil {
		return err
	}
	if err := utils.VerifyToken(req.RefreshToken); err != nil {
		return err
	}
	return nil
}

func GenerateStateOauthCookie(ctx context.Context) string {
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)

	return state
}

func getGoogleUserInfo(ctx context.Context, accessToken string) ([]byte, error) {
	token, err := utils.GoogleConfig.Exchange(ctx, accessToken)
	if err != nil {
		return nil, err
	}
	resp, err := http.Get("https://www.googleapis.com/oauth2/v1/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return content, nil
}
