package usecase

import (
	"context"
	_interface "main/features/game/model/interface"
	"main/features/game/model/request"
	"time"
)

type ReportGameUseCase struct {
	Repository     _interface.IReportGameRepository
	ContextTimeout time.Duration
}

func NewReportGameUseCase(repo _interface.IReportGameRepository, timeout time.Duration) _interface.IReportGameUseCase {
	return &ReportGameUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *ReportGameUseCase) Report(c context.Context, userID uint, req *request.ReqReport) error {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()

	//디비에 저장
	reportDTO := CreateReportDTO(userID, req)
	err := d.Repository.SaveReport(ctx, reportDTO)
	if err != nil {
		return err
	}
	//이메일 전송
	

	return nil
}
