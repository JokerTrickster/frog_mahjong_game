package usecase

import (
	"main/features/auth/model/request"
	"main/utils/db/mysql"
)

func CreateSignupUser(req *request.ReqSignup) mysql.Users {
	return mysql.Users{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Score:    30,
		RoomID:   1,
		State:    "ready",
	}
}
