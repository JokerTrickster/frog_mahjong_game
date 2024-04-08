package usecase

import (
	"context"
	"fmt"
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

func (d *ScoreCalculateGameUseCase) ScoreCalculate(c context.Context, userID uint, req *request.ReqScoreCalculate) error {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()
	fmt.Println(ctx)

	return nil
}
