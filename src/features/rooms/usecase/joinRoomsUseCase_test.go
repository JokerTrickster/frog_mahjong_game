package usecase
/*
TODO 트랜잭션 처리
// create join room usecase test code write
// table driven test
// Path: src/features/rooms/usecase/joinRoomUseCase_test.go

import (
	"context"
	"testing"
	"time"

	_errors "main/features/rooms/model/errors"
	"main/features/rooms/model/interface/mocks"
	"main/features/rooms/model/request"
	"main/features/rooms/model/response"
	"main/utils"
	"main/utils/db/mysql"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestJoinRoomsUseCase_Join(t *testing.T) {
	tests := []struct {
		name        string
		req         request.ReqJoin
		res         response.ResJoinRoom
		roomDTO     mysql.Rooms
		roomUserDTO mysql.RoomUsers
		err         error
	}{
		{
			name: "success",
			req: request.ReqJoin{
				RoomID: 1,
			},
			res: response.ResJoinRoom{
				RoomID: 1,
			},
			roomDTO: mysql.Rooms{
				CurrentCount: 1,
				MaxCount:     2,
				MinCount:     2,
				Name:         "test",
				Password:     "",
				State:        "wait",
				OwnerID:      1,
			},
			roomUserDTO: mysql.RoomUsers{
				UserID:         1,
				RoomID:         18,
				PlayerState:    "ready",
				Score:          0,
				OwnedCardCount: 0,
			},
			err: nil,
		},
		{
			name: "fail",
			req: request.ReqJoin{
				RoomID: 1,
			},
			res: response.ResJoinRoom{
				RoomID: 0,
			},
			roomDTO: mysql.Rooms{
				CurrentCount: 2,
				MaxCount:     2,
				MinCount:     2,
				Name:         "test",
				Password:     "",
				State:        "wait",
				OwnerID:      1,
			},
			roomUserDTO: mysql.RoomUsers{
				UserID:         1,
				RoomID:         18,
				PlayerState:    "ready",
				Score:          0,
				OwnedCardCount: 0,
			},
			err: utils.ErrorMsg(context.TODO(), utils.ErrRoomFull, utils.Trace(), _errors.ErrRoomFull.Error(), utils.ErrFromClient),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			mockJoinRoomRepository := new(mocks.IJoinRoomsRepository)
			mockJoinRoomRepository.On("FindOneRoom", mock.Anything, mock.Anything).Return(tt.roomDTO, tt.err)
			if tt.err == nil {
				mockJoinRoomRepository.On("InsertOneRoomUser", mock.Anything, mock.Anything).Return(nil)
				mockJoinRoomRepository.On("FindOneAndUpdateRoom", mock.Anything, mock.Anything).Return(nil)
				mockJoinRoomRepository.On("FindOneAndUpdateUser", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			}
			us := NewJoinRoomsUseCase(mockJoinRoomRepository, 8*time.Second)
			// When
			res, err := us.Join(context.Background(), 1, "ryan@gamil.com", &tt.req)
			// Then
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.res, res)
		})
	}
}

*/