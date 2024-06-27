package usecase

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"main/features/auth/model/request"
	"main/utils"
	"main/utils/db/mysql"
	"net/http"
)

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

func getGoogleUserInfo(ctx context.Context, accessToken string) (string, error) {
	fmt.Println("accessToken: ", accessToken)
	token, err := utils.GoogleConfig.Exchange(ctx, accessToken)
	fmt.Println(token)
	if err != nil {
		return "", err
	}
	resp, err := http.Get("https://www.googleapis.com/oauth2/v1/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(content), nil
}
