package _interface

import (
	"context"
	"main/features/game/model/request"
)

type IStartGameUseCase interface {
	Start(c context.Context, email string, req *request.ReqStart) error
}

type IDoraGameUseCase interface {
	Dora(c context.Context, userID int, req *request.ReqDora) error
}
type IOwnershipGameUseCase interface {
	Ownership(c context.Context, req *request.ReqOwnership) error
}

type IDiscardGameUseCase interface {
	Discard(c context.Context, userID int, req *request.ReqDiscard) error
}
