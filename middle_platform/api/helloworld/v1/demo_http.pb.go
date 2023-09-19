// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.6.2
// - protoc             v4.23.3
// source: helloworld/v1/demo.proto

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

const OperationDemoGetDemo = "/api.helloworld.v1.Demo/GetDemo"

type DemoHTTPServer interface {
	GetDemo(context.Context, *GetDemoRequest) (*GetDemoReply, error)
}

func RegisterDemoHTTPServer(s *http.Server, srv DemoHTTPServer) {
	r := s.Route("/")
	r.GET("/demo", _Demo_GetDemo0_HTTP_Handler(srv))
}

func _Demo_GetDemo0_HTTP_Handler(srv DemoHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetDemoRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationDemoGetDemo)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetDemo(ctx, req.(*GetDemoRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetDemoReply)
		return ctx.Result(200, reply)
	}
}

type DemoHTTPClient interface {
	GetDemo(ctx context.Context, req *GetDemoRequest, opts ...http.CallOption) (rsp *GetDemoReply, err error)
}

type DemoHTTPClientImpl struct {
	cc *http.Client
}

func NewDemoHTTPClient(client *http.Client) DemoHTTPClient {
	return &DemoHTTPClientImpl{client}
}

func (c *DemoHTTPClientImpl) GetDemo(ctx context.Context, in *GetDemoRequest, opts ...http.CallOption) (*GetDemoReply, error) {
	var out GetDemoReply
	pattern := "/demo"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationDemoGetDemo))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}
