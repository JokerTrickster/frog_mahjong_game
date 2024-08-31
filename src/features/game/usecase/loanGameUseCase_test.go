package usecase

// TestLoanGameUseCase_Loan 함수는 LoanGameUseCase 의 Loan 메서드를 테스트합니다.
// 테이블 기반 테스트를 사용하여 여러 시나리오를 테스트합니다.
// given-when-then 패턴을 사용하여 테스트를 작성합니다.
// 테스트 케이스:
// - loan 을 성공하는 경우
// - loan을 실패하는 경우 (카드가 없는 경우)
// 테스트 경로: src/features/game/usecase/loanGameUseCase_test.go

import (
	"context"
	"main/features/game/model/interface/mocks"
	"main/features/game/model/request"
	"main/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"gopkg.in/go-playground/assert.v1"
)

func TestLoanGameUseCase_Loan(t *testing.T) {
	tests := []struct {
		name    string
		userID  uint
		req     *request.ReqLoan
		wantErr error
	}{
		{
			name:   "loan success",
			userID: 2,
			req: &request.ReqLoan{
				RoomID:     1,
				LoanUserID: 1,
				LoanCardID: 1,
			},
			wantErr: nil,
		},
		{
			name:   "loan fail",
			userID: 1,
			req: &request.ReqLoan{
				RoomID:     1,
				LoanUserID: 1,
				LoanCardID: 2,
			},
			wantErr: utils.ErrorMsg(context.TODO(), utils.ErrInternalDB, utils.Trace(), "Internal DB Error", utils.ErrFromMysqlDB),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			mockLoanGameRepository := new(mocks.ILoanGameRepository)
			mockLoanGameRepository.On("CheckLoan", mock.Anything, mock.Anything).Return(tt.wantErr)
			mockLoanGameRepository.On("Loan", mock.Anything, mock.Anything).Return(tt.wantErr)
			mockLoanGameRepository.On("UpdateRoomUserCardCount", mock.Anything, mock.Anything, mock.Anything).Return(tt.wantErr)
			uc := NewLoanGameUseCase(mockLoanGameRepository, 8*time.Second)

			// when
			err := uc.Loan(context.Background(), tt.userID, tt.req)

			// then
			assert.Equal(t, err, tt.wantErr)
		})
	}
}
