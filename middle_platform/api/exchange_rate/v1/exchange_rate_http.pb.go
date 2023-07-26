// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.6.3
// - protoc             v4.23.3
// source: exchange_rate/v1/exchange_rate.proto

package v1

import (
	context "context"
	http "github.com/go-kratos/kratos/v2/transport/http"
	binding "github.com/go-kratos/kratos/v2/transport/http/binding"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.EncodeURL

const _ = http.SupportPackageIsVersion1

const OperationExchangeRateSupportedCurrencies = "/exchange_rate.v1.ExchangeRate/SupportedCurrencies"

type ExchangeRateHTTPServer interface {
	// SupportedCurrencies get Supported Currencies
	SupportedCurrencies(context.Context, *RateRequest) (*RateReply, error)
}

func RegisterExchangeRateHTTPServer(s *http.Server, srv ExchangeRateHTTPServer) {
	r := s.Route("/")
	r.GET("/supported-currencies", _ExchangeRate_SupportedCurrencies0_HTTP_Handler(srv))
}

func _ExchangeRate_SupportedCurrencies0_HTTP_Handler(srv ExchangeRateHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in RateRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationExchangeRateSupportedCurrencies)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.SupportedCurrencies(ctx, req.(*RateRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*RateReply)
		return ctx.Result(200, reply)
	}
}

type ExchangeRateHTTPClient interface {
	SupportedCurrencies(ctx context.Context, req *RateRequest, opts ...http.CallOption) (rsp *RateReply, err error)
}

type ExchangeRateHTTPClientImpl struct {
	cc *http.Client
}

func NewExchangeRateHTTPClient(client *http.Client) ExchangeRateHTTPClient {
	return &ExchangeRateHTTPClientImpl{client}
}

func (c *ExchangeRateHTTPClientImpl) SupportedCurrencies(ctx context.Context, in *RateRequest, opts ...http.CallOption) (*RateReply, error) {
	var out RateReply
	pattern := "/supported-currencies"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationExchangeRateSupportedCurrencies))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}
