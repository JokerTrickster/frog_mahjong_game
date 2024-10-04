package usecase

import (
	"context"
	"encoding/base64"
	"io/ioutil"
	"main/features/auth/model/request"
	_errors "main/features/game/model/errors"
	"main/utils"
	"main/utils/db/mysql"
	"net/http"
	"time"

	"golang.org/x/exp/rand"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// 랜덤 값 생성 함수
func GeneratePasswordAuthCode() string {
	seed := rand.NewSource(uint64(time.Now().UnixNano()))
	r := rand.New(seed)
	b := make([]byte, 6)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b)
}

func CreateUserSQL(email string) *mysql.Users {
	return &mysql.Users{
		Name:      "임시개굴맨",
		Email:     email,
		State:     "wait",
		Coin:      30,
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

func CreateSignupUser(req *request.ReqSignup) *mysql.Users {
	return &mysql.Users{
		Name:      req.Name,
		Email:     req.Email,
		Password:  req.Password,
		Coin:      30,
		RoomID:    1,
		State:     "logout",
		Provider:  "email",
		ProfileID: 1,
	}
}

func VerifyAccessAndRefresh(req *request.ReqReissue) error {
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

func CreateUserProfileDTOList(userID uint, profileIDList []*mysql.Profiles) []*mysql.UserProfiles {
	userProfileDTOList := make([]*mysql.UserProfiles, 0)
	for _, profile := range profileIDList {
		userProfileDTOList = append(userProfileDTOList, &mysql.UserProfiles{
			UserID:       int(userID),
			ProfileID:    int(profile.ID),
			IsAchieved:   true,
			CurrentCount: 0,
		})
	}
	return userProfileDTOList
}
