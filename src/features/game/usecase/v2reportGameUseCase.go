package usecase

import (
	"context"
	_interface "main/features/game/model/interface"
	"main/features/game/model/request"
	"main/utils/aws"
	"strconv"
	"time"
)

type V2ReportGameUseCase struct {
	Repository     _interface.IV2ReportGameRepository
	ContextTimeout time.Duration
}

func NewV2ReportGameUseCase(repo _interface.IV2ReportGameRepository, timeout time.Duration) _interface.IV2ReportGameUseCase {
	return &V2ReportGameUseCase{Repository: repo, ContextTimeout: timeout}
}

func (d *V2ReportGameUseCase) V2Report(c context.Context, userID uint, req *request.ReqV2Report) error {
	ctx, cancel := context.WithTimeout(c, d.ContextTimeout)
	defer cancel()

	//디비에 저장
	V2ReportDTO := CreateV2ReportDTO(userID, req)
	err := d.Repository.SaveReport(ctx, V2ReportDTO)
	if err != nil {
		return err
	}
	//이메일 전송
	reqV2Report := &aws.ReqReportSES{
		UserID:       strconv.Itoa(int(userID)),
		TargetUserID: strconv.Itoa(int(req.TargetUserID)),
		CategoryID:   strconv.Itoa(int(req.CategoryID)),
		Reason:       string(req.Reason),
	}
	go aws.EmailSendReport([]string{"pkjhj485@gmail.com", "kkukileon305@gmail.com"}, reqV2Report)

	return nil
}
