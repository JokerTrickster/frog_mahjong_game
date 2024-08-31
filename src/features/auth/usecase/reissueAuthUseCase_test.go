package usecase

import (
	"context"
	"main/features/auth/model/interface/mocks"
	"main/features/auth/model/request"
	"main/features/auth/model/response"
	"main/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"gopkg.in/go-playground/assert.v1"
)

// TestReissueAuthUseCase_Reissue 함수는 ReissueAuthUseCase 의 Reissue 메서드를 테스트합니다.
// 테이블 기반 테스트를 사용하여 여러 시나리오를 테스트합니다.
// given-when-then 패턴을 사용하여 테스트를 작성합니다.
// 매개변수 : req *request.ReqReissue
// 반환값 : response.ResReissue, error
// 테스트 케이스:
// 1. 액세스 토큰과 리프레시 토큰 재발급 성공
// 2. 액세스 토큰 검증 실패
// 테스트 경로: src/features/auth/usecase/reissueAuthUseCase_test.go

func TestReissueAuthUseCase_Reissue(t *testing.T) {
	tests := []struct {
		name    string
		req     *request.ReqReissue
		want    response.ResReissue
		wantErr error
	}{
		{
			name: "액세스 토큰과 리프레시 토큰 재발급 성공",
			want: response.ResReissue{
				AccessToken:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOjEsImVtYWlsIjoic2VjcmV0QG1haWwuY29tIn0.8LWJv0Xy7jgj4q7v",
				RefreshToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOjEsImVtYWlsIjoic2VjcmV0QG1haWwuY29tIn0.8LWJv0Xy7jgj4q7v",
			},
			wantErr: nil,
		},
		{
			name:    "액세스 토큰 검증 실패",
			want:    response.ResReissue{},
			wantErr: utils.ErrorMsg(context.TODO(), utils.ErrBadParameter, utils.Trace(), "no access token in header", utils.ErrFromClient),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			accessTkn, _, refreshTkn, _, _ := utils.GenerateToken("ryan@gmail.com", 1)
			tt.req = &request.ReqReissue{
				AccessToken:  accessTkn,
				RefreshToken: refreshTkn,
			}
			mockReissueAuthRepository := new(mocks.IReissueAuthRepository)
			mockReissueAuthRepository.On("DeleteToken", mock.Anything, mock.Anything).Return(tt.wantErr) //mock
			mockReissueAuthRepository.On("SaveToken", mock.Anything, mock.Anything).Return(tt.wantErr)   //mock
			us := NewReissueAuthUseCase(mockReissueAuthRepository, 8*time.Second)

			// when
			_, err := us.Reissue(context.Background(), tt.req)
			// then
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
