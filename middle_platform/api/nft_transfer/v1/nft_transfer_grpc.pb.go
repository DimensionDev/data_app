// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.23.3
// source: nft_transfer/v1/nft_transfer.proto

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
	NftTransfer_GetNftTransfer_FullMethodName            = "/api.nft_transfer.v1.NftTransfer/GetNftTransfer"
	NftTransfer_GetReportSpam_FullMethodName             = "/api.nft_transfer.v1.NftTransfer/GetReportSpam"
	NftTransfer_PostReportSpam_FullMethodName            = "/api.nft_transfer.v1.NftTransfer/PostReportSpam"
	NftTransfer_GetTransferNft_FullMethodName            = "/api.nft_transfer.v1.NftTransfer/GetTransferNft"
	NftTransfer_PostReportAccountMute_FullMethodName     = "/api.nft_transfer.v1.NftTransfer/PostReportAccountMute"
	NftTransfer_AddWhitelistCollection_FullMethodName    = "/api.nft_transfer.v1.NftTransfer/AddWhitelistCollection"
	NftTransfer_DeleteWhitelistCollection_FullMethodName = "/api.nft_transfer.v1.NftTransfer/DeleteWhitelistCollection"
	NftTransfer_ListWhitelistCollections_FullMethodName  = "/api.nft_transfer.v1.NftTransfer/ListWhitelistCollections"
	NftTransfer_AddWhitelistAddress_FullMethodName       = "/api.nft_transfer.v1.NftTransfer/AddWhitelistAddress"
	NftTransfer_DeleteWhitelistAddress_FullMethodName    = "/api.nft_transfer.v1.NftTransfer/DeleteWhitelistAddress"
	NftTransfer_ListWhitelistAddresses_FullMethodName    = "/api.nft_transfer.v1.NftTransfer/ListWhitelistAddresses"
)

// NftTransferClient is the client API for NftTransfer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NftTransferClient interface {
	GetNftTransfer(ctx context.Context, in *GetNftTransferRequest, opts ...grpc.CallOption) (*GetNftTransferReply, error)
	GetReportSpam(ctx context.Context, in *GetReportSpamRequest, opts ...grpc.CallOption) (*GetReportSpamReply, error)
	PostReportSpam(ctx context.Context, in *PostReportSpamRequest, opts ...grpc.CallOption) (*PostReportSpamReply, error)
	GetTransferNft(ctx context.Context, in *GetTransferNftRequest, opts ...grpc.CallOption) (*GetTransferNftReply, error)
	PostReportAccountMute(ctx context.Context, in *PostReportAccountMuteRequest, opts ...grpc.CallOption) (*PostReportAccountMuteReply, error)
	AddWhitelistCollection(ctx context.Context, in *AddWhitelistCollectionRequest, opts ...grpc.CallOption) (*AddWhitelistCollectionReply, error)
	DeleteWhitelistCollection(ctx context.Context, in *DeleteWhitelistCollectionRequest, opts ...grpc.CallOption) (*DeleteWhitelistCollectionReply, error)
	ListWhitelistCollections(ctx context.Context, in *ListWhitelistCollectionsRequest, opts ...grpc.CallOption) (*ListWhitelistCollectionsReply, error)
	AddWhitelistAddress(ctx context.Context, in *AddWhitelistAddressRequest, opts ...grpc.CallOption) (*AddWhitelistAddressReply, error)
	DeleteWhitelistAddress(ctx context.Context, in *DeleteWhitelistAddressRequest, opts ...grpc.CallOption) (*DeleteWhitelistAddressReply, error)
	ListWhitelistAddresses(ctx context.Context, in *ListWhitelistAddressesRequest, opts ...grpc.CallOption) (*ListWhitelistAddressesReply, error)
}

type nftTransferClient struct {
	cc grpc.ClientConnInterface
}

func NewNftTransferClient(cc grpc.ClientConnInterface) NftTransferClient {
	return &nftTransferClient{cc}
}

