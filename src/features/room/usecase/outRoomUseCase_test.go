package usecase

import (
	"context"
	"testing"
	"time"

	_errors "main/features/room/model/errors"
	"main/features/room/model/interface/mocks"
	"main/features/room/model/request"
	"main/utils"

	"github.com/stretchr/testify/mock"
	"gopkg.in/go-playground/assert.v1"
)

func TestOutRoomUseCase_Out(t *testing.T) {
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
			mockOutRoomRepository := new(mocks.IOutRoomRepository)
			mockOutRoomRepository.On("FindOneAndDeleteRoomUser", mock.Anything, mock.Anything, mock.Anything).Return(tc.err) //mock
			mockOutRoomRepository.On("FindOneAndUpdateRoom", mock.Anything, mock.Anything).Return(nil)
			mockOutRoomRepository.On("FindOneAndUpdateUser", mock.Anything, mock.Anything).Return(nil)

			us := NewOutRoomUseCase(mockOutRoomRepository, 8*time.Second)
			//when
			err := us.Out(context.TODO(), tc.uID, tc.req)

			//then
			assert.Equal(t, tc.err, err)
		})
	}
}
