// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.23.3
// source: exchange_rate/v1/exchange_rate.proto

package v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	ExchangeRate_SupportedCurrencies_FullMethodName = "/exchange_rate.v1.ExchangeRate/SupportedCurrencies"
	ExchangeRate_BaseCurrency_FullMethodName        = "/exchange_rate.v1.ExchangeRate/BaseCurrency"
)

// ExchangeRateClient is the client API for ExchangeRate service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ExchangeRateClient interface {
	// get Supported Currencies
	SupportedCurrencies(ctx context.Context, in *RateRequest, opts ...grpc.CallOption) (*RateReply, error)
	BaseCurrency(ctx context.Context, in *BaseCurrencyRequest, opts ...grpc.CallOption) (*BaseCurrencyReply, error)
}

type exchangeRateClient struct {
	cc grpc.ClientConnInterface
}

func NewExchangeRateClient(cc grpc.ClientConnInterface) ExchangeRateClient {
	return &exchangeRateClient{cc}
}

func (c *exchangeRateClient) SupportedCurrencies(ctx context.Context, in *RateRequest, opts ...grpc.CallOption) (*RateReply, error) {
	out := new(RateReply)
	err := c.cc.Invoke(ctx, ExchangeRate_SupportedCurrencies_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *exchangeRateClient) BaseCurrency(ctx context.Context, in *BaseCurrencyRequest, opts ...grpc.CallOption) (*BaseCurrencyReply, error) {
	out := new(BaseCurrencyReply)
	err := c.cc.Invoke(ctx, ExchangeRate_BaseCurrency_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ExchangeRateServer is the server API for ExchangeRate service.
// All implementations must embed UnimplementedExchangeRateServer
// for forward compatibility
type ExchangeRateServer interface {
	// get Supported Currencies
	SupportedCurrencies(context.Context, *RateRequest) (*RateReply, error)
	BaseCurrency(context.Context, *BaseCurrencyRequest) (*BaseCurrencyReply, error)
	mustEmbedUnimplementedExchangeRateServer()
}

// UnimplementedExchangeRateServer must be embedded to have forward compatible implementations.
type UnimplementedExchangeRateServer struct {
}

func (UnimplementedExchangeRateServer) SupportedCurrencies(context.Context, *RateRequest) (*RateReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SupportedCurrencies not implemented")
}
func (UnimplementedExchangeRateServer) BaseCurrency(context.Context, *BaseCurrencyRequest) (*BaseCurrencyReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BaseCurrency not implemented")
}
func (UnimplementedExchangeRateServer) mustEmbedUnimplementedExchangeRateServer() {}

// UnsafeExchangeRateServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ExchangeRateServer will
// result in compilation errors.
type UnsafeExchangeRateServer interface {
	mustEmbedUnimplementedExchangeRateServer()
}

func RegisterExchangeRateServer(s grpc.ServiceRegistrar, srv ExchangeRateServer) {
	s.RegisterService(&ExchangeRate_ServiceDesc, srv)
}

func _ExchangeRate_SupportedCurrencies_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExchangeRateServer).SupportedCurrencies(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ExchangeRate_SupportedCurrencies_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExchangeRateServer).SupportedCurrencies(ctx, req.(*RateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ExchangeRate_BaseCurrency_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BaseCurrencyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExchangeRateServer).BaseCurrency(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ExchangeRate_BaseCurrency_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExchangeRateServer).BaseCurrency(ctx, req.(*BaseCurrencyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ExchangeRate_ServiceDesc is the grpc.ServiceDesc for ExchangeRate service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ExchangeRate_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "exchange_rate.v1.ExchangeRate",
	HandlerType: (*ExchangeRateServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SupportedCurrencies",
			Handler:    _ExchangeRate_SupportedCurrencies_Handler,
		},
		{
			MethodName: "BaseCurrency",
			Handler:    _ExchangeRate_BaseCurrency_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "exchange_rate/v1/exchange_rate.proto",
}