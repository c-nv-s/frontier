// Code generated by mockery v2.20.2. DO NOT EDIT.

package mocks

import (
	context "context"

	session "github.com/odpf/shield/core/authenticate/session"
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// SessionService is an autogenerated mock type for the SessionService type
type SessionService struct {
	mock.Mock
}

type SessionService_Expecter struct {
	mock *mock.Mock
}

func (_m *SessionService) EXPECT() *SessionService_Expecter {
	return &SessionService_Expecter{mock: &_m.Mock}
}

// Create provides a mock function with given fields: ctx, userID
func (_m *SessionService) Create(ctx context.Context, userID string) (*session.Session, error) {
	ret := _m.Called(ctx, userID)

	var r0 *session.Session
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*session.Session, error)); ok {
		return rf(ctx, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *session.Session); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*session.Session)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SessionService_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type SessionService_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - ctx context.Context
//   - userID string
func (_e *SessionService_Expecter) Create(ctx interface{}, userID interface{}) *SessionService_Create_Call {
	return &SessionService_Create_Call{Call: _e.mock.On("Create", ctx, userID)}
}

func (_c *SessionService_Create_Call) Run(run func(ctx context.Context, userID string)) *SessionService_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *SessionService_Create_Call) Return(_a0 *session.Session, _a1 error) *SessionService_Create_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *SessionService_Create_Call) RunAndReturn(run func(context.Context, string) (*session.Session, error)) *SessionService_Create_Call {
	_c.Call.Return(run)
	return _c
}

// Delete provides a mock function with given fields: ctx, sessionID
func (_m *SessionService) Delete(ctx context.Context, sessionID uuid.UUID) error {
	ret := _m.Called(ctx, sessionID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) error); ok {
		r0 = rf(ctx, sessionID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SessionService_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type SessionService_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - ctx context.Context
//   - sessionID uuid.UUID
func (_e *SessionService_Expecter) Delete(ctx interface{}, sessionID interface{}) *SessionService_Delete_Call {
	return &SessionService_Delete_Call{Call: _e.mock.On("Delete", ctx, sessionID)}
}

func (_c *SessionService_Delete_Call) Run(run func(ctx context.Context, sessionID uuid.UUID)) *SessionService_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *SessionService_Delete_Call) Return(_a0 error) *SessionService_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *SessionService_Delete_Call) RunAndReturn(run func(context.Context, uuid.UUID) error) *SessionService_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// ExtractFromContext provides a mock function with given fields: ctx
func (_m *SessionService) ExtractFromContext(ctx context.Context) (*session.Session, error) {
	ret := _m.Called(ctx)

	var r0 *session.Session
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (*session.Session, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) *session.Session); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*session.Session)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SessionService_ExtractFromContext_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ExtractFromContext'
type SessionService_ExtractFromContext_Call struct {
	*mock.Call
}

// ExtractFromContext is a helper method to define mock.On call
//   - ctx context.Context
func (_e *SessionService_Expecter) ExtractFromContext(ctx interface{}) *SessionService_ExtractFromContext_Call {
	return &SessionService_ExtractFromContext_Call{Call: _e.mock.On("ExtractFromContext", ctx)}
}

func (_c *SessionService_ExtractFromContext_Call) Run(run func(ctx context.Context)) *SessionService_ExtractFromContext_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *SessionService_ExtractFromContext_Call) Return(_a0 *session.Session, _a1 error) *SessionService_ExtractFromContext_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *SessionService_ExtractFromContext_Call) RunAndReturn(run func(context.Context) (*session.Session, error)) *SessionService_ExtractFromContext_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewSessionService interface {
	mock.TestingT
	Cleanup(func())
}

// NewSessionService creates a new instance of SessionService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewSessionService(t mockConstructorTestingTNewSessionService) *SessionService {
	mock := &SessionService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}