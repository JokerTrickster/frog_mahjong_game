package _interface

import (
	"context"
)

type IGetUsersUseCase interface {
	Get(c context.Context) error
}
