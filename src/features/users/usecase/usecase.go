package usecase

import (
	"main/features/users/model/entity"
	"main/features/users/model/request"
	"main/features/users/model/response"
	"main/utils/db/mysql"
)

func CreateUpdateUsersEntitySQL(userID uint, req *request.ReqUpdateUsers) entity.UpdateUsersEntitySQL {
	result := entity.UpdateUsersEntitySQL{
		UserID: userID,
	}
	if req.Name != "" {
		result.Name = req.Name
	}
	if req.Password != "" {
		result.Password = req.Password
	}
	if req.ProfileID != 0 {
		result.ProfileID = req.ProfileID
	}

	return result
}

func CreateResGetUser(userDTO mysql.Users) response.ResGetUser {
	return response.ResGetUser{
		UserID:    int(userDTO.ID),
		Email:     userDTO.Email,
		Name:      userDTO.Name,
		Coin:      userDTO.Coin,
		ProfileID: userDTO.ProfileID,
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

func CreateResProfileList(profileList []*mysql.UserProfiles) response.ResListProfileUser {
	res := response.ResListProfileUser{}
	for _, profile := range profileList {
		res.Profiles = append(res.Profiles, response.Profile{
			ProfileID:    int(profile.ProfileID),
			CurrentCount: profile.CurrentCount,
			IsAchieved:   profile.IsAchieved,
		})
	}
	return res
}
