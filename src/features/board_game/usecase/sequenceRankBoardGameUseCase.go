package usecase

import (
	"context"
	_interface "main/features/board_game/model/interface"
	"main/features/board_game/model/response"
	"time"
)

type SequenceRankBoardGameUseCase struct {
	Repository     _interface.ISequenceRankBoardGameRepository
	ContextTimeout time.Duration
}

func NewSequenceRankBoardGameUseCase(repo _interface.ISequenceRankBoardGameRepository, timeout time.Duration) _interface.ISequenceRankBoardGameUseCase {
	return &SequenceRankBoardGameUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *SequenceRankBoardGameUseCase) SequenceRank(c context.Context) (response.ResSequenceRank, error) {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()

	// 승리 횟수가 가장 많은 순으로 3명 랭킹을 가져온다.
	userCorrectDTOList, err := d.Repository.FindTop3User(ctx)
	if err != nil {
		return response.ResSequenceRank{}, err
	}
	res := response.ResSequenceRank{}
	rankUserList := make([]response.RankUser, 0)
	// 유저 정보를 가져온다.
	for i, userCorrectDTO := range userCorrectDTOList {
		// 유저 정보를 가져온다.
		userDTO, err := d.Repository.FindOneUser(ctx, userCorrectDTO.UserID)
		if err != nil {
			return response.ResSequenceRank{}, err
		}
		rankUser := CreateSequenceRankUser(userDTO, userCorrectDTO, i+1)
		rankUserList = append(rankUserList, rankUser)
	}
	res.RankUserList = rankUserList

	return res, nil
}
