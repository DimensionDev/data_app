package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	pb "middle_platform/api/exchange_rate/v1"
	biz "middle_platform/internal/biz"
)

type ExchangeRateService struct {
	pb.UnimplementedExchangeRateServer
	usecase *biz.RateUsecase
}

func NewExchangeRateService(rate_info *biz.RateUsecase, logger log.Logger) *ExchangeRateService {
	return &ExchangeRateService{}
}

func (s *ExchangeRateService) SupportedCurrencies(ctx context.Context, req *pb.RateRequest) (*pb.RateReply, error) {
	res, err := s.usecase.ListAll(ctx, req)

	return res, err
}
