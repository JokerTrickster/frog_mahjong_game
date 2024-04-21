package handler

import (
	"context"
	_interface "main/features/auth/model/interface"
	"main/features/auth/model/request"
	"main/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ReissueAuthHandler struct {
	UseCase _interface.IReissueAuthUseCase
}

func NewReissueAuthHandler(c *echo.Echo, useCase _interface.IReissueAuthUseCase) _interface.IReissueAuthHandler {
	handler := &ReissueAuthHandler{
		UseCase: useCase,
	}
	c.POST("/v0.1/auth/reissue", handler.Reissue)
	return handler
}

// 액세스 토큰 재발급
// @Router /v0.1/auth/reissue [post]
// @Summary 액세스 토큰 재발급
// @Description
// @Description ■ errCode with 400
// @Description PARAM_BAD : 파라미터 오류
// @Description USER_NOT_EXIST : 유저가 존재하지 않음
// @Description USER_ALREADY_EXISTED : 유저가 이미 존재
// @Description
// @Description ■ errCode with 500
// @Description INTERNAL_SERVER : 내부 로직 처리 실패
// @Description INTERNAL_DB : DB 처리 실패
// @Param json body request.ReqReissue true "액세스 토큰, 리프레시 토큰"
// @Produce json
// @Success 200 {object} response.ResReissue
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags auth
func (d *ReissueAuthHandler) Reissue(c echo.Context) error {
	ctx := context.Background()
	req := &request.ReqReissue{}
	if err := utils.ValidateReq(c, req); err != nil {
		return err
	}
	res, err := d.UseCase.Reissue(ctx, req)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, res)
}
