package usecase

import (
	"context"
	_interface "main/features/game/model/interface"
	"main/features/game/model/request"
	"main/utils"
	"time"
)

type ResultGameUseCase struct {
	Repository     _interface.IResultGameRepository
	ContextTimeout time.Duration
}

func NewResultGameUseCase(repo _interface.IResultGameRepository, timeout time.Duration) _interface.IResultGameUseCase {
	return &ResultGameUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *ResultGameUseCase) Result(c context.Context, userID uint, req *request.ReqResult) (int, []string, error) {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()

	//카드가 6장 소유했는지 체크하고 카드 정보를 가져온다.
	if len(req.Cards) != 6 {
		return 0, []string{}, utils.ErrorMsg(c, utils.ErrNotEnoughCard, utils.Trace(), "cards length is not 6", utils.ErrFromClient)
	}
	entitySQL := CreateResultEntitySQL(userID, req)

	// 카드 정보를 가져온다.
	cardsDTO, err := d.Repository.FindCards(ctx, entitySQL)
	if err != nil {
		return 0, []string{}, err
	}

	// dora 카드를 가져온다.
	doraCard, err := d.Repository.GetDoraCard(ctx, req)
	if err != nil {
		return 0, []string{}, err
	}

	// 요청받은 카드 순서에 맞게 엔티티를 만든다.
	entity := CreateResultEntity(cardsDTO, req.Cards)

	//전달받은 카드들의 점수를 계산한다.
	score, bonuses, err := ScoreCalculate(entity, doraCard)
	if err != nil {
		return 0, []string{}, err
	}

	return score, bonuses, nil
}
