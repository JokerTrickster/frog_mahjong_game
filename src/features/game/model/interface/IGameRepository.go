package _interface

import (
	"context"
	"main/features/game/model/request"
	"main/utils/db/mysql"
)

type IStartGameRepository interface {
	CheckOwner(c context.Context, email string, roomID uint) error
	CheckReady(c context.Context, roomID uint) ([]mysql.RoomUsers, error)
	UpdateRoomUser(c context.Context, roomID uint, state string) error
	UpdateRoom(c context.Context, roomID uint, state string) error
	CreateCards(c context.Context, roomID uint, cards []mysql.Cards) error
}

type IDoraGameRepository interface {
	CheckOwner(c context.Context, userID int, roomID int) error
	UpdateDoraCard(c context.Context, req *request.ReqDora) error
}