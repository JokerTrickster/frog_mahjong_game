package usecase

import (
	"context"
	_interface "main/features/game/model/interface"
	"main/features/game/model/request"
	"main/utils/aws"
	"time"
)

type SaveCardInfoGameUseCase struct {
	Repository     _interface.ISaveCardInfoGameRepository
	ContextTimeout time.Duration
}

func NewSaveCardInfoGameUseCase(repo _interface.ISaveCardInfoGameRepository, timeout time.Duration) _interface.ISaveCardInfoGameUseCase {
	return &SaveCardInfoGameUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *SaveCardInfoGameUseCase) SaveCardInfo(c context.Context, req *request.ReqSaveCardInfo) error {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()

	// bird cards DTO 생성
	birdCardsDTO := CreateBirdCardsDTO(req)

	// db 저장
	err := d.Repository.SaveCardInfo(ctx, birdCardsDTO)
	if err != nil {
		return err
	}
	cardList := []string{}
	for _, v := range birdCardsDTO {
		cardList = append(cardList, v.Name)
	}

	// email 전송
	aws.EmailSendCardInfo([]string{"pkjhj485@gmail.com", "kkukileon305@gmail.com", "ohhyejin1213@naver.com"}, cardList)

	return nil

}
