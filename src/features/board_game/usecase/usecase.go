package usecase

import (
	"context"
	"fmt"
	"main/features/board_game/model/entity"
	"main/features/board_game/model/response"
	_aws "main/utils/aws"

	"main/utils/db/mysql"
)

func CreateSoloPlayGameInfo(imageDTO *mysql.FindItImages, correctPositions []*mysql.FindItImageCorrectPositions, round int) response.SoloPlayGameInfo {
	res := response.SoloPlayGameInfo{
		ImageID: int(imageDTO.ID),
		Round:   round,
	}
	normalUrl, err := _aws.ImageGetSignedURL(context.TODO(), imageDTO.NormalImageUrl, _aws.ImgTypeFindIt)
	if err != nil {
		fmt.Println(err)
	}
	abnormalUrl, err := _aws.ImageGetSignedURL(context.TODO(), imageDTO.AbnormalImageUrl, _aws.ImgTypeFindIt)
	if err != nil {
		fmt.Println(err)
	}
	res.NormalUrl = normalUrl
	res.AbnormalUrl = abnormalUrl
	for _, correctPosition := range correctPositions {
		position := response.Position{
			X: correctPosition.XPosition,
			Y: correctPosition.YPosition,
		}
		res.CorrectPositions = append(res.CorrectPositions, position)
	}

	return res
}

func CreateRankUser(userDTO *mysql.GameUsers, correctDTO *entity.FindItRankEntity, rank int) response.RankUser {
	return response.RankUser{
		UserID:    int(userDTO.ID),
		Name:      userDTO.Name,
		Score:     correctDTO.CorrectCount,
		Rank:      rank,
		ProfileID: int(userDTO.ProfileID),
	}
}

func CreateSlimeWarRankUser(userDTO *mysql.GameUsers, correctDTO *entity.SlimeWarRankEntity, rank int) response.RankUser {
	return response.RankUser{
		UserID:    int(userDTO.ID),
		Name:      userDTO.Name,
		Score:     correctDTO.Score,
		Rank:      rank,
		ProfileID: int(userDTO.ProfileID),
	}
}
