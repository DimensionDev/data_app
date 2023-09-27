package service

import (
	"context"

	pb "middle_platform/api/exchange_rate/v1"
	biz "middle_platform/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type ExchangeRateService struct {
	pb.UnimplementedExchangeRateServer
	usecase *biz.RateUsecase
}

func NewExchangeRateService(rate_info *biz.RateUsecase, logger log.Logger) *ExchangeRateService {
	return &ExchangeRateService{usecase: rate_info}
}

func (s *ExchangeRateService) SupportedCurrencies(ctx context.Context, req *pb.RateRequest) (*pb.RateReply, error) {
	res, err := s.usecase.SupportedCurrencies(ctx, req)

	return res, err
}

func (s *ExchangeRateService) BaseCurrency(ctx context.Context, req *pb.BaseCurrencyRequest) (*pb.BaseCurrencyReply, error) {
	res, err := s.usecase.BaseCurrency(ctx, req)
	return res, err
}
