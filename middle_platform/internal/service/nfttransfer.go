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
