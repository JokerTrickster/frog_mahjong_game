package usecase

// createRoomUseCase.go test code write
// Path: src/features/rooms/usecase/createRoomUseCase_test.go

import (
	"context"
	"main/features/rooms/model/interface/mocks"
	"main/features/rooms/model/request"
	"main/features/rooms/model/response"
	"main/utils"
	"testing"
	"time"

	_errors "main/features/rooms/model/errors"

	"github.com/stretchr/testify/mock"
	"gopkg.in/go-playground/assert.v1"
)

func TestCreateRoomsUseCase_Create(t *testing.T) {

	// CreateRoomUseCase.Create test code write
	// Path: src/features/rooms/usecase/createRoomUseCase_test.go
	tests := []struct {
		name string
		req  request.ReqCreate
		res  response.ResCreateRoom
		err  error
	}{
		{"success1", request.ReqCreate{Name: "test", MaxCount: 4, MinCount: 2, Password: "test"}, response.ResCreateRoom{RoomID: 5}, nil},
		{"fail1", request.ReqCreate{Name: "test", MaxCount: 1, MinCount: 2, Password: ""}, response.ResCreateRoom{RoomID: 0}, utils.ErrorMsg(context.TODO(), utils.ErrUserNotFound, utils.Trace(), _errors.ErrBadRequest.Error(), utils.ErrFromClient)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//given
			mockCreateRoomRepository := new(mocks.ICreateRoomsRepository)
			if tt.err == nil {
				mockCreateRoomRepository.On("InsertOneRoom", mock.Anything, mock.Anything).Return(5, nil) //mock
			}

			mockCreateRoomRepository.On("InsertOneRoom", mock.Anything, mock.Anything).Return(0, tt.err) //mock
			mockCreateRoomRepository.On("InsertOneRoomUser", mock.Anything, mock.Anything).Return(nil)   //mock
			mockCreateRoomRepository.On("FindOneAndUpdateUser", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			us := NewCreateRoomsUseCase(mockCreateRoomRepository, 8*time.Second)
			//when
			res, err := us.Create(context.TODO(), 1, "ryan@gmail.com", &tt.req)
			//then
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.res, res)
		})
	}
}

// CreateRoomDTO test code write
// Path: src/features/rooms/usecase/createRoomUseCase_test.go
func TestCreateRoomDTO(t *testing.T) {
	tests := []struct {
		name string
		req  request.ReqCreate
		uID  uint
		err  error
	}{
		{"success1", request.ReqCreate{Name: "test", MaxCount: 4, MinCount: 2, Password: ""}, 1, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//given
			//when
			_, err := CreateRoomDTO(context.TODO(), &tt.req, tt.uID)
			//then
			assert.Equal(t, tt.err, err)
		})
	}
}

// CreateRoomUserDTO test code write
// Path: src/features/rooms/usecase/createRoomUseCase_test.go
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
