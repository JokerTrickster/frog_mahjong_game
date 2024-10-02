package usecase

import (
	"main/features/profiles/model/response"
	"main/utils/db/mysql"
)

func CreateResProfileList(profileList []*mysql.UserProfiles) response.ResListProfile {
	res := response.ResListProfile{}
	for _, profile := range profileList {
		res.ProfileList = append(res.ProfileList, response.Profile{
			ProfileID:    int(profile.ID),
			CurrentCount: profile.Earned,
			IsAchieved:   profile.IsAchieved,
		})
	}
	return res
}
