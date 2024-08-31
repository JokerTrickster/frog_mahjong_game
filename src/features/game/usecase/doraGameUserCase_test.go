package usecase

import (
	"context"
	_errors "main/features/game/model/errors"
	"main/features/game/model/interface/mocks"
	"main/features/game/model/request"
	"main/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"gopkg.in/go-playground/assert.v1"
)

// TestDoraGameUseCase_Dora 함수는 DoraGameUseCase 의 Dora 메서드를 테스트합니다.
// 테이블 기반 테스트를 사용하여 여러 시나리오를 테스트합니다.
// given-when-then 패턴을 사용하여 테스트를 작성합니다.
// 테스트 케이스:
// - 플레이어 상태가 1인 유저가 도라를 선택한 경우
// - 플레이어 상태가 1이 아닌 유저가 도라를 선택하는 경우
// 테스트 경로: src/features/game/usecase/doraGameUseCase_test.go

func TestDoraGameUseCase_Dora(t *testing.T) {
	tests := []struct {
		name    string
		userID  int
		req     *request.ReqDora
		wantErr error
	}{
		{
			name:   "player state is 1",
			userID: 1,
			req: &request.ReqDora{
				RoomID: 1,
				CardID: 1,
			},
			wantErr: nil,
		},
		{
			name:   "player state is not 1",
			userID: 2,
			req: &request.ReqDora{
				RoomID: 1,
				CardID: 1,
			},
			wantErr: utils.ErrorMsg(context.TODO(), utils.ErrBadParameter, utils.Trace(), _errors.ErrNotFirstPlayer.Error(), utils.ErrFromClient),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			mockDoraGameRepository := new(mocks.IDoraGameRepository)
			mockDoraGameRepository.On("CheckFirstPlayer", mock.Anything, tt.userID, tt.req.RoomID).Return(tt.wantErr)
			mockDoraGameRepository.On("UpdateDoraCard", mock.Anything, tt.req).Return(tt.wantErr)
			uc := NewDoraGameUseCase(mockDoraGameRepository, 8*time.Second)

			// when
			err := uc.Dora(context.Background(), tt.userID, tt.req)

			// then
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
