package usecase

// TestScoreCalculateGameUseCase_ScoreCalculate 함수는 ScoreCalculateGameUseCase 의 ScoreCalculate 메서드를 테스트합니다.
// 테이블 기반 테스트를 사용하여 여러 시나리오를 테스트합니다.
// given-when-then 패턴을 사용하여 테스트를 작성합니다.
// 요청 Cards는 ScoreCard 구조체를 가지고 있으며, RoomID, CardID, Color, Name, State를 가지고 있습니다.
// 테스트 케이스:
// - 카드 6장을 요청을 받고 점수 계산 성공하는 경우 (dora 카드가 없는 경우)
// - 카드 6장을 요청을 받고 점수 계산 성공하는 경우 (dora 카드가 있는 경우)
// 테스트 경로: src/features/game/usecase/scoreCalculateGameUseCase_test.go

import (
	"context"
	"main/features/game/model/interface/mocks"
	"main/features/game/model/request"
	"main/utils/db/mysql"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"gopkg.in/go-playground/assert.v1"
)

func TestScoreCalculateGameUseCase_ScoreCalculate(t *testing.T) {
	tests := []struct {
		name    string
		userID  uint
		req     *request.ReqScoreCalculate
		dora    mysql.Cards
		wantErr error
	}{
		{
			name:   "success",
			userID: 1,
			dora: mysql.Cards{
				RoomID: 1,
				Name:   "1",
				State:  "dora",
			},
			req: &request.ReqScoreCalculate{
				RoomID: 1,
				Cards: []request.ScoreCard{
					{
						CardID: 1,
						Color:  "red",
						Name:   "1",
						State:  "none",
					},
					{
						CardID: 2,
						Color:  "red",
						Name:   "2",
						State:  "none",
					},
					{
						CardID: 3,
						Color:  "red",
						Name:   "3",
						State:  "none",
					},
					{
						CardID: 4,
						Color:  "red",
						Name:   "4",
						State:  "none",
					},
					{
						CardID: 5,
						Color:  "red",
						Name:   "5",
						State:  "none",
					},
					{
						CardID: 6,
						Color:  "red",
						Name:   "6",
						State:  "none",
					},
				},
			},
			wantErr: nil,
		},
		{
			name:   "success",
			userID: 1,
			dora: mysql.Cards{
				RoomID: 1,
				Name:   "1",
				State:  "dora",
			},
			req: &request.ReqScoreCalculate{
				RoomID: 1,
				Cards: []request.ScoreCard{
					{
						CardID: 1,
						Color:  "red",
						Name:   "1",
						State:  "none",
					},
					{
						CardID: 2,
						Color:  "red",
						Name:   "2",
						State:  "none",
					},
					{
						CardID: 3,
						Color:  "red",
						Name:   "3",
						State:  "none",
					},
					{
						CardID: 4,
						Color:  "red",
						Name:   "4",
						State:  "none",
					},
					{
						CardID: 5,
						Color:  "red",
						Name:   "5",
						State:  "none",
					},
					{
						CardID: 6,
						Color:  "red",
						Name:   "6",
						State:  "none",
					},
				},
			},

			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			mockScoreCalculateRepo := new(mocks.IScoreCalculateGameRepository)
			mockScoreCalculateRepo.On("CheckCardCount", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			mockScoreCalculateRepo.On("GetDoraCard", mock.Anything, mock.Anything).Return(tt.dora, tt.wantErr)

			uc := NewScoreCalculateGameUseCase(mockScoreCalculateRepo, 8*time.Second)

			// when
			_, _, err := uc.ScoreCalculate(context.Background(), tt.userID, tt.req)

			// then
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
