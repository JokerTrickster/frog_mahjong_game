package usecase

// createRoomUseCase.go test code write
// Path: src/features/room/usecase/createRoomUseCase_test.go

import (
	"context"
	"main/features/room/model/interface/mocks"
	"main/features/room/model/request"
	"main/utils"
	"testing"
	"time"

	_errors "main/features/room/model/errors"

	"github.com/stretchr/testify/mock"
	"gopkg.in/go-playground/assert.v1"
)

func TestCreateRoomUseCase_Create(t *testing.T) {

	// CreateRoomUseCase.Create test code write
	// Path: src/features/room/usecase/createRoomUseCase_test.go
	tests := []struct {
		name string
		req  request.ReqCreate
		err  error
	}{
		{"success1", request.ReqCreate{Name: "test", MaxCount: 4, MinCount: 2, Password: "test"}, nil},
		{"fail1", request.ReqCreate{Name: "test", MaxCount: 1, MinCount: 2, Password: ""}, utils.ErrorMsg(context.TODO(), utils.ErrUserNotFound, utils.Trace(), _errors.ErrBadRequest.Error(), utils.ErrFromClient)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//given
			mockCreateRoomRepository := new(mocks.ICreateRoomRepository)
			if tt.err == nil {
				mockCreateRoomRepository.On("InsertOneRoom", mock.Anything, mock.Anything).Return(1, nil) //mock
			}

			mockCreateRoomRepository.On("InsertOneRoom", mock.Anything, mock.Anything).Return(0, tt.err) //mock
			mockCreateRoomRepository.On("InsertOneRoomUser", mock.Anything, mock.Anything).Return(nil)   //mock
			mockCreateRoomRepository.On("FindOneAndUpdateUser", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			us := NewCreateRoomUseCase(mockCreateRoomRepository, 8*time.Second)
			//when
			err := us.Create(context.TODO(), 1, "ryan@gmail.com", &tt.req)
			//then
			assert.Equal(t, tt.err, err)
		})
	}
}

// CreateRoomDTO test code write
// Path: src/features/room/usecase/createRoomUseCase_test.go
func TestCreateRoomDTO(t *testing.T) {
	tests := []struct {
		name  string
		req   request.ReqCreate
		email string
		err   error
	}{
		{"success1", request.ReqCreate{Name: "test", MaxCount: 4, MinCount: 2, Password: ""}, "ryan@gmail.com", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//given
			//when
			_, err := CreateRoomDTO(context.TODO(), &tt.req, tt.email)
			//then
			assert.Equal(t, tt.err, err)
		})
	}
}

// CreateRoomUserDTO test code write
// Path: src/features/room/usecase/createRoomUseCase_test.go
func TestCreateRoomUserDTO(t *testing.T) {
	// write table driven test
	tests := []struct {
		name   string
		uID    uint
		roomID int
		err    error
	}{
		{"success1", 1, 1, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//given
			//when
			_, err := CreateRoomUserDTO(tt.uID, tt.roomID, "ready")
			//then
			assert.Equal(t, tt.err, err)
		})
	}
}
