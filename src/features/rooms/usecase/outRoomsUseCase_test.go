package usecase

import (
	"context"
	"testing"
	"time"

	_errors "main/features/rooms/model/errors"
	"main/features/rooms/model/interface/mocks"
	"main/features/rooms/model/request"
	"main/utils"

	"github.com/stretchr/testify/mock"
	"gopkg.in/go-playground/assert.v1"
)

func TestOutRoomsUseCase_Out(t *testing.T) {
	testCases := []struct {
		name string
		uID  uint
		req  *request.ReqOut
		err  error
	}{
		{
			name: "Test Case 1 success",
			uID:  1,
			req:  &request.ReqOut{},
			err:  nil,
		},
		{
			name: "Test Case 2 Failed",
			uID:  2,
			req:  &request.ReqOut{},
			err:  utils.ErrorMsg(context.TODO(), utils.ErrBadParameter, utils.Trace(), _errors.ErrRoomUserNotFound.Error(), utils.ErrFromClient),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//given
			mockOutRoomRepository := new(mocks.IOutRoomsRepository)
			mockOutRoomRepository.On("FindOneAndDeleteRoomUser", mock.Anything, mock.Anything, mock.Anything).Return(tc.err) //mock
			mockOutRoomRepository.On("FindOneAndUpdateRoom", mock.Anything, mock.Anything).Return(nil)
			mockOutRoomRepository.On("FindOneAndUpdateUser", mock.Anything, mock.Anything).Return(nil)

			us := NewOutRoomsUseCase(mockOutRoomRepository, 8*time.Second)
			//when
			err := us.Out(context.TODO(), tc.uID, tc.req)

			//then
			assert.Equal(t, tc.err, err)
		})
	}
}
