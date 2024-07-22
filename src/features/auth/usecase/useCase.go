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


// func NativeValidate(ctx context.Context, tokenString string) (domain.OAuthData, error) {
// 	oauthData := domain.OAuthData{}
// 	claims, iErr := common.JwtVerifyWithKeySet(ctx, "google", tokenString, "https://www.googleapis.com/oauth2/v3/certs", common.ErrFromProviderGoogle)
// 	if iErr != nil {
// 		return domain.OAuthData{}, iErr
// 	}
// 	aud, okAud := claims["aud"].(string)
// 	azp, _ := claims["azp"].(string)
// 	iss, okIss := claims["iss"].(string)
// 	sub, okSub := claims["sub"].(string)
// 	email, okEmail := claims["email"].(string)
// 	hd, okHd := claims["hd"].(string)

// 	if !okAud || !okIss || !okSub || !okEmail || !okHd ||
// 		(aud != domain.WAuthMeta.GoogleClientID && azp != domain.WAuthMeta.GoogleClientID) ||
// 		(iss != "accounts.google.com" && iss != "https://accounts.google.com") ||
// 		(hd != "breathings.co.kr") {
// 		return domain.OAuthData{}, errorSystem.ErrorMsg(errorSystem.ErrBadToken, common.Trace(), fmt.Sprintf("not valid token claims from provider Google - %+v", claims), errorSystem.ErrFromClient)
// 	}
// 	oauthData = domain.OAuthData{
// 		ID:    sub,
// 		Email: email,
// 	}
// 	oauthData.Provider = domain.AuthProviderName[0]
// 	if err := common.ValidateStruct(oauthData); err != nil {
// 		return domain.OAuthData{}, errorSystem.ErrorMsg(errorSystem.ErrInternalServer, common.Trace(), fmt.Sprintf("wrong generated oauthData - %+v", oauthData), errorSystem.ErrFromInternal)
// 	}

// 	return oauthData, nil
// }


// type UserOAuth struct {
// 	ID             primitive.ObjectID `bson:"_id,omitempty"`
// 	Email          string             `bson:"email,omitempty" validate:"require,email"`
// 	LastSignin     time.Time          `bson:"lastSignin,omitempty" validate:"require"`
// 	LastTokenIssue time.Time          `bson:"lastTokenIssue,omitempty"` // 유저 마지막 토큰 발행시간
// 	Provider       string             `bson:"provider,omitempty" validate:"required,oneof=google"`
// 	ProviderID     string             `bson:"providerID,omitempty" validate:"required"`
// 	Role           string             `bson:"role,omitempty"` // 유저 권한
// }

// type WebAuthMeta struct {
// 	GoogleClientID  string `json:"google_client_id" validate:"required"`
// 	GoogleSecretKey string `json:"google_secret_key" validate:"required"`
// 	Hd              string `json:"hd" validate:"required"`
// }

// const (
// 	CallBackURL  = "https://dev-admin-api.breathings.net/v0.1/auth/google/signin/callback"
// 	ScopeEmail   = "https://www.googleapis.com/auth/userinfo.email"
// 	ScopeProfile = "https://www.googleapis.com/auth/userinfo.profile"
// )

// var OAuthConf *oauth2.Config
// var WAuthMeta WebAuthMeta

// type OAuthData struct {
// 	ID       string `validate:"required"`
// 	Email    string `validate:"omitempty,email"`
// 	Provider string `validate:"required"`
// }

// type AuthProvider uint8

// const (
// 	AuthProviderGoogle = AuthProvider(0)
// )

// var AuthProviderName map[AuthProvider]string = map[AuthProvider]string{
// 	AuthProviderGoogle: "google",
// }
