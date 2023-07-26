package service

import (
	"context"

	pb "middle_platform/api/helloworld/v1"
)

type DemoService struct {
	pb.UnimplementedDemoServer
}

func NewDemoService() *DemoService {
	return &DemoService{}
}

func (s *DemoService) GetDemo(ctx context.Context, req *pb.GetDemoRequest) (*pb.GetDemoReply, error) {
	return &pb.GetDemoReply{}, nil
}
