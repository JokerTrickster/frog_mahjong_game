package usecase

import (
	"context"
	_interface "main/features/board_game/model/interface"
	"main/features/board_game/model/response"
	"main/utils/aws"
	"time"
)

type SlimeWarGetsCardBoardGameUseCase struct {
	Repository     _interface.ISlimeWarGetsCardBoardGameRepository
	ContextTimeout time.Duration
}

func NewSlimeWarGetsCardBoardGameUseCase(repo _interface.ISlimeWarGetsCardBoardGameRepository, timeout time.Duration) _interface.ISlimeWarGetsCardBoardGameUseCase {
	return &SlimeWarGetsCardBoardGameUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *SlimeWarGetsCardBoardGameUseCase) SlimeWarGetsCard(c context.Context) (response.ResSlimeWarGetsCardBoardGame, error) {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()

	cardList, err := d.Repository.FindCardList(ctx)
	if err != nil {
		return response.ResSlimeWarGetsCardBoardGame{}, err
	}

	resCardList := make([]response.Card, 0)
	for _, card := range cardList {
		// 서명된 Url로 개선 필요
		imageUrl, err := aws.ImageGetSignedURL(ctx, card.Image, aws.ImgTypeSlimeWar)
		if err != nil {
			return response.ResSlimeWarGetsCardBoardGame{}, err
		}
		resCardList = append(resCardList, response.Card{ID: int(card.ID), Direction: card.Direction, ImageUrl: imageUrl, Move: card.Move})
	}

	return response.ResSlimeWarGetsCardBoardGame{CardList: resCardList}, nil
}
