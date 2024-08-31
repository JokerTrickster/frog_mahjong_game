package usecase

import (
	"context"
	"main/features/game/model/interface/mocks"
	"main/features/game/model/request"
	"main/utils/db/mysql"

	"testing"

	"github.com/stretchr/testify/mock"
	"gopkg.in/go-playground/assert.v1"
)

// TestDiscardGameUseCase_Discard 함수는 DiscardGameUseCase 의 Discard 메서드를 테스트합니다.
// 테이블 기반 테스트를 사용하여 여러 시나리오를 테스트합니다.
// given-when-then 패턴을 사용하여 테스트를 작성합니다.
// 테스트 케이스:
// - 카드 1개를 요청으로 받고 업데이트 성공하는 경우
// 테스트 경로: src/features/game/usecase/discardGameUseCase_test.go

func TestDiscardGameUseCase_Discard(t *testing.T) {
	tests := []struct {
		name     string
		userID   int
		req      *request.ReqDiscard
		roomUser mysql.RoomUsers
		wantErr  error
	}{
		{
			name:   "카드 1개를 요청으로 받고 업데이트 성공하는 경우",
			userID: 1,
			req: &request.ReqDiscard{
				CardID:    1,
				RoomID:    1,
				UserID:    1,
				CardState: "discard",
			},
			roomUser: mysql.RoomUsers{
				UserID:         1,
				RoomID:         1,
				Score:          0,
				OwnedCardCount: 5,
				PlayerState:    "play",
				TurnNumber:     1,
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			mockDiscardGameRepository := new(mocks.IDiscardGameRepository)
			mockDiscardGameRepository.On("PlayerCheckTurn", mock.Anything, mock.Anything).Return(tt.roomUser, nil)
			mockDiscardGameRepository.On("UpdateCardStateDiscard", mock.Anything, mock.Anything).Return(tt.wantErr)
			mockDiscardGameRepository.On("UpdateRoomUser", mock.Anything, mock.Anything).Return(tt.wantErr)
			DiscardGameUseCase := NewDiscardGameUseCase(mockDiscardGameRepository, 1)

			// when
			err := DiscardGameUseCase.Discard(context.Background(), tt.userID, tt.req)
			// then
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

// TestCreateUpdateRoomUser 함수는 DiscardGameUseCase 의 CreateUpdateRoomUser 메서드를 테스트합니다.
// 테이블 기반 테스트를 사용하여 여러 시나리오를 테스트합니다.
// given-when-then 패턴을 사용하여 테스트를 작성합니다.
// 테스트 케이스:
// - roomUser 의 OwnedCardCount 를 -1 하고 PlayerState 를 play_wait 로 변경하는 경우
// 테스트 경로: src/features/game/usecase/discardGameUseCase_test.go

func TestCreateUpdateRoomUser(t *testing.T) {
	tests := []struct {
		name     string
		roomUser mysql.RoomUsers
		req      *request.ReqDiscard
		want     mysql.RoomUsers
	}{
		{
			name: "roomUser 의 OwnedCardCount 를 -1 하고 PlayerState 를 play_wait 로 변경하는 경우",
			roomUser: mysql.RoomUsers{
				UserID:         1,
				RoomID:         1,
				Score:          0,
				OwnedCardCount: 5,
				PlayerState:    "play",
				TurnNumber:     1,
			},
			req: &request.ReqDiscard{
				CardID:    1,
				RoomID:    1,
				UserID:    1,
				CardState: "discard",
			},
			want: mysql.RoomUsers{
				UserID:         1,
				RoomID:         1,
				Score:          0,
				OwnedCardCount: 4,
				PlayerState:    "play_wait",
				TurnNumber:     1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// when
			got := CreateUpdateRoomUser(tt.roomUser, tt.req)
			// then
			assert.Equal(t, tt.want, got)
		})
	}
}
