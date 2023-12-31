// Code generated by mockery v2.32.3. DO NOT EDIT.

package mocks

import (
	middleware "github.com/SawitProRecruitment/UserService/middleware"
	mock "github.com/stretchr/testify/mock"
)

// JwtInterface is an autogenerated mock type for the JwtInterface type
type JwtInterface struct {
	mock.Mock
}

// CreateToken provides a mock function with given fields: jwtData, expireInHour
func (_m *JwtInterface) CreateToken(jwtData middleware.UserJwtPayload, expireInHour int) (string, error) {
	ret := _m.Called(jwtData, expireInHour)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(middleware.UserJwtPayload, int) (string, error)); ok {
		return rf(jwtData, expireInHour)
	}
	if rf, ok := ret.Get(0).(func(middleware.UserJwtPayload, int) string); ok {
		r0 = rf(jwtData, expireInHour)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(middleware.UserJwtPayload, int) error); ok {
		r1 = rf(jwtData, expireInHour)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsValid provides a mock function with given fields: tokenString
func (_m *JwtInterface) IsValid(tokenString string) (bool, error) {
	ret := _m.Called(tokenString)

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (bool, error)); ok {
		return rf(tokenString)
	}
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(tokenString)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(tokenString)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ParseToken provides a mock function with given fields: tokenString
func (_m *JwtInterface) ParseToken(tokenString string) (*middleware.JwtParsedPayload, error) {
	ret := _m.Called(tokenString)

	var r0 *middleware.JwtParsedPayload
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*middleware.JwtParsedPayload, error)); ok {
		return rf(tokenString)
	}
	if rf, ok := ret.Get(0).(func(string) *middleware.JwtParsedPayload); ok {
		r0 = rf(tokenString)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*middleware.JwtParsedPayload)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(tokenString)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewJwtInterface creates a new instance of JwtInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewJwtInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *JwtInterface {
	mock := &JwtInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
