package usecase

// TestLogoutAuthUseCase_LogoutAuth 함수는 LogoutAuthUseCase 의 LogoutAuth 메서드를 테스트합니다.
// 테이블 기반 테스트를 사용하여 여러 시나리오를 테스트합니다.
// given-when-then 패턴을 사용하여 테스트를 작성합니다.
// 테스트 매개변수 :
// 테스트 케이스:
// 1. 로그아웃 성공
// 2. 로그아웃 실패
// 테스트 경로: src/features/auth/usecase/LogoutAuthUseCase_test.go

import (
	"context"
	"main/features/auth/model/interface/mocks"
	"main/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"gopkg.in/go-playground/assert.v1"
)

// TestLogoutAuthUseCase_LogoutAuth 함수는 LogoutAuthUseCase 의 LogoutAuth 메서드를 테스트합니다.
func TestLogoutAuthUseCase_LogoutAuth(t *testing.T) {
	// table driven test
	tests := []struct {
		name string
		uid  uint
		err  error
	}{
		{"로그아웃 성공", 1, nil},
		{"로그아웃 실패", 1, utils.ErrorMsg(context.TODO(), utils.ErrUserNotFound, utils.Trace(), string(utils.ErrBadParameter), utils.ErrFromClient)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//given
			mockLogoutAuthRepository := new(mocks.ILogoutAuthRepository)
			mockLogoutAuthRepository.On("FindOneAndUpdateUser", mock.Anything, mock.Anything).Return(tt.err) //mock
			mockLogoutAuthRepository.On("DeleteToken", mock.Anything, mock.Anything).Return(tt.err)          //mock
			us := NewLogoutAuthUseCase(mockLogoutAuthRepository, 8*time.Second)

			//when
			err := us.Logout(context.TODO(), 1)
			//then
			assert.Equal(t, tt.err, err)
		})
	}
}
