// Code generated by mockery v2.42.2. DO NOT EDIT.

package mocks

import (
	echo "github.com/labstack/echo/v4"
	mock "github.com/stretchr/testify/mock"
)

// IScoreCalculateGameHandler is an autogenerated mock type for the IScoreCalculateGameHandler type
type IScoreCalculateGameHandler struct {
	mock.Mock
}

// ScoreCalculate provides a mock function with given fields: c
func (_m *IScoreCalculateGameHandler) ScoreCalculate(c echo.Context) error {
	ret := _m.Called(c)

	if len(ret) == 0 {
		panic("no return value specified for ScoreCalculate")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(echo.Context) error); ok {
		r0 = rf(c)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewIScoreCalculateGameHandler creates a new instance of IScoreCalculateGameHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIScoreCalculateGameHandler(t interface {
	mock.TestingT
	Cleanup(func())
}) *IScoreCalculateGameHandler {
	mock := &IScoreCalculateGameHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}