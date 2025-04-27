package usecase

import (
	"context"
	_interface "main/features/board_game/model/interface"
	"main/features/board_game/model/response"
	"time"
)

type SlimeWarRankBoardGameUseCase struct {
	Repository     _interface.ISlimeWarRankBoardGameRepository
	ContextTimeout time.Duration
}

func NewSlimeWarRankBoardGameUseCase(repo _interface.ISlimeWarRankBoardGameRepository, timeout time.Duration) _interface.ISlimeWarRankBoardGameUseCase {
	return &SlimeWarRankBoardGameUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *SlimeWarRankBoardGameUseCase) SlimeWarRank(c context.Context) (response.ResSlimeWarRankBoardGame, error) {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()

	// 승리 횟수가 가장 많은 순으로 3명 랭킹을 가져온다.
	userCorrectDTOList, err := d.Repository.FindTop3User(ctx)
	if err != nil {
		return response.ResSlimeWarRankBoardGame{}, err
	}
	res := response.ResSlimeWarRankBoardGame{}
	rankUserList := make([]response.RankUser, 0)
	// 유저 정보를 가져온다.
	for i, userCorrectDTO := range userCorrectDTOList {
		// 유저 정보를 가져온다.
		userDTO, err := d.Repository.FindOneUser(ctx, userCorrectDTO.UserID)
		if err != nil {
			return response.ResSlimeWarRankBoardGame{}, err
		}
		rankUser := CreateSlimeWarRankUser(userDTO, userCorrectDTO, i+1)
		rankUserList = append(rankUserList, rankUser)
	}
	res.RankUserList = rankUserList

	return res, nil
}
