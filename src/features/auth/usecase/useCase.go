package usecase

import (
	"main/features/auth/model/request"
	"main/utils"
	"main/utils/db/mysql"
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
