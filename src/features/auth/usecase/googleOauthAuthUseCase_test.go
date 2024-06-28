package usecase

import (
	"context"
	"main/features/auth/model/interface/mocks"
	"testing"
	"time"

	"gopkg.in/go-playground/assert.v1"
)

// TestGoogleOauthUseCase_GoogleOauth 함수는 GoogleOauthUseCase 의 GoogleOauth 메서드를 테스트합니다.
// 테이블 기반 테스트를 사용하여 여러 시나리오를 테스트합니다.
// given-when-then 패턴을 사용하여 테스트를 작성합니다.
// 매개변수 :
// 반환값 : string, error
// 테스트 케이스:
// 1. 구글 oauth state 생성 성공
// 테스트 경로: src/features/auth/usecase/googleOauthAuthUseCase_test.go

func TestGoogleOauthUseCase_GoogleOauth(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr error
	}{
		{
			name:    "구글 oauth state 생성 성공",
			want:    "sWXhTT0K_qRcqRpMQt2HLQ==",
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			mockGoogleOauthAuthRepository := new(mocks.IGoogleOauthAuthRepository)
			us := NewGoogleOauthAuthUseCase(mockGoogleOauthAuthRepository, 8*time.Second)

			// when
			_, err := us.GoogleOauth(context.Background())

			// then
			if (err != nil) && (tt.wantErr != nil) {
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			}
		})
	}
}
