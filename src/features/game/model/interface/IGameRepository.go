package _interface

import (
	"context"
	"main/features/game/model/request"
	"main/utils/db/mysql"
)

type IStartGameRepository interface {
	CheckOwner(c context.Context, email string, roomID uint) error
	CheckReady(c context.Context, roomID uint) ([]mysql.RoomUsers, error)
	UpdateRoomUser(c context.Context, updateRoomUsers []mysql.RoomUsers) error
	UpdateRoom(c context.Context, roomID uint, state string) error
	CreateCards(c context.Context, roomID uint, cards []mysql.Cards) error
}

type IDoraGameRepository interface {
	CheckFirstPlayer(c context.Context, userID int, roomID int) error
	UpdateDoraCard(c context.Context, req *request.ReqDora) error
}
type IOwnershipGameRepository interface {
	UpdateCardState(c context.Context, req *request.ReqOwnership) error
	UpdateRoomUserCardCount(c context.Context, req *request.ReqOwnership) error
}

type IDiscardGameRepository interface {
	PlayerCheckTurn(c context.Context, req *request.ReqDiscard) (mysql.RoomUsers, error)
	UpdateCardStateDiscard(c context.Context, req *request.ReqDiscard) error
	UpdateRoomUser(c context.Context, updateRoomUser mysql.RoomUsers) error
}

type INextTurnGameRepository interface {
	UpdatePlayerNextTurn(c context.Context, req *request.ReqNextTurn) error
}

type ILoanGameRepository interface {
}
