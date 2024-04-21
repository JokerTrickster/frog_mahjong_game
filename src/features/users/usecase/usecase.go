package usecase

import (
	"main/features/users/model/response"
	"main/utils/db/mysql"
)

func CreateResGetUser(userDTO mysql.Users) response.ResGetUser {
	return response.ResGetUser{
		UserID: int(userDTO.ID),
		Email:  userDTO.Email,
		Name:   userDTO.Name,
	}
}
