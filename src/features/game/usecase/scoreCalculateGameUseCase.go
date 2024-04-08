package usecase

import (
	"context"
	_interface "main/features/game/model/interface"
	"main/features/game/model/request"
	"time"
)

type ScoreCalculateGameUseCase struct {
	Repository     _interface.IScoreCalculateGameRepository
	ContextTimeout time.Duration
}

func NewScoreCalculateGameUseCase(repo _interface.IScoreCalculateGameRepository, timeout time.Duration) _interface.IScoreCalculateGameUseCase {
	return &ScoreCalculateGameUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *ScoreCalculateGameUseCase) ScoreCalculate(c context.Context, userID uint, req *request.ReqScoreCalculate) (int, []string, error) {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()

	//카드가 6장 소유했는지 체크
	err := d.Repository.CheckCardCount(ctx, userID, req)
	if err != nil {
		return 0, []string{}, err
	}

	// dora 카드를 가져온다.
	doraCard, err := d.Repository.GetDoraCard(ctx, req)
	if err != nil {
		return 0, []string{}, err
	}

	//전달받은 카드들의 점수를 계산한다.
	score, bonuses, err := ScoreCalculate(req, doraCard)
	if err != nil {
		return 0, []string{}, err
	}

	return score, bonuses, nil
}
