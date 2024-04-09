package _interface

import "github.com/labstack/echo/v4"

type IStartGameHandler interface {
	Start(c echo.Context) error
}

type IDoraGameHandler interface {
	Dora(c echo.Context) error
}

type IOwnershipGameHandler interface {
	Ownership(c echo.Context) error
}

type IDiscardGameHandler interface {
	Discard(c echo.Context) error
}

type INextTurnGameHandler interface {
	NextTurn(c echo.Context) error
}

type ILoanGameHandler interface {
	Loan(c echo.Context) error
}

type IScoreCalculateGameHandler interface {
	ScoreCalculate(c echo.Context) error
}

type IWinRequestGameHandler interface {
	WinRequest(c echo.Context) error
}
