// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// SingleHealthChecker is an autogenerated mock type for the SingleHealthChecker type
type SingleHealthChecker struct {
	mock.Mock
}

// Check provides a mock function with given fields:
func (_m *SingleHealthChecker) Check() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewSingleHealthChecker interface {
	mock.TestingT
	Cleanup(func())
}

// NewSingleHealthChecker creates a new instance of SingleHealthChecker. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewSingleHealthChecker(t mockConstructorTestingTNewSingleHealthChecker) *SingleHealthChecker {
	mock := &SingleHealthChecker{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
