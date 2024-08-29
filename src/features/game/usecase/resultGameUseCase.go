package usecase

import (
	"context"
	_interface "main/features/game/model/interface"
	"main/features/game/model/request"
	"main/features/game/model/response"
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

func (d *ResultGameUseCase) Result(c context.Context, userID uint, req *request.ReqResult) (response.ResResult, error) {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()

	if len(req.Cards) == 0 {
		//무승부 처리
		return response.ResResult{}, nil
	} else if len(req.Cards) != 6 {
		//카드가 6장이 아닌 경우 에러 처리
		return response.ResResult{}, utils.ErrorMsg(c, utils.ErrNotEnoughCard, utils.Trace(), "cards length is not 6", utils.ErrFromClient)
	}
	entitySQL := CreateResultEntitySQL(userID, req)

	// 카드 정보를 가져온다.
	cardsDTO, err := d.Repository.FindCards(ctx, entitySQL)
	if err != nil {
		return response.ResResult{}, err
	}

	// dora 카드를 가져온다.
	doraCard, err := d.Repository.GetDoraCard(ctx, req)
	if err != nil {
		return response.ResResult{}, err
	}

	// 요청받은 카드 순서에 맞게 엔티티를 만든다.
	entity := CreateResultEntity(cardsDTO, req.Cards)

	//전달받은 카드들의 점수를 계산한다.
	score, bonuses, err := ScoreCalculate(entity, doraCard)
	if err != nil {
		return response.ResResult{}, err
	}
	res := response.ResResult{
		Score:   score,
		Winner:  uint64(cardsDTO[0].UserID),
		Bonuses: bonuses,
	}

	return res, nil
}
