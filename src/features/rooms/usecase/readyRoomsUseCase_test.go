package usecase

// readyRoomsUseCase.go test code write.
// i want to write it as table driven test code.
// Path: src/features/rooms/usecase/readyRoomsUseCase_test.go

import (
	"context"
	"main/features/rooms/model/interface/mocks"
	"main/features/rooms/model/request"
	"main/utils"
	"testing"
	_errors "main/features/rooms/model/errors"
	"time"

	"github.com/stretchr/testify/mock"
	"gopkg.in/go-playground/assert.v1"
)

func TestReadyRoomsUseCase_Ready(t *testing.T) {
	tests := []struct {
		name string
		uID  uint
		req  *request.ReqReady
		err  error
	}{
		// Add test cases here
		{
			name: "Test Case 1 success",
			uID:  1,
			req: &request.ReqReady{
				RoomID:      28,
				PlayerState: "ready",
			},
			err: nil,
		},
		{
			name: "Test Case 2 fail",
			uID:  2,
			req: &request.ReqReady{
				RoomID:      29,
				PlayerState: "test",
			},
			err: utils.ErrorMsg(context.TODO(), utils.ErrBadParameter, utils.Trace(), _errors.ErrPlayerStateFailed.Error(), utils.ErrFromClient),
		},
		// Add more test cases here
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//given
			mockReadyRoomRepository := new(mocks.IReadyRoomsRepository)
			mockReadyRoomRepository.On("FindOneAndUpdateRoomUser", mock.Anything, mock.Anything, mock.Anything).Return(tt.err) //mock
			us := NewReadyRoomsUseCase(mockReadyRoomRepository, 8*time.Second)

			//when
			err := us.Ready(context.TODO(), tt.uID, tt.req)
			//then
			assert.Equal(t, tt.err, err)
		})
	}
}
