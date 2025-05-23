package usecase

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	_errors "main/features/game/model/errors"
	"main/features/game_auth/model/request"
	"main/utils"
	"main/utils/db/mysql"
	"net/http"
	"time"

	"golang.org/x/exp/rand"
)

const charset = "0123456789"

// 랜덤 값 생성 함수
func GeneratePasswordAuthCode() string {
	seed := rand.NewSource(uint64(time.Now().UnixNano()))
	r := rand.New(seed)
	b := make([]byte, 4)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b)
}

func CreateUserAuth(email string, authCode string) *mysql.UserAuths {
	return &mysql.UserAuths{
		Email:    email,
		AuthCode: authCode,
		Type:     "signup",
		Project:  "board_game",
		IsActive: false,
	}
}

func CreateUserSQL(email string) *mysql.GameUsers {
	return &mysql.GameUsers{
		Name:      "보린이",
		Email:     email,
		State:     "wait",
		Coin:      1000,
		RoomID:    1,
		Provider:  "google",
		ProfileID: 1,
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

func CreateSignupUser(req *request.ReqGameSignup) *mysql.GameUsers {
	return &mysql.GameUsers{
		Name:      req.Name,
		Email:     req.Email,
		Password:  req.Password,
		Coin:      1000,
		RoomID:    1,
		State:     "logout",
		Provider:  "email",
		ProfileID: 1,
	}
}

func VerifyAccessAndRefresh(req *request.ReqGameReissue) error {
	// if err := utils.VerifyToken(req.AccessToken); err != nil {
	// 	return err
	// }
	accessTokenUserID, accessTokenEmail, err := utils.ParseToken(req.AccessToken)
	if err != nil {
		return err
	}
	refresdhTokenUserID, refreshTokenEmail, err := utils.ParseToken(req.RefreshToken)
	if err != nil {
		return err
	}
	if accessTokenUserID != refresdhTokenUserID || accessTokenEmail != refreshTokenEmail {
		return utils.ErrorMsg(context.TODO(), utils.ErrBadParameter, utils.Trace(), "access token and refresh token are not matched", utils.ErrFromClient)
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
		return nil, utils.ErrorMsg(ctx, utils.ErrInternalServer, utils.Trace(), _errors.ErrInvalidGoogleCode.Error()+err.Error(), utils.ErrFromInternal)
	}
	resp, err := http.Get("https://www.googleapis.com/oauth2/v1/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, utils.ErrorMsg(ctx, utils.ErrInternalServer, utils.Trace(), _errors.ErrInvalidGoogleCode.Error()+err.Error(), utils.ErrFromInternal)
	}

	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, utils.ErrorMsg(ctx, utils.ErrInternalServer, utils.Trace(), _errors.ErrInvalidGoogleCode.Error()+err.Error(), utils.ErrFromInternal)
	}

	return content, nil
}

func getAppGoogleUserInfo(ctx context.Context, accessToken string) ([]byte, error) {
	fmt.Println(accessToken)
	token, err := utils.AppGoogleConfig.Exchange(ctx, accessToken)
	if err != nil {
		return nil, utils.ErrorMsg(ctx, utils.ErrInternalServer, utils.Trace(), _errors.ErrInvalidGoogleCode.Error()+err.Error(), utils.ErrFromInternal)
	}
	resp, err := http.Get("https://www.googleapis.com/oauth2/v1/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, utils.ErrorMsg(ctx, utils.ErrInternalServer, utils.Trace(), _errors.ErrInvalidGoogleCode.Error()+err.Error(), utils.ErrFromInternal)
	}

	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, utils.ErrorMsg(ctx, utils.ErrInternalServer, utils.Trace(), _errors.ErrInvalidGoogleCode.Error()+err.Error(), utils.ErrFromInternal)
	}

	return content, nil
}
func CreateUserProfileDTOList(userID uint, profileIDList []*mysql.GameProfiles) []*mysql.GameUserProfiles {
	userProfileDTOList := make([]*mysql.GameUserProfiles, 0)
	for _, profile := range profileIDList {
		userProfileDTOList = append(userProfileDTOList, &mysql.GameUserProfiles{
			UserID:     int(userID),
			ProfileID:  int(profile.ID),
			IsAchieved: true,
		})
	}
	return userProfileDTOList
}
