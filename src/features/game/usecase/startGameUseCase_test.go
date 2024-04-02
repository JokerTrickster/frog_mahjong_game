package usecase

// TestStartGameUseCase_Start 함수는 StartGameUseCase의 Start 메서드를 테스트합니다.
// 테이블 기반 테스트를 사용하여 여러 시나리오를 테스트합니다.
// 테스트 케이스:
// - 게임이 시작되는 경우
// - 유저가 모두 준비하지 않은 경우
// - 방장이 시작을 요청하지 않은 경우
// 테스트 경로: src/features/game/usecase/startGameUseCase_test.go

import (
	"context"
	_errors "main/features/game/model/errors"
	"testing"
	"time"

	"main/features/game/model/interface/mocks"
	"main/features/game/model/request"
	"main/utils/db/mysql"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestStartGameUseCase_Start(t *testing.T) {

	tests := []struct {
		name        string
		email       string
		req         request.ReqStart
		roomDTO     mysql.Rooms
		roomUserDTO []mysql.RoomUsers
		cardDTO     []mysql.Cards
		err         error
	}{
		{
			name:  "game start success",
			email: "ryan@gmail.com",
			req: request.ReqStart{
				RoomID: 1,
				State:  "play",
			},
			roomDTO: mysql.Rooms{
				CurrentCount: 2,
				MaxCount:     2,
				MinCount:     2,
				Name:         "test",
				Password:     "",
				State:        "wait",
				Owner:        "ryan@gamil.com",
			},
			roomUserDTO: []mysql.RoomUsers{
				{
					UserID:         1,
					RoomID:         1,
					PlayerState:    "ready",
					Score:          0,
					OwnedCardCount: 0,
				},
				{
					UserID:         2,
					RoomID:         1,
					PlayerState:    "ready",
					Score:          0,
					OwnedCardCount: 0,
				},
			},
			cardDTO: []mysql.Cards{
				{
					RoomID: 1,
					Name:   "1",
					Color:  "red",
					State:  "none",
				},
				{
					RoomID: 1,
					Name:   "2",
					Color:  "red",
					State:  "none",
				},
				{
					RoomID: 1,
					Name:   "3",
					Color:  "red",
					State:  "none",
				},
				{
					RoomID: 1,
					Name:   "4",
					Color:  "red",
					State:  "none",
				},
				{
					RoomID: 1,
					Name:   "5",
					Color:  "red",
					State:  "none",
				},
				{
					RoomID: 1,
					Name:   "6",
					Color:  "red",
					State:  "none",
				},
				{
					RoomID: 1,
					Name:   "7",
					Color:  "red",
					State:  "none",
				},
			}, err: nil,
		},
		{
			name:  "owner did not request start",
			email: "joker@gmail.com",
			req: request.ReqStart{
				RoomID: 1,
				State:  "play",
			},
			roomDTO: mysql.Rooms{
				CurrentCount: 2,
				MaxCount:     2,
				MinCount:     2,
				Name:         "test",
				Password:     "",
				State:        "wait",
				Owner:        "ryan@gmail.com",
			},
			roomUserDTO: []mysql.RoomUsers{
				{
					UserID:         1,
					RoomID:         1,
					PlayerState:    "ready",
					Score:          0,
					OwnedCardCount: 0,
				},
				{
					UserID:         2,
					RoomID:         1,
					PlayerState:    "ready",
					Score:          0,
					OwnedCardCount: 0,
				},
			},
			cardDTO: nil,
			err:     _errors.ErrNotOwner,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			mockStartGameRepository := new(mocks.IStartGameRepository)
			mockStartGameRepository.On("CheckOwner", mock.Anything, mock.Anything, mock.Anything).Return(tt.err)
			mockStartGameRepository.On("CheckReady", mock.Anything, mock.Anything).Return(tt.roomUserDTO, nil)
			mockStartGameRepository.On("UpdateRoomUser", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			mockStartGameRepository.On("UpdateRoom", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			mockStartGameRepository.On("CreateCards", mock.Anything, mock.Anything, mock.Anything).Return(nil)

			uc := NewStartGameUseCase(mockStartGameRepository, 8*time.Second)

			// When
			err := uc.Start(context.Background(), tt.email, &tt.req)

			// Then
			assert.Equal(t, tt.err, err)
		})
	}
}

// CheckRoomUsersReady 함수에 대한 테스트 코드 작성합니다.
// 테스트 케이스:
// - 모든 유저가 준비한 경우
// - 모든 유저가 준비하지 않은 경우
// 테스트 경로: src/features/game/usecase/startGameUseCase_test.go

func TestCheckRoomUsersReady(t *testing.T) {
	tests := []struct {
		name     string
		roomUser []mysql.RoomUsers
		want     bool
	}{
		{
			name: "all users are ready",
			roomUser: []mysql.RoomUsers{
				{
					UserID:         1,
					RoomID:         1,
					PlayerState:    "ready",
					Score:          0,
					OwnedCardCount: 0,
				},
				{
					UserID:         2,
					RoomID:         1,
					PlayerState:    "ready",
					Score:          0,
					OwnedCardCount: 0,
				},
			},
			want: true,
		},
		{
			name: "not all users are ready",
			roomUser: []mysql.RoomUsers{
				{
					UserID:         1,
					RoomID:         1,
					PlayerState:    "ready",
					Score:          0,
					OwnedCardCount: 0,
				},
				{
					UserID:         2,
					RoomID:         1,
					PlayerState:    "wait",
					Score:          0,
					OwnedCardCount: 0,
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// When
			got := CheckRoomUsersReady(tt.roomUser)

			// Then
			assert.Equal(t, tt.want, got)
		})
	}
}
