// Code generated by mockery v2.42.2. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// ILogoutAuthUseCase is an autogenerated mock type for the ILogoutAuthUseCase type
type ILogoutAuthUseCase struct {
	mock.Mock
}

// Logout provides a mock function with given fields: c, uID
func (_m *ILogoutAuthUseCase) Logout(c context.Context, uID uint) error {
	ret := _m.Called(c, uID)

	if len(ret) == 0 {
		panic("no return value specified for Logout")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint) error); ok {
		r0 = rf(c, uID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewILogoutAuthUseCase creates a new instance of ILogoutAuthUseCase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewILogoutAuthUseCase(t interface {
	mock.TestingT
	Cleanup(func())
}) *ILogoutAuthUseCase {
	mock := &ILogoutAuthUseCase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}