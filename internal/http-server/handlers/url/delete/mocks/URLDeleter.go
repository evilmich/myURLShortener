// Code generated by mockery v2.28.2. DO NOT EDIT.

package mocks

import (
	delete "myURLShortener/internal/http-server/handlers/url/delete"

	mock "github.com/stretchr/testify/mock"
)

// URLDeleter is an autogenerated mock type for the URLDeleter type
type URLDeleter struct {
	mock.Mock
}

// DeleteAliasByURL provides a mock function with given fields: url
func (_m *URLDeleter) DeleteAliasByURL(url string) ([]*delete.AliasData, error) {
	ret := _m.Called(url)

	var r0 []*delete.AliasData
	var r1 error
	if rf, ok := ret.Get(0).(func(string) ([]*delete.AliasData, error)); ok {
		return rf(url)
	}
	if rf, ok := ret.Get(0).(func(string) []*delete.AliasData); ok {
		r0 = rf(url)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*delete.AliasData)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(url)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteByAliasAndURL provides a mock function with given fields: alias, url
func (_m *URLDeleter) DeleteByAliasAndURL(alias string, url string) error {
	ret := _m.Called(alias, url)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(alias, url)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteURLByAlias provides a mock function with given fields: alias
func (_m *URLDeleter) DeleteURLByAlias(alias string) error {
	ret := _m.Called(alias)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(alias)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAliasAndURL provides a mock function with given fields: alias, url
func (_m *URLDeleter) GetAliasAndURL(alias string, url string) (string, string, error) {
	ret := _m.Called(alias, url)

	var r0 string
	var r1 string
	var r2 error
	if rf, ok := ret.Get(0).(func(string, string) (string, string, error)); ok {
		return rf(alias, url)
	}
	if rf, ok := ret.Get(0).(func(string, string) string); ok {
		r0 = rf(alias, url)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string, string) string); ok {
		r1 = rf(alias, url)
	} else {
		r1 = ret.Get(1).(string)
	}

	if rf, ok := ret.Get(2).(func(string, string) error); ok {
		r2 = rf(alias, url)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetURL provides a mock function with given fields: alias
func (_m *URLDeleter) GetURL(alias string) (string, error) {
	ret := _m.Called(alias)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (string, error)); ok {
		return rf(alias)
	}
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(alias)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(alias)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewURLDeleter interface {
	mock.TestingT
	Cleanup(func())
}

// NewURLDeleter creates a new instance of URLDeleter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewURLDeleter(t mockConstructorTestingTNewURLDeleter) *URLDeleter {
	mock := &URLDeleter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
