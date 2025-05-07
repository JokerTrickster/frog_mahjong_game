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
type ISlimeWarRankBoardGameHandler interface {
	SlimeWarRank(c echo.Context) error
}

type ISequenceResultBoardGameHandler interface {
	SequenceResult(c echo.Context) error
}

type ISequenceRankBoardGameHandler interface {
	SequenceRank(c echo.Context) error
}
type IGameOverBoardGameHandler interface {
	GameOver(c echo.Context) error
}
