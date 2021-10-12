// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2021 Intel Corporation

// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	context "context"

	flow "github.com/otcshare/intel-ethernet-operator/pkg/flowconfig/rpc/v1/flow"
	grpc "google.golang.org/grpc"

	mock "github.com/stretchr/testify/mock"
)

// FlowServiceClient is an autogenerated mock type for the FlowServiceClient type
type FlowServiceClient struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, in, opts
func (_m *FlowServiceClient) Create(ctx context.Context, in *flow.RequestFlowCreate, opts ...grpc.CallOption) (*flow.ResponseFlowCreate, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *flow.ResponseFlowCreate
	if rf, ok := ret.Get(0).(func(context.Context, *flow.RequestFlowCreate, ...grpc.CallOption) *flow.ResponseFlowCreate); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.ResponseFlowCreate)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *flow.RequestFlowCreate, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Destroy provides a mock function with given fields: ctx, in, opts
func (_m *FlowServiceClient) Destroy(ctx context.Context, in *flow.RequestFlowofPort, opts ...grpc.CallOption) (*flow.ResponseFlow, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *flow.ResponseFlow
	if rf, ok := ret.Get(0).(func(context.Context, *flow.RequestFlowofPort, ...grpc.CallOption) *flow.ResponseFlow); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.ResponseFlow)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *flow.RequestFlowofPort, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Flush provides a mock function with given fields: ctx, in, opts
func (_m *FlowServiceClient) Flush(ctx context.Context, in *flow.RequestofPort, opts ...grpc.CallOption) (*flow.ResponseFlow, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *flow.ResponseFlow
	if rf, ok := ret.Get(0).(func(context.Context, *flow.RequestofPort, ...grpc.CallOption) *flow.ResponseFlow); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.ResponseFlow)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *flow.RequestofPort, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Isolate provides a mock function with given fields: ctx, in, opts
func (_m *FlowServiceClient) Isolate(ctx context.Context, in *flow.RequestIsolate, opts ...grpc.CallOption) (*flow.ResponseFlow, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *flow.ResponseFlow
	if rf, ok := ret.Get(0).(func(context.Context, *flow.RequestIsolate, ...grpc.CallOption) *flow.ResponseFlow); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.ResponseFlow)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *flow.RequestIsolate, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields: ctx, in, opts
func (_m *FlowServiceClient) List(ctx context.Context, in *flow.RequestofPort, opts ...grpc.CallOption) (*flow.ResponseFlowList, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *flow.ResponseFlowList
	if rf, ok := ret.Get(0).(func(context.Context, *flow.RequestofPort, ...grpc.CallOption) *flow.ResponseFlowList); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.ResponseFlowList)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *flow.RequestofPort, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListPorts provides a mock function with given fields: ctx, in, opts
func (_m *FlowServiceClient) ListPorts(ctx context.Context, in *flow.RequestListPorts, opts ...grpc.CallOption) (*flow.ResponsePortList, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *flow.ResponsePortList
	if rf, ok := ret.Get(0).(func(context.Context, *flow.RequestListPorts, ...grpc.CallOption) *flow.ResponsePortList); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.ResponsePortList)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *flow.RequestListPorts, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Query provides a mock function with given fields: ctx, in, opts
func (_m *FlowServiceClient) Query(ctx context.Context, in *flow.RequestFlowofPort, opts ...grpc.CallOption) (*flow.ResponseFlowQuery, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *flow.ResponseFlowQuery
	if rf, ok := ret.Get(0).(func(context.Context, *flow.RequestFlowofPort, ...grpc.CallOption) *flow.ResponseFlowQuery); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.ResponseFlowQuery)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *flow.RequestFlowofPort, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Validate provides a mock function with given fields: ctx, in, opts
func (_m *FlowServiceClient) Validate(ctx context.Context, in *flow.RequestFlowCreate, opts ...grpc.CallOption) (*flow.ResponseFlow, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *flow.ResponseFlow
	if rf, ok := ret.Get(0).(func(context.Context, *flow.RequestFlowCreate, ...grpc.CallOption) *flow.ResponseFlow); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*flow.ResponseFlow)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *flow.RequestFlowCreate, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
