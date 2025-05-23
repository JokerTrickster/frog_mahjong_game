package handler

import (
	"context"
	"encoding/json"
	"fmt"
	_interface "main/features/board_game/model/interface"
	"main/features/board_game/model/request"
	"main/features/board_game/model/response"
	"main/utils"
	_redis "main/utils/db/redis"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type FrogCardListBoardGameHandler struct {
	UseCase _interface.IFrogCardListBoardGameUseCase
}

func NewFrogCardListBoardGameHandler(c *echo.Echo, useCase _interface.IFrogCardListBoardGameUseCase) _interface.IFrogCardListBoardGameHandler {
	handler := &FrogCardListBoardGameHandler{
		UseCase: useCase,
	}
	c.GET("/frog/v0.1/game/cards", handler.FrogCardList)
	return handler
}

// 개굴작 카드 정보 리스트 가져오기
// @Router /frog/v0.1/game/cards  [get]
// @Summary 개굴작 카드 정보 리스트 가져오기
// @Description
// @Description ■ errCode with 400
// @Description PARAM_BAD : 파라미터 오류
// @Description NOT_OWNER : 방장이 시작 요청을 하지 않음
// @Description NOT_FIRST_PLAYER : 첫 플레이어가 아님
// @Description
// @Description ■ errCode with 500
// @Description INTERNAL_SERVER : 내부 로직 처리 실패
// @Description INTERNAL_DB : DB 처리 실패
// @Produce json
// @Success 200 {object} response.ResFrogCardListBoardGame
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Tags app/board-game/game
func (d *FrogCardListBoardGameHandler) FrogCardList(c echo.Context) error {
	ctx := context.Background()
	req := &request.ReqFrogCardListBoardGame{}
	if err := utils.ValidateReq(c, req); err != nil {
		return err
	}
	//business logic
	//redis
	cacheKey := fmt.Sprintf("game:frog:cards:%d", req.RoomID)
	cardData, err := _redis.Client.Get(ctx, cacheKey).Result()
	if cardData == "" {
		// 2. 캐시에 데이터가 없을 경우 UseCase에서 조회
		res, err := d.UseCase.FrogCardList(ctx)
		if err != nil {
			return err
		}

		// 3. 조회된 데이터를 Redis에 캐시 (예: 1시간 TTL)
		data, err := json.Marshal(res)
		if err != nil {
			return err
		}
		err = _redis.Client.Set(ctx, cacheKey, data, 1*time.Hour).Err()
		if err != nil {
			return err
		}

		// 캐시 히트 여부
		c.Response().Header().Set("X-Cache-Hit", "false")

		// 4. DB에서 조회한 데이터 반환
		return c.JSON(http.StatusOK, res)
	} else if err != nil {
		// Redis 오류 처리
		return err
	}

	// 5. 캐시된 데이터가 있을 경우 반환
	var res response.ResFrogCardListBoardGame
	if err := json.Unmarshal([]byte(cardData), &res); err != nil {
		return err
	}

	// 캐시 히트 여부
	c.Response().Header().Set("X-Cache-Hit", "true")

	return c.JSON(http.StatusOK, res)
}
