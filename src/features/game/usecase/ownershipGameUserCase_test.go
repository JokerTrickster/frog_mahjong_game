package usecase

import (
	"context"
	"main/features/game/model/interface/mocks"
	"main/features/game/model/request"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"gopkg.in/go-playground/assert.v1"
)

// TestOwnershipGameUseCase_Dora 함수는 OwnershipGameUseCase 의 Ownership 메서드를 테스트합니다.
// 테이블 기반 테스트를 사용하여 여러 시나리오를 테스트합니다.
// given-when-then 패턴을 사용하여 테스트를 작성합니다.
// 테스트 케이스:
// - 카드 5개를 요청으로 받고 업데이트 성공하는 경우
// 테스트 경로: src/features/game/usecase/ownershipGameUseCase_test.go

func TestOwnershipGameUseCase_Ownership(t *testing.T) {
	tests := []struct {
		name    string
		userID  int
		req     *request.ReqOwnership
		wantErr error
	}{
		{
			name:   "success",
			userID: 1,
			req: &request.ReqOwnership{
				Cards: []request.Card{
					{
						RoomID: 1,
						CardID: 1,
						State:  "none",
						UserID: 1,
					},
					{
						RoomID: 1,
						CardID: 2,
						State:  "none",
						UserID: 1,
					},
					{
						RoomID: 1,
						CardID: 3,
						State:  "none",
						UserID: 1,
					},
					{
						RoomID: 1,
						CardID: 4,
						State:  "none",
						UserID: 1,
					},
					{
						RoomID: 1,
						CardID: 5,
						State:  "none",
						UserID: 1,
					},
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			mockOwnershipGameRepository := new(mocks.IOwnershipGameRepository)
			mockOwnershipGameRepository.On("UpdateCardState", mock.Anything, mock.Anything).Return(tt.wantErr)
			mockOwnershipGameRepository.On("UpdateRoomUserCardCount", mock.Anything, mock.Anything).Return(tt.wantErr)
			uc := NewOwnershipGameUseCase(mockOwnershipGameRepository, 8*time.Second)

			// when
			err := uc.Ownership(context.Background(), tt.req)

			// then
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
