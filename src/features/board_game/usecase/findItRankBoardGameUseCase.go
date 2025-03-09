package usecase

import (
	"context"
	_interface "main/features/board_game/model/interface"
	"main/features/board_game/model/response"
	"time"
)

type FindItRankBoardGameUseCase struct {
	Repository     _interface.IFindItRankBoardGameRepository
	ContextTimeout time.Duration
}

func NewFindItRankBoardGameUseCase(repo _interface.IFindItRankBoardGameRepository, timeout time.Duration) _interface.IFindItRankBoardGameUseCase {
	return &FindItRankBoardGameUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *FindItRankBoardGameUseCase) FindItRank(c context.Context) (response.ResFindItRankBoardGame,error) {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()

	// 맞춘 횟수가 가장 많은 순으로 3명 랭킹을 가져온다.
	userCorrectDTOList, err := d.Repository.FindTop3UserCorrect(ctx)
	if err != nil {
		return response.ResFindItRankBoardGame{}, err
	}
	res := response.ResFindItRankBoardGame{}
	rankUserList := make([]response.RankUser, 0)
	// 유저 정보를 가져온다.
	for i, userCorrectDTO := range userCorrectDTOList {
		// 유저 정보를 가져온다.
		userDTO, err := d.Repository.FindOneUser(ctx, userCorrectDTO.UserID)
		if err != nil {
			return response.ResFindItRankBoardGame{}, err
		}
		rankUser := CreateRankUser(userDTO, userCorrectDTO, i+1)
		rankUserList = append(rankUserList, rankUser)
	}
	res.RankUserList = rankUserList
	return res, nil
}
