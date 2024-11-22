package handler

import (
	"context"

	_interface "main/features/game/model/interface"
	"main/features/game/model/request"
	"main/utils"

	"net/http"

	"github.com/labstack/echo/v4"
)

type ReportImageUploadGameHandler struct {
	UseCase _interface.IReportImageUploadGameUseCase
}

func NewReportImageUploadGameHandler(c *echo.Echo, useCase _interface.IReportImageUploadGameUseCase) _interface.IReportImageUploadGameHandler {
	handler := &ReportImageUploadGameHandler{
		UseCase: useCase,
	}
	c.POST("/v0.1/games/images/upload-report", handler.ReportImageUpload)
	return handler
}

// 새 카드 이미지 업로드 리포트
// @Router /v0.1/games/images/upload-report [post]
// @Summary 새 카드 이미지 업로드 리포트
// @Description
// @Description ■ errCode with 400
// @Description PARAM_BAD : 파라미터 오류
// @Description USER_NOT_FOUND : 유저가 존재하지 않음
// @Description ■ errCode with 401
// @Description INVALID_AUTH_CODE : 인증 코드 검증 실패
// @Description TOKEN_BAD : 잘못된 토큰
// @Description INVALID_ACCESS_TOKEN : 잘못된 액세스 토큰
// @Description
// @Description ■ errCode with 500
// @Description INTERNAL_SERVER : 내부 로직 처리 실패
// @Description INTERNAL_DB : DB 처리 실패
// @Description GEMINI_INTERNAL_SERVER : Gemini 서버 내부 오류
// @Param type body request.ReqReportImageUploadGame true "type"
// @Produce json
// @Success 200 {object} bool
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags game
func (d *ReportImageUploadGameHandler) ReportImageUpload(c echo.Context) error {
	ctx := context.Background()
	req := &request.ReqReportImageUploadGame{}
	if err := utils.ValidateReq(c, req); err != nil {
		return err
	}

	//business logic
	err := d.UseCase.ReportImageUpload(ctx, req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, true)
}
