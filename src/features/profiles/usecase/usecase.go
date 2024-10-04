package usecase

import (
	"context"
	"fmt"
	"main/features/profiles/model/entity"
	"main/features/profiles/model/response"
	"main/utils/aws"
	_aws "main/utils/aws"
	"main/utils/db/mysql"
)

func CreateResProfileList(profileList []*mysql.Profiles) response.ResListProfile {
	res := response.ResListProfile{}
	for _, profile := range profileList {
		p := response.Profile{
			ProfileID:  int(profile.ID),
			Name:       profile.Name,
			TotalCount: profile.TotalCount,
			Image:      profile.Image,
		}
		imageUrl, err := _aws.ImageGetSignedURL(context.TODO(), profile.Image, aws.ImgTypeProfile)
		if err != nil {
			fmt.Println(err)
			return response.ResListProfile{}
		}
		p.Image = imageUrl

		res.Profiles = append(res.Profiles, p)
	}
	return res
}
func CreateProfileDTO(entity entity.ImageUploadProfileEntity, fileName string) *mysql.Profiles {
	return &mysql.Profiles{
		Name:        entity.Name,
		TotalCount:  entity.TotalCount,
		Description: entity.Description,
		Image:       fileName,
	}
}
