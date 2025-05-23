// Code generated by MockGen. DO NOT EDIT.
// Source: services.go
//
// Generated by this command:
//
//	mockgen -source=services.go -destination=../../mocks/services/auth_service.go -package=mockservices -exclude_interfaces=UsersService
//

// Package mockservices is a generated GoMock package.
package mockservices

import (
	context "context"
	reflect "reflect"
	time "time"

	entities "github.com/DKhorkov/hmtm-sso/internal/entities"
	gomock "go.uber.org/mock/gomock"
)

// MockAuthService is a mock of AuthService interface.
type MockAuthService struct {
	ctrl     *gomock.Controller
	recorder *MockAuthServiceMockRecorder
	isgomock struct{}
}

// MockAuthServiceMockRecorder is the mock recorder for MockAuthService.
type MockAuthServiceMockRecorder struct {
	mock *MockAuthService
}

// NewMockAuthService creates a new mock instance.
func NewMockAuthService(ctrl *gomock.Controller) *MockAuthService {
	mock := &MockAuthService{ctrl: ctrl}
	mock.recorder = &MockAuthServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthService) EXPECT() *MockAuthServiceMockRecorder {
	return m.recorder
}

// ChangePassword mocks base method.
func (m *MockAuthService) ChangePassword(ctx context.Context, userID uint64, newPassword string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangePassword", ctx, userID, newPassword)
	ret0, _ := ret[0].(error)
	return ret0
}

// ChangePassword indicates an expected call of ChangePassword.
func (mr *MockAuthServiceMockRecorder) ChangePassword(ctx, userID, newPassword any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangePassword", reflect.TypeOf((*MockAuthService)(nil).ChangePassword), ctx, userID, newPassword)
}

// CreateRefreshToken mocks base method.
func (m *MockAuthService) CreateRefreshToken(ctx context.Context, userID uint64, refreshToken string, ttl time.Duration) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRefreshToken", ctx, userID, refreshToken, ttl)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateRefreshToken indicates an expected call of CreateRefreshToken.
func (mr *MockAuthServiceMockRecorder) CreateRefreshToken(ctx, userID, refreshToken, ttl any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRefreshToken", reflect.TypeOf((*MockAuthService)(nil).CreateRefreshToken), ctx, userID, refreshToken, ttl)
}

// ExpireRefreshToken mocks base method.
func (m *MockAuthService) ExpireRefreshToken(ctx context.Context, refreshToken string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExpireRefreshToken", ctx, refreshToken)
	ret0, _ := ret[0].(error)
	return ret0
}

// ExpireRefreshToken indicates an expected call of ExpireRefreshToken.
func (mr *MockAuthServiceMockRecorder) ExpireRefreshToken(ctx, refreshToken any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExpireRefreshToken", reflect.TypeOf((*MockAuthService)(nil).ExpireRefreshToken), ctx, refreshToken)
}

// ForgetPassword mocks base method.
func (m *MockAuthService) ForgetPassword(ctx context.Context, userID uint64, newPassword string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ForgetPassword", ctx, userID, newPassword)
	ret0, _ := ret[0].(error)
	return ret0
}

// ForgetPassword indicates an expected call of ForgetPassword.
func (mr *MockAuthServiceMockRecorder) ForgetPassword(ctx, userID, newPassword any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ForgetPassword", reflect.TypeOf((*MockAuthService)(nil).ForgetPassword), ctx, userID, newPassword)
}

// GetRefreshTokenByUserID mocks base method.
func (m *MockAuthService) GetRefreshTokenByUserID(ctx context.Context, userID uint64) (*entities.RefreshToken, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRefreshTokenByUserID", ctx, userID)
	ret0, _ := ret[0].(*entities.RefreshToken)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRefreshTokenByUserID indicates an expected call of GetRefreshTokenByUserID.
func (mr *MockAuthServiceMockRecorder) GetRefreshTokenByUserID(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRefreshTokenByUserID", reflect.TypeOf((*MockAuthService)(nil).GetRefreshTokenByUserID), ctx, userID)
}

// RegisterUser mocks base method.
func (m *MockAuthService) RegisterUser(ctx context.Context, userData entities.RegisterUserDTO) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterUser", ctx, userData)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RegisterUser indicates an expected call of RegisterUser.
func (mr *MockAuthServiceMockRecorder) RegisterUser(ctx, userData any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterUser", reflect.TypeOf((*MockAuthService)(nil).RegisterUser), ctx, userData)
}

// VerifyUserEmail mocks base method.
func (m *MockAuthService) VerifyUserEmail(ctx context.Context, userID uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyUserEmail", ctx, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// VerifyUserEmail indicates an expected call of VerifyUserEmail.
func (mr *MockAuthServiceMockRecorder) VerifyUserEmail(ctx, userID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyUserEmail", reflect.TypeOf((*MockAuthService)(nil).VerifyUserEmail), ctx, userID)
}
