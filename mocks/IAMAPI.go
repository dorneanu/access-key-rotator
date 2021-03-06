// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	context "context"

	iam "github.com/aws/aws-sdk-go-v2/service/iam"

	mock "github.com/stretchr/testify/mock"
)

// IAMAPI is an autogenerated mock type for the IAMAPI type
type IAMAPI struct {
	mock.Mock
}

// CreateAccessKey provides a mock function with given fields: ctx, params, optFns
func (_m *IAMAPI) CreateAccessKey(ctx context.Context, params *iam.CreateAccessKeyInput, optFns ...func(*iam.Options)) (*iam.CreateAccessKeyOutput, error) {
	_va := make([]interface{}, len(optFns))
	for _i := range optFns {
		_va[_i] = optFns[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, params)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *iam.CreateAccessKeyOutput
	if rf, ok := ret.Get(0).(func(context.Context, *iam.CreateAccessKeyInput, ...func(*iam.Options)) *iam.CreateAccessKeyOutput); ok {
		r0 = rf(ctx, params, optFns...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*iam.CreateAccessKeyOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *iam.CreateAccessKeyInput, ...func(*iam.Options)) error); ok {
		r1 = rf(ctx, params, optFns...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteAccessKey provides a mock function with given fields: ctx, params, optFns
func (_m *IAMAPI) DeleteAccessKey(ctx context.Context, params *iam.DeleteAccessKeyInput, optFns ...func(*iam.Options)) (*iam.DeleteAccessKeyOutput, error) {
	_va := make([]interface{}, len(optFns))
	for _i := range optFns {
		_va[_i] = optFns[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, params)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *iam.DeleteAccessKeyOutput
	if rf, ok := ret.Get(0).(func(context.Context, *iam.DeleteAccessKeyInput, ...func(*iam.Options)) *iam.DeleteAccessKeyOutput); ok {
		r0 = rf(ctx, params, optFns...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*iam.DeleteAccessKeyOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *iam.DeleteAccessKeyInput, ...func(*iam.Options)) error); ok {
		r1 = rf(ctx, params, optFns...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListAccessKeys provides a mock function with given fields: ctx, params, optFns
func (_m *IAMAPI) ListAccessKeys(ctx context.Context, params *iam.ListAccessKeysInput, optFns ...func(*iam.Options)) (*iam.ListAccessKeysOutput, error) {
	_va := make([]interface{}, len(optFns))
	for _i := range optFns {
		_va[_i] = optFns[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, params)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *iam.ListAccessKeysOutput
	if rf, ok := ret.Get(0).(func(context.Context, *iam.ListAccessKeysInput, ...func(*iam.Options)) *iam.ListAccessKeysOutput); ok {
		r0 = rf(ctx, params, optFns...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*iam.ListAccessKeysOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *iam.ListAccessKeysInput, ...func(*iam.Options)) error); ok {
		r1 = rf(ctx, params, optFns...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
