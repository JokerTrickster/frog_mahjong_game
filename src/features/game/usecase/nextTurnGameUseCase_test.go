package usecase

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

// TestNextTurnGameUseCase_NextTurn 함수는 NextTurnGameUseCase 의 NextTurn 메서드를 테스트합니다.
// 테이블 기반 테스트를 사용하여 여러 시나리오를 테스트합니다.
// given-when-then 패턴을 사용하여 테스트를 작성합니다.
// 테스트 케이스:
// -  해당 턴 넘버를 가진 room user가 play_wait인지 확인 후 플레이 상태를 play로 변경 성공한다.
// -  해당 턴 넘버를 가진 room user가 play_wait인지 확인 후 플레이 상태를 play로 변경 실패한다.
// 테스트 경로: src/features/game/usecase/nextTurnGameUseCase_test.go

func TestNextTurnGameUseCase_NextTurn(t *testing.T) {
	tests := []struct {
		name    string
		req     *request.ReqNextTurn
		wantErr error
	}{
		{
			name: "player next turn success",
			req: &request.ReqNextTurn{
				RoomID:      1,
				UserID:      1,
				TurnNumber:  1,
				PlayerState: "play",
			},
			wantErr: nil,
		},
		{
			name: "player next turn fail",
			req: &request.ReqNextTurn{
				RoomID:      1,
				UserID:      1,
				TurnNumber:  1,
				PlayerState: "play",
			},
			wantErr: utils.ErrorMsg(context.TODO(), utils.ErrInternalDB, utils.Trace(), "Internal DB Error", utils.ErrFromMysqlDB),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			mockNextTurnGameRepository := new(mocks.INextTurnGameRepository)
			mockNextTurnGameRepository.On("UpdatePlayerNextTurn", mock.Anything, tt.req).Return(tt.wantErr)
			uc := NewNextTurnGameUseCase(mockNextTurnGameRepository, 8*time.Second)

			// when
			err := uc.NextTurn(context.Background(), tt.req)

			// then
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
