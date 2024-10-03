package usecase

import (
	"main/features/profiles/model/response"
	"main/utils/db/mysql"
)

func CreateResProfileList(profileList []*mysql.Profiles) response.ResListProfile {
	res := response.ResListProfile{}
	for _, profile := range profileList {
		res.Profiles = append(res.Profiles, response.Profile{
			ProfileID:  int(profile.ID),
			Name:       profile.Name,
			TotalCount: profile.TotalCount,
		})
	}
	return res
}
