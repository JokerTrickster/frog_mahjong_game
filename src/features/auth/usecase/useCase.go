package usecase

import (
	"main/features/auth/model/request"
	"main/utils/db/mysql"
	"time"
)

func CreateSignupUserDTO(name, email string) (mysql.GormUserDTO, error) {
	result := mysql.GormUserDTO{
		GormModel: mysql.GormModel{
			ID: mysql.PKIDGenerate(),
		},
		Name:  name,
		Email: email,
	}
	return result, nil
}
func CreateSignupUserAuthDTO(userID string, req *request.ReqSignup) (mysql.GormUserAuthDTO, error) {
	now := time.Now()
	result := mysql.GormUserAuthDTO{
		GormModel: mysql.GormModel{
			ID: mysql.PKIDGenerate(),
		},
		Provider:   "email",
		Name:       req.Name,
		Email:      req.Email,
		Password:   req.Password,
		UserID:     userID,
		LastSignIn: mysql.TimeToEpoch(now),
	}
	return result, nil
}
