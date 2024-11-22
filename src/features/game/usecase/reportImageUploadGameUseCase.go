package usecase

import (
	"context"
	_interface "main/features/game/model/interface"

	"main/features/game/model/request"
	"main/utils/aws"

	"time"
)

type ReportImageUploadGameUseCase struct {
	Repository     _interface.IReportImageUploadGameRepository
	ContextTimeout time.Duration
}

func NewReportImageUploadGameUseCase(repo _interface.IReportImageUploadGameRepository, timeout time.Duration) _interface.IReportImageUploadGameUseCase {
	return &ReportImageUploadGameUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *ReportImageUploadGameUseCase) ReportImageUpload(c context.Context, req *request.ReqReportImageUploadGame) error {
	_, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()

	aws.EmailSendCardUploadReport([]string{"pkjhj485@gmail.com", "kkukileon305@gmail.com", "ohhyejin1213@naver.com"}, req.SuccessList, req.FailedList)

	return nil
}
