package service

import (
	"context"
	"middle_platform/internal/biz"

	"github.com/go-kratos/kratos/v2/log"

	pb "middle_platform/api/nft_transfer/v1"
)

type NftTransferService struct {
	pb.UnimplementedNftTransferServer
	uc  *biz.NftTransferUsecase
	log *log.Helper
}

func NewNftTransferService(nftinfo *biz.NftTransferUsecase, logger log.Logger) *NftTransferService {
	return &NftTransferService{uc: nftinfo, log: log.NewHelper(logger)}
}

func (s *NftTransferService) GetNftTransfer(ctx context.Context, req *pb.GetNftTransferRequest) (*pb.GetNftTransferReply, error) {

	res, err := s.uc.GetHandleNftinfo(ctx, req)
	return res, err
	/*
		return &pb.GetNftTransferReply{
			Code:    200,
			Reason:  "SUCCESS",
			Message: "SUCCESS",
		}, nil

	*/

}

func (s *NftTransferService) GetReportSpam(ctx context.Context, req *pb.GetReportSpamRequest) (*pb.GetReportSpamReply, error) {
	res, err := s.uc.GetSpamReport(ctx, req)
	return res, err
}

// PostSpamReport
func (s *NftTransferService) PostReportSpam(ctx context.Context, req *pb.PostReportSpamRequest) (*pb.PostReportSpamReply, error) {
	res, err := s.uc.PostSpamReport(ctx, req)
	return res, err
}

func (s *NftTransferService) GetTransferNft(ctx context.Context, req *pb.GetTransferNftRequest) (*pb.GetTransferNftReply, error) {
	res, err := s.uc.GetTransferNft(ctx, req)
	return res, err
}

func (s *NftTransferService) PostReportAccountMute(ctx context.Context, req *pb.PostReportAccountMuteRequest) (*pb.PostReportAccountMuteReply, error) {
	res, err := s.uc.PostNftMute(ctx, req)
	return res, err
}

// AddWhitelistCollection 添加白名单collection
func (s *NftTransferService) AddWhitelistCollection(ctx context.Context, req *pb.AddWhitelistCollectionRequest) (*pb.AddWhitelistCollectionReply, error) {
	return s.uc.AddWhitelistCollection(ctx, req)
}

// DeleteWhitelistCollection 删除白名单collection
func (s *NftTransferService) DeleteWhitelistCollection(ctx context.Context, req *pb.DeleteWhitelistCollectionRequest) (*pb.DeleteWhitelistCollectionReply, error) {
	return s.uc.DeleteWhitelistCollection(ctx, req)
}

// ListWhitelistCollections 获取白名单collection列表
func (s *NftTransferService) ListWhitelistCollections(ctx context.Context, req *pb.ListWhitelistCollectionsRequest) (*pb.ListWhitelistCollectionsReply, error) {
	return s.uc.ListWhitelistCollections(ctx, req)
}

// AddWhitelistAddress 添加白名单地址
func (s *NftTransferService) AddWhitelistAddress(ctx context.Context, req *pb.AddWhitelistAddressRequest) (*pb.AddWhitelistAddressReply, error) {
	return s.uc.AddWhitelistAddress(ctx, req)
}

// DeleteWhitelistAddress 删除白名单地址
func (s *NftTransferService) DeleteWhitelistAddress(ctx context.Context, req *pb.DeleteWhitelistAddressRequest) (*pb.DeleteWhitelistAddressReply, error) {
	return s.uc.DeleteWhitelistAddress(ctx, req)
}

// ListWhitelistAddresses 获取白名单地址列表
func (s *NftTransferService) ListWhitelistAddresses(ctx context.Context, req *pb.ListWhitelistAddressesRequest) (*pb.ListWhitelistAddressesReply, error) {
	return s.uc.ListWhitelistAddresses(ctx, req)
}
