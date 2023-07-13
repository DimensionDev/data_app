package service

import (
	"context"
	"nft_transfer/internal/biz"

	"github.com/go-kratos/kratos/v2/log"

	pb "nft_transfer/api/nft_transfer/v1"
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
