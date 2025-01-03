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
	// Add whitelist collection methods
	AddWhitelistCollection(ctx context.Context, req *pb.AddWhitelistCollectionRequest) (*pb.AddWhitelistCollectionReply, error)
	DeleteWhitelistCollection(ctx context.Context, req *pb.DeleteWhitelistCollectionRequest) (*pb.DeleteWhitelistCollectionReply, error)
	ListWhitelistCollections(ctx context.Context, req *pb.ListWhitelistCollectionsRequest) (*pb.ListWhitelistCollectionsReply, error)
	IsCollectionWhitelisted(ctx context.Context, collectionID string) (bool, error)
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

// AddWhitelistCollection implements the whitelist collection addition
func (uc *NftTransferUsecase) AddWhitelistCollection(ctx context.Context, req *pb.AddWhitelistCollectionRequest) (*pb.AddWhitelistCollectionReply, error) {
	return uc.repo.AddWhitelistCollection(ctx, req)
}

// DeleteWhitelistCollection implements the whitelist collection deletion
func (uc *NftTransferUsecase) DeleteWhitelistCollection(ctx context.Context, req *pb.DeleteWhitelistCollectionRequest) (*pb.DeleteWhitelistCollectionReply, error) {
	return uc.repo.DeleteWhitelistCollection(ctx, req)
}

// ListWhitelistCollections implements the whitelist collection listing
func (uc *NftTransferUsecase) ListWhitelistCollections(ctx context.Context, req *pb.ListWhitelistCollectionsRequest) (*pb.ListWhitelistCollectionsReply, error) {
	return uc.repo.ListWhitelistCollections(ctx, req)
}
