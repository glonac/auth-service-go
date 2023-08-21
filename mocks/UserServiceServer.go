// Code generated by mockery v2.32.0. DO NOT EDIT.

package mocks

import (
	userGrpc "auth-service/pkg/grpc/userGrpc"
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// UserServiceServer is an autogenerated mock type for the UserServiceServer type
type UserServiceServer struct {
	mock.Mock
}

// GetConsumerByClientAccountId provides a mock function with given fields: _a0, _a1
func (_m *UserServiceServer) GetConsumerByClientAccountId(_a0 context.Context, _a1 *userGrpc.GetConsumerByClientAccountId_Request) (*userGrpc.GetConsumerByClientAccountId_Response, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *userGrpc.GetConsumerByClientAccountId_Response
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *userGrpc.GetConsumerByClientAccountId_Request) (*userGrpc.GetConsumerByClientAccountId_Response, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *userGrpc.GetConsumerByClientAccountId_Request) *userGrpc.GetConsumerByClientAccountId_Response); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*userGrpc.GetConsumerByClientAccountId_Response)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *userGrpc.GetConsumerByClientAccountId_Request) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSellerByClientAccountId provides a mock function with given fields: _a0, _a1
func (_m *UserServiceServer) GetSellerByClientAccountId(_a0 context.Context, _a1 *userGrpc.GetSellerByClientAccountId_Request) (*userGrpc.GetSellerByClientAccountId_Response, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *userGrpc.GetSellerByClientAccountId_Response
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *userGrpc.GetSellerByClientAccountId_Request) (*userGrpc.GetSellerByClientAccountId_Response, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *userGrpc.GetSellerByClientAccountId_Request) *userGrpc.GetSellerByClientAccountId_Response); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*userGrpc.GetSellerByClientAccountId_Response)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *userGrpc.GetSellerByClientAccountId_Request) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSellersByTags provides a mock function with given fields: _a0, _a1
func (_m *UserServiceServer) GetSellersByTags(_a0 context.Context, _a1 *userGrpc.GetSellersByTags_Request) (*userGrpc.GetSellersByTags_Response, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *userGrpc.GetSellersByTags_Response
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *userGrpc.GetSellersByTags_Request) (*userGrpc.GetSellersByTags_Response, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *userGrpc.GetSellersByTags_Request) *userGrpc.GetSellersByTags_Response); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*userGrpc.GetSellersByTags_Response)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *userGrpc.GetSellersByTags_Request) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserByClientAccountId provides a mock function with given fields: _a0, _a1
func (_m *UserServiceServer) GetUserByClientAccountId(_a0 context.Context, _a1 *userGrpc.GetUserByClientAccountId_Request) (*userGrpc.GetUserByClientAccountId_Response, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *userGrpc.GetUserByClientAccountId_Response
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *userGrpc.GetUserByClientAccountId_Request) (*userGrpc.GetUserByClientAccountId_Response, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *userGrpc.GetUserByClientAccountId_Request) *userGrpc.GetUserByClientAccountId_Response); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*userGrpc.GetUserByClientAccountId_Response)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *userGrpc.GetUserByClientAccountId_Request) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserById provides a mock function with given fields: _a0, _a1
func (_m *UserServiceServer) GetUserById(_a0 context.Context, _a1 *userGrpc.GetUserById_Request) (*userGrpc.GetUserById_Response, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *userGrpc.GetUserById_Response
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *userGrpc.GetUserById_Request) (*userGrpc.GetUserById_Response, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *userGrpc.GetUserById_Request) *userGrpc.GetUserById_Response); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*userGrpc.GetUserById_Response)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *userGrpc.GetUserById_Request) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUsersList provides a mock function with given fields: _a0, _a1
func (_m *UserServiceServer) GetUsersList(_a0 context.Context, _a1 *userGrpc.GetUsersList_Request) (*userGrpc.GetUsersList_Response, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *userGrpc.GetUsersList_Response
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *userGrpc.GetUsersList_Request) (*userGrpc.GetUsersList_Response, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *userGrpc.GetUsersList_Request) *userGrpc.GetUsersList_Response); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*userGrpc.GetUsersList_Response)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *userGrpc.GetUsersList_Request) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RegisterConsumer provides a mock function with given fields: _a0, _a1
func (_m *UserServiceServer) RegisterConsumer(_a0 context.Context, _a1 *userGrpc.RegisterConsumer_Request) (*userGrpc.RegisterConsumer_Response, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *userGrpc.RegisterConsumer_Response
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *userGrpc.RegisterConsumer_Request) (*userGrpc.RegisterConsumer_Response, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *userGrpc.RegisterConsumer_Request) *userGrpc.RegisterConsumer_Response); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*userGrpc.RegisterConsumer_Response)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *userGrpc.RegisterConsumer_Request) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RegisterSeller provides a mock function with given fields: _a0, _a1
func (_m *UserServiceServer) RegisterSeller(_a0 context.Context, _a1 *userGrpc.RegisterSeller_Request) (*userGrpc.RegisterSeller_Response, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *userGrpc.RegisterSeller_Response
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *userGrpc.RegisterSeller_Request) (*userGrpc.RegisterSeller_Response, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *userGrpc.RegisterSeller_Request) *userGrpc.RegisterSeller_Response); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*userGrpc.RegisterSeller_Response)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *userGrpc.RegisterSeller_Request) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateConsumer provides a mock function with given fields: _a0, _a1
func (_m *UserServiceServer) UpdateConsumer(_a0 context.Context, _a1 *userGrpc.UpdateConsumer_Request) (*userGrpc.UpdateConsumer_Response, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *userGrpc.UpdateConsumer_Response
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *userGrpc.UpdateConsumer_Request) (*userGrpc.UpdateConsumer_Response, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *userGrpc.UpdateConsumer_Request) *userGrpc.UpdateConsumer_Response); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*userGrpc.UpdateConsumer_Response)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *userGrpc.UpdateConsumer_Request) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateSeller provides a mock function with given fields: _a0, _a1
func (_m *UserServiceServer) UpdateSeller(_a0 context.Context, _a1 *userGrpc.UpdateSeller_Request) (*userGrpc.UpdateSeller_Response, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *userGrpc.UpdateSeller_Response
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *userGrpc.UpdateSeller_Request) (*userGrpc.UpdateSeller_Response, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *userGrpc.UpdateSeller_Request) *userGrpc.UpdateSeller_Response); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*userGrpc.UpdateSeller_Response)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *userGrpc.UpdateSeller_Request) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateSellerSetting provides a mock function with given fields: _a0, _a1
func (_m *UserServiceServer) UpdateSellerSetting(_a0 context.Context, _a1 *userGrpc.UpdateSellerSetting_Request) (*userGrpc.UpdateSellerSetting_Response, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *userGrpc.UpdateSellerSetting_Response
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *userGrpc.UpdateSellerSetting_Request) (*userGrpc.UpdateSellerSetting_Response, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *userGrpc.UpdateSellerSetting_Request) *userGrpc.UpdateSellerSetting_Response); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*userGrpc.UpdateSellerSetting_Response)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *userGrpc.UpdateSellerSetting_Request) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mustEmbedUnimplementedUserServiceServer provides a mock function with given fields:
func (_m *UserServiceServer) mustEmbedUnimplementedUserServiceServer() {
	_m.Called()
}

// NewUserServiceServer creates a new instance of UserServiceServer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserServiceServer(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserServiceServer {
	mock := &UserServiceServer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
