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

func CreateResListUser(users []mysql.Users, total int) response.ResListUser {
	res := response.ResListUser{}
	UserList := make([]response.User, 0, len(users))
	for i := 0; i < len(users); i++ {
		user := response.User{
			UserID: int(users[i].ID),
			Name:   users[i].Name,
			Email:  users[i].Email,
			State:  users[i].State,
			Coin:   users[i].Coin,
		}
		UserList = append(UserList, user)
	}
	res.Users = UserList
	res.Total = total
	return res
}
