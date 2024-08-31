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

	//카드가 6장 소유했는지 체크하고 카드 정보를 가져온다.
	entitySQL := CreateScoreCalculateEntitySQL(userID, req)
	cardsDTO, err := d.Repository.FindOwnedCards(ctx, entitySQL)
	if err != nil {
		return 0, []string{}, err
	}

	//카드 검증을 한다.
	if err := CardValidation(cardsDTO, req.Cards); err != nil {
		return 0, []string{}, err
	}

	// dora 카드를 가져온다.
	doraCard, err := d.Repository.GetDoraCard(ctx, req)
	if err != nil {
		return 0, []string{}, err
	}

	// 요청받은 카드 순서에 맞게 엔티티를 만든다.
	entity := CreateScoreCalculateEntity(cardsDTO, req.Cards)

	//전달받은 카드들의 점수를 계산한다.
	score, bonuses, err := ScoreCalculate(entity, doraCard)
	if err != nil {
		return 0, []string{}, err
	}

	return score, bonuses, nil
}
