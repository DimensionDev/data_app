package biz

import (
	"context"
	pb "middle_platform/api/nft_transfer/v1"

	"github.com/go-kratos/kratos/v2/log"
)

// NftTransferRepo is a NftTransfer repo.
type NftTransferRepo interface {
	GetHandleNftinfo(ctx context.Context, req *pb.GetNftTransferRequest) (*pb.GetNftTransferReply, error)
	GetSpamReport(ctx context.Context, req *pb.GetReportSpamRequest) (*pb.GetReportSpamReply, error)
	GetTransferNft(ctx context.Context, req *pb.GetTransferNftRequest) (*pb.GetTransferNftReply, error)
	PostSpamReport(ctx context.Context, req *pb.PostReportSpamRequest) (*pb.PostReportSpamReply, error)
	PostNftMute(ctx context.Context, req *pb.PostReportAccountMuteRequest) (*pb.PostReportAccountMuteReply, error)
}

// NftTransferUsecase
type NftTransferUsecase struct {
	repo NftTransferRepo
	log  *log.Helper
}

// NewNftTransferUsecase
func NewNftTransferUsecase(repo NftTransferRepo, logger log.Logger) *NftTransferUsecase {
	return &NftTransferUsecase{repo: repo, log: log.NewHelper(logger)}
}

// GetHandleNftinfo
func (uc *NftTransferUsecase) GetHandleNftinfo(ctx context.Context, req *pb.GetNftTransferRequest) (*pb.GetNftTransferReply, error) {
	res, err := uc.repo.GetHandleNftinfo(ctx, req)
	/*if err != nil {
		return nil, err
	}*/
	//fmt.Print("bizyyyyyyy:", res)
	return res, err
}

func (uc *NftTransferUsecase) GetSpamReport(ctx context.Context, req *pb.GetReportSpamRequest) (*pb.GetReportSpamReply, error) {
	res, err := uc.repo.GetSpamReport(ctx, req)
	return res, err
}

func (uc *NftTransferUsecase) PostSpamReport(ctx context.Context, req *pb.PostReportSpamRequest) (*pb.PostReportSpamReply, error) {
	res, err := uc.repo.PostSpamReport(ctx, req)
	return res, err
}

func (uc *NftTransferUsecase) GetTransferNft(ctx context.Context, req *pb.GetTransferNftRequest) (*pb.GetTransferNftReply, error) {
	res, err := uc.repo.GetTransferNft(ctx, req)
	return res, err
}

func (uc *NftTransferUsecase) PostNftMute(ctx context.Context, req *pb.PostReportAccountMuteRequest) (*pb.PostReportAccountMuteReply, error) {
	res, err := uc.repo.PostNftMute(ctx, req)
	return res, err
}
