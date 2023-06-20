package service

import (
	"context"

	pb "nft_transfer/api/nft_transfer/v1"
)

type NftTransferService struct {
	pb.UnimplementedNftTransferServer
}

func NewNftTransferService() *NftTransferService {
	return &NftTransferService{}
}

func (s *NftTransferService) GetNftTransfer(ctx context.Context, req *pb.GetNftTransferRequest) (*pb.GetNftTransferReply, error) {

	return &pb.GetNftTransferReply{}, nil
}
