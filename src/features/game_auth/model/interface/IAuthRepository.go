package _interface

import (
	"context"
	"main/features/game_auth/model/entity"
	"main/utils/db/mysql"
)

type ISignupAuthRepository interface {
	UserCheckByEmail(ctx context.Context, email string) error
	InsertOneUser(ctx context.Context, user *mysql.GameUsers) error
	VerifyAuthCode(ctx context.Context, email, code string) error
	FindAllBasicProfile(ctx context.Context) ([]*mysql.Profiles, error)
	InsertOneUserProfile(ctx context.Context, userProfileDTOList []*mysql.UserProfiles) error
}

type ISigninAuthRepository interface {
	FindOneAndUpdateUser(ctx context.Context, email, password string) (mysql.GameUsers, error)
	SaveToken(ctx context.Context, uID uint, accessToken, refreshToken string, refreshTknExpiredAt int64) error
	DeleteToken(ctx context.Context, uID uint) error
	CheckToken(ctx context.Context, uID uint) (*mysql.Tokens, error)
}

type ILogoutAuthRepository interface {
	FindOneAndUpdateUser(ctx context.Context, uID uint) error
	DeleteToken(ctx context.Context, uID uint) error
}

type IReissueAuthRepository interface {
	SaveToken(ctx context.Context, token mysql.Tokens) error
	DeleteToken(ctx context.Context, uID uint) error
	CheckToken(ctx context.Context, uID uint, refreshToken string) error
}

type IGoogleOauthAuthRepository interface {
}

type IGoogleOauthCallbackAuthRepository interface {
	FindOneAndUpdateUser(ctx context.Context, googleOauthCallbackSQLQuery *entity.GoogleOauthCallbackSQLQuery) (*mysql.GameUsers, error)
	SaveToken(ctx context.Context, uID uint, accessToken, refreshToken string, refreshTknExpiredAt int64) error
	DeleteToken(ctx context.Context, uID uint) error
	CreateUser(ctx context.Context, user *mysql.GameUsers) (*mysql.GameUsers, error)
}

type IRequestPasswordAuthRepository interface {
	FindOneUserByEmail(ctx context.Context, email string) error
	InsertAuthCode(ctx context.Context, userAuthDTO mysql.UserAuths) error
}

type IRequestSignupAuthRepository interface {
	FindOneUserByEmail(ctx context.Context, email string) error
	InsertAuthCode(ctx context.Context, userAuthDTO mysql.UserAuths) error
	DeleteAuthCodeByEmail(ctx context.Context, email string) error
}
type IValidatePasswordAuthRepository interface {
	CheckAuthCode(ctx context.Context, email, code string) error
	UpdatePassword(ctx context.Context, user mysql.GameUsers) error
	DeleteAuthCode(ctx context.Context, email string) error
}

type IV02GoogleOauthCallbackAuthRepository interface {
	FindOneAndUpdateUser(ctx context.Context, googleOauthCallbackSQLQuery *entity.V02GoogleOauthCallbackSQLQuery) (*mysql.GameUsers, error)
	SaveToken(ctx context.Context, uID uint, accessToken, refreshToken string, refreshTknExpiredAt int64) error
	DeleteToken(ctx context.Context, uID uint) error
	CreateUser(ctx context.Context, user *mysql.GameUsers) (*mysql.GameUsers, error)
	FindAllBasicProfile(ctx context.Context) ([]*mysql.Profiles, error)
	InsertOneUserProfile(ctx context.Context, userProfileDTOList []*mysql.UserProfiles) error
	CheckToken(ctx context.Context, uID uint) (*mysql.Tokens, error)
}

type IFCMTokenAuthRepository interface {
	SaveFCMToken(ctx context.Context, userID uint, token string) error
}

type ICheckSigninAuthRepository interface {
	FindOneAndUpdateUser(ctx context.Context, email, password string) (mysql.GameUsers, error)
	CheckToken(ctx context.Context, uID uint) (*mysql.Tokens, error)
}

type INameCheckAuthRepository interface {
	CheckName(ctx context.Context, name string) error
}
