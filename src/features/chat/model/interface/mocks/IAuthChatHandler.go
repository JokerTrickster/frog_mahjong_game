// Code generated by mockery v2.42.2. DO NOT EDIT.

package mocks

import (
	echo "github.com/labstack/echo/v4"
	mock "github.com/stretchr/testify/mock"
)

// IAuthChatHandler is an autogenerated mock type for the IAuthChatHandler type
type IAuthChatHandler struct {
	mock.Mock
}

// Auth provides a mock function with given fields: c
func (_m *IAuthChatHandler) Auth(c echo.Context) error {
	ret := _m.Called(c)

	if len(ret) == 0 {
		panic("no return value specified for Auth")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(echo.Context) error); ok {
		r0 = rf(c)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewIAuthChatHandler creates a new instance of IAuthChatHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIAuthChatHandler(t interface {
	mock.TestingT
	Cleanup(func())
}) *IAuthChatHandler {
	mock := &IAuthChatHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
