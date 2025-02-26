package usecase

import (
	"main/features/game_users/model/entity"
	"main/features/game_users/model/request"
	"main/features/game_users/model/response"
	"main/utils/db/mysql"
)

func CreateUpdateUsersEntitySQL(userID uint, req *request.ReqUpdateGameUsers) entity.UpdateUsersEntitySQL {
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

func CreateResGetUser(userDTO mysql.GameUsers, disconnected int64) response.ResGetGameUser {
	result := response.ResGetGameUser{
		UserID:       int(userDTO.ID),
		Email:        userDTO.Email,
		Name:         userDTO.Name,
		Coin:         userDTO.Coin,
		ProfileID:    userDTO.ProfileID,
		Disconnected: disconnected,
	}

	return result
}

func CreateResListUser(users []mysql.GameUsers, total int) response.ResListGameUser {
	res := response.ResListGameUser{}
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

func CreateResProfileList(profileList []*mysql.GameUserProfiles) response.ResListProfileGameUser {
	res := response.ResListProfileGameUser{}
	for _, profile := range profileList {
		res.Profiles = append(res.Profiles, response.Profile{
			ProfileID:    int(profile.ProfileID),
			IsAchieved:   profile.IsAchieved,
		})
	}
	return res
}
