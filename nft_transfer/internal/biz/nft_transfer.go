package biz

import (
	"context"
	"fmt"
	pb "nft_transfer/api/nft_transfer/v1"

	"github.com/go-kratos/kratos/v2/log"
)

// NftTransferRepo is a NftTransfer repo.
type NftTransferRepo interface {
	GetHandleNftinfo(ctx context.Context, req *pb.GetNftTransferRequest) (*pb.GetNftTransferReply, error)
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
	fmt.Print("bizyyyyyyy:", res)
	return res, err
}
