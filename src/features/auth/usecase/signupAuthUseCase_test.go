package usecase

// Write code to test the signupAuthUseCase.go file

import (
	"context"
	"main/features/auth/model/interface/mocks"
	"main/features/auth/model/request"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"gopkg.in/go-playground/assert.v1"
)

func TestSignupAuthUseCase_Signup(t *testing.T) {
	// table driven test
	tests := []struct {
		name string
		req  request.ReqSignup
		want error
	}{
		{"success1", request.ReqSignup{Email: "test01@test.com", Password: "1234", Name: "ryan"}, nil},
		{"success2", request.ReqSignup{Email: "test02@test.com", Password: "1234", Name: "test"}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//given
			mockSignupAuthRepository := new(mocks.ISignupAuthRepository)
			mockSignupAuthRepository.On("UserCheckByEmail", mock.Anything, mock.Anything).Return(nil) //mock
			mockSignupAuthRepository.On("InsertOneUser", mock.Anything, mock.Anything).Return(nil)    //mock
			us := NewSignupAuthUseCase(mockSignupAuthRepository, 8*time.Second)

			//when
			got := us.Signup(context.TODO(), &tt.req)
			//then
			assert.Equal(t, tt.want, got)
		})
	}
}