func (c *nftTransferClient) GetNftTransfer(ctx context.Context, in *GetNftTransferRequest, opts ...grpc.CallOption) (*GetNftTransferReply, error) {
	out := new(GetNftTransferReply)
	err := c.cc.Invoke(ctx, NftTransfer_GetNftTransfer_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nftTransferClient) GetReportSpam(ctx context.Context, in *GetReportSpamRequest, opts ...grpc.CallOption) (*GetReportSpamReply, error) {
	out := new(GetReportSpamReply)
	err := c.cc.Invoke(ctx, NftTransfer_GetReportSpam_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nftTransferClient) PostReportSpam(ctx context.Context, in *PostReportSpamRequest, opts ...grpc.CallOption) (*PostReportSpamReply, error) {
	out := new(PostReportSpamReply)
	err := c.cc.Invoke(ctx, NftTransfer_PostReportSpam_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nftTransferClient) GetTransferNft(ctx context.Context, in *GetTransferNftRequest, opts ...grpc.CallOption) (*GetTransferNftReply, error) {
	out := new(GetTransferNftReply)
	err := c.cc.Invoke(ctx, NftTransfer_GetTransferNft_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nftTransferClient) PostReportAccountMute(ctx context.Context, in *PostReportAccountMuteRequest, opts ...grpc.CallOption) (*PostReportAccountMuteReply, error) {
	out := new(PostReportAccountMuteReply)
	err := c.cc.Invoke(ctx, NftTransfer_PostReportAccountMute_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nftTransferClient) AddWhitelistCollection(ctx context.Context, in *AddWhitelistCollectionRequest, opts ...grpc.CallOption) (*AddWhitelistCollectionReply, error) {
	out := new(AddWhitelistCollectionReply)
	err := c.cc.Invoke(ctx, NftTransfer_AddWhitelistCollection_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nftTransferClient) DeleteWhitelistCollection(ctx context.Context, in *DeleteWhitelistCollectionRequest, opts ...grpc.CallOption) (*DeleteWhitelistCollectionReply, error) {
	out := new(DeleteWhitelistCollectionReply)
	err := c.cc.Invoke(ctx, NftTransfer_DeleteWhitelistCollection_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nftTransferClient) ListWhitelistCollections(ctx context.Context, in *ListWhitelistCollectionsRequest, opts ...grpc.CallOption) (*ListWhitelistCollectionsReply, error) {
	out := new(ListWhitelistCollectionsReply)
	err := c.cc.Invoke(ctx, NftTransfer_ListWhitelistCollections_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nftTransferClient) AddWhitelistAddress(ctx context.Context, in *AddWhitelistAddressRequest, opts ...grpc.CallOption) (*AddWhitelistAddressReply, error) {
	out := new(AddWhitelistAddressReply)
	err := c.cc.Invoke(ctx, NftTransfer_AddWhitelistAddress_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nftTransferClient) DeleteWhitelistAddress(ctx context.Context, in *DeleteWhitelistAddressRequest, opts ...grpc.CallOption) (*DeleteWhitelistAddressReply, error) {
	out := new(DeleteWhitelistAddressReply)
	err := c.cc.Invoke(ctx, NftTransfer_DeleteWhitelistAddress_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nftTransferClient) ListWhitelistAddresses(ctx context.Context, in *ListWhitelistAddressesRequest, opts ...grpc.CallOption) (*ListWhitelistAddressesReply, error) {
	out := new(ListWhitelistAddressesReply)
	err := c.cc.Invoke(ctx, NftTransfer_ListWhitelistAddresses_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NftTransferServer is the server API for NftTransfer service.
// All implementations must embed UnimplementedNftTransferServer
// for forward compatibility
type NftTransferServer interface {
	GetNftTransfer(context.Context, *GetNftTransferRequest) (*GetNftTransferReply, error)
	GetReportSpam(context.Context, *GetReportSpamRequest) (*GetReportSpamReply, error)
	PostReportSpam(context.Context, *PostReportSpamRequest) (*PostReportSpamReply, error)
	GetTransferNft(context.Context, *GetTransferNftRequest) (*GetTransferNftReply, error)
	PostReportAccountMute(context.Context, *PostReportAccountMuteRequest) (*PostReportAccountMuteReply, error)
	AddWhitelistCollection(context.Context, *AddWhitelistCollectionRequest) (*AddWhitelistCollectionReply, error)
	DeleteWhitelistCollection(context.Context, *DeleteWhitelistCollectionRequest) (*DeleteWhitelistCollectionReply, error)
	ListWhitelistCollections(context.Context, *ListWhitelistCollectionsRequest) (*ListWhitelistCollectionsReply, error)
	AddWhitelistAddress(context.Context, *AddWhitelistAddressRequest) (*AddWhitelistAddressReply, error)
	DeleteWhitelistAddress(context.Context, *DeleteWhitelistAddressRequest) (*DeleteWhitelistAddressReply, error)
	ListWhitelistAddresses(context.Context, *ListWhitelistAddressesRequest) (*ListWhitelistAddressesReply, error)
	mustEmbedUnimplementedNftTransferServer()
}

// UnimplementedNftTransferServer must be embedded to have forward compatible implementations.
type UnimplementedNftTransferServer struct {
}

func (UnimplementedNftTransferServer) GetNftTransfer(context.Context, *GetNftTransferRequest) (*GetNftTransferReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetNftTransfer not implemented")
}
func (UnimplementedNftTransferServer) GetReportSpam(context.Context, *GetReportSpamRequest) (*GetReportSpamReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetReportSpam not implemented")
}
func (UnimplementedNftTransferServer) PostReportSpam(context.Context, *PostReportSpamRequest) (*PostReportSpamReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PostReportSpam not implemented")
}
func (UnimplementedNftTransferServer) GetTransferNft(context.Context, *GetTransferNftRequest) (*GetTransferNftReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTransferNft not implemented")
}
func (UnimplementedNftTransferServer) PostReportAccountMute(context.Context, *PostReportAccountMuteRequest) (*PostReportAccountMuteReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PostReportAccountMute not implemented")
}
func (UnimplementedNftTransferServer) AddWhitelistCollection(context.Context, *AddWhitelistCollectionRequest) (*AddWhitelistCollectionReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddWhitelistCollection not implemented")
}
func (UnimplementedNftTransferServer) DeleteWhitelistCollection(context.Context, *DeleteWhitelistCollectionRequest) (*DeleteWhitelistCollectionReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteWhitelistCollection not implemented")
}
func (UnimplementedNftTransferServer) ListWhitelistCollections(context.Context, *ListWhitelistCollectionsRequest) (*ListWhitelistCollectionsReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListWhitelistCollections not implemented")
}
func (UnimplementedNftTransferServer) AddWhitelistAddress(context.Context, *AddWhitelistAddressRequest) (*AddWhitelistAddressReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddWhitelistAddress not implemented")
}
func (UnimplementedNftTransferServer) DeleteWhitelistAddress(context.Context, *DeleteWhitelistAddressRequest) (*DeleteWhitelistAddressReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteWhitelistAddress not implemented")
}
func (UnimplementedNftTransferServer) ListWhitelistAddresses(context.Context, *ListWhitelistAddressesRequest) (*ListWhitelistAddressesReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListWhitelistAddresses not implemented")
}
func (UnimplementedNftTransferServer) mustEmbedUnimplementedNftTransferServer() {}

// UnsafeNftTransferServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NftTransferServer will
// result in compilation errors.
type UnsafeNftTransferServer interface {
	mustEmbedUnimplementedNftTransferServer()
}

func RegisterNftTransferServer(s grpc.ServiceRegistrar, srv NftTransferServer) {
	s.RegisterService(&NftTransfer_ServiceDesc, srv)
}

func _NftTransfer_GetNftTransfer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetNftTransferRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NftTransferServer).GetNftTransfer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NftTransfer_GetNftTransfer_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NftTransferServer).GetNftTransfer(ctx, req.(*GetNftTransferRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NftTransfer_GetReportSpam_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetReportSpamRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NftTransferServer).GetReportSpam(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NftTransfer_GetReportSpam_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NftTransferServer).GetReportSpam(ctx, req.(*GetReportSpamRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NftTransfer_PostReportSpam_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PostReportSpamRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NftTransferServer).PostReportSpam(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NftTransfer_PostReportSpam_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NftTransferServer).PostReportSpam(ctx, req.(*PostReportSpamRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NftTransfer_GetTransferNft_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTransferNftRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NftTransferServer).GetTransferNft(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NftTransfer_GetTransferNft_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NftTransferServer).GetTransferNft(ctx, req.(*GetTransferNftRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NftTransfer_PostReportAccountMute_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PostReportAccountMuteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NftTransferServer).PostReportAccountMute(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NftTransfer_PostReportAccountMute_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NftTransferServer).PostReportAccountMute(ctx, req.(*PostReportAccountMuteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NftTransfer_AddWhitelistCollection_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddWhitelistCollectionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NftTransferServer).AddWhitelistCollection(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NftTransfer_AddWhitelistCollection_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NftTransferServer).AddWhitelistCollection(ctx, req.(*AddWhitelistCollectionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NftTransfer_DeleteWhitelistCollection_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteWhitelistCollectionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NftTransferServer).DeleteWhitelistCollection(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NftTransfer_DeleteWhitelistCollection_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NftTransferServer).DeleteWhitelistCollection(ctx, req.(*DeleteWhitelistCollectionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NftTransfer_ListWhitelistCollections_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListWhitelistCollectionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NftTransferServer).ListWhitelistCollections(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NftTransfer_ListWhitelistCollections_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NftTransferServer).ListWhitelistCollections(ctx, req.(*ListWhitelistCollectionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NftTransfer_AddWhitelistAddress_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddWhitelistAddressRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NftTransferServer).AddWhitelistAddress(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NftTransfer_AddWhitelistAddress_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NftTransferServer).AddWhitelistAddress(ctx, req.(*AddWhitelistAddressRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NftTransfer_DeleteWhitelistAddress_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteWhitelistAddressRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NftTransferServer).DeleteWhitelistAddress(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NftTransfer_DeleteWhitelistAddress_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NftTransferServer).DeleteWhitelistAddress(ctx, req.(*DeleteWhitelistAddressRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NftTransfer_ListWhitelistAddresses_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListWhitelistAddressesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NftTransferServer).ListWhitelistAddresses(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NftTransfer_ListWhitelistAddresses_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NftTransferServer).ListWhitelistAddresses(ctx, req.(*ListWhitelistAddressesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// NftTransfer_ServiceDesc is the grpc.ServiceDesc for NftTransfer service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var NftTransfer_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.nft_transfer.v1.NftTransfer",
	HandlerType: (*NftTransferServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetNftTransfer",
			Handler:    _NftTransfer_GetNftTransfer_Handler,
		},
		{
			MethodName: "GetReportSpam",
			Handler:    _NftTransfer_GetReportSpam_Handler,
		},
		{
			MethodName: "PostReportSpam",
			Handler:    _NftTransfer_PostReportSpam_Handler,
		},
		{
			MethodName: "GetTransferNft",
			Handler:    _NftTransfer_GetTransferNft_Handler,
		},
		{
			MethodName: "PostReportAccountMute",
			Handler:    _NftTransfer_PostReportAccountMute_Handler,
		},
		{
			MethodName: "AddWhitelistCollection",
			Handler:    _NftTransfer_AddWhitelistCollection_Handler,
		},
		{
			MethodName: "DeleteWhitelistCollection",
			Handler:    _NftTransfer_DeleteWhitelistCollection_Handler,
		},
		{
			MethodName: "ListWhitelistCollections",
			Handler:    _NftTransfer_ListWhitelistCollections_Handler,
		},
		{
			MethodName: "AddWhitelistAddress",
			Handler:    _NftTransfer_AddWhitelistAddress_Handler,
		},
		{
			MethodName: "DeleteWhitelistAddress",
			Handler:    _NftTransfer_DeleteWhitelistAddress_Handler,
		},
		{
			MethodName: "ListWhitelistAddresses",
			Handler:    _NftTransfer_ListWhitelistAddresses_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "nft_transfer/v1/nft_transfer.proto",
}
