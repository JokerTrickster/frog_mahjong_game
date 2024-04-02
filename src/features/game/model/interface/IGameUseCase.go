package _interface

import (
	"context"
	"main/features/game/model/request"
)

type IStartGameUseCase interface {
	Start(c context.Context, email string, req *request.ReqStart) error
}
