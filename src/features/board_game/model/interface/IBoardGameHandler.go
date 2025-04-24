package _interface

import "github.com/labstack/echo/v4"

type IFindItSoloPlayBoardGameHandler interface {
	FindItSoloPlay(c echo.Context) error
}

type IFindItRankBoardGameHandler interface {
	FindItRank(c echo.Context) error
}
type IFindItCoinBoardGameHandler interface {
	FindItCoin(c echo.Context) error
}

type IFindItPasswordCheckBoardGameHandler interface {
	FindItPasswordCheck(c echo.Context) error
}

type ISlimeWarGetsCardBoardGameHandler interface {
	SlimeWarGetsCard(c echo.Context) error
}

type ISlimeWarResultBoardGameHandler interface {
	SlimeWarResult(c echo.Context) error
}
