package usecase

import (
	"context"
	"fmt"
	"main/features/game_profiles/model/entity"
	"main/features/game_profiles/model/response"
	"main/utils/aws"
	_aws "main/utils/aws"
	"main/utils/db/mysql"
)

func CreateResProfileList(profileList []*mysql.GameProfiles) response.ResListGameProfile {
	res := response.ResListGameProfile{}
	for _, profile := range profileList {
		p := response.Profile{
			ProfileID: int(profile.ID),
			Name:      profile.Name,
			Image:     profile.Image,
		}
		imageUrl, err := _aws.ImageGetSignedURL(context.TODO(), profile.Image, aws.ImgTypeProfile)
		if err != nil {
			fmt.Println(err)
			return response.ResListGameProfile{}
		}
		p.Image = imageUrl

		res.Profiles = append(res.Profiles, p)
	}
	return res
}
func CreateProfileDTO(entity entity.ImageUploadProfileEntity, fileName string) *mysql.GameProfiles {
	return &mysql.GameProfiles{
		Name:        entity.Name,
		Description: entity.Description,
		Image:       fileName,
	}
}
