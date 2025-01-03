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

type IResultGameHandler interface {
	Result(c echo.Context) error
}

type IReportGameHandler interface {
	Report(c echo.Context) error
}
type IMetaGameHandler interface {
	Meta(c echo.Context) error
}
type IDeckCardGameHandler interface {
	DeckCard(c echo.Context) error
}

type IListMissionGameHandler interface {
	ListMission(c echo.Context) error
}

type ICreateMissionGameHandler interface {
	CreateMission(c echo.Context) error
}

type IListCardGameHandler interface {
	ListCard(c echo.Context) error
}

type IV2ListCardGameHandler interface {
	V2ListCard(c echo.Context) error
}

// v2
type IV2DeckCardGameHandler interface {
	V2DeckCard(c echo.Context) error
}

type IV2ResultGameHandler interface {
	V2Result(c echo.Context) error
}

type IV2ReportGameHandler interface {
	V2Report(c echo.Context) error
}

type ISaveCardInfoGameHandler interface {
	SaveCardInfo(c echo.Context) error
}

type ISaveCardImageGameHandler interface {
	SaveCardImage(c echo.Context) error
}

type IUpdateCardGameHandler interface {
	UpdateCard(c echo.Context) error
}

type IReportImageUploadGameHandler interface {
	ReportImageUpload(c echo.Context) error
}
