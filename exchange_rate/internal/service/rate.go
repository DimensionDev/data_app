package service

import (
	"context"
	v1 "exchange_rate/api/rate/v1"
	"exchange_rate/internal/biz"
)

type RateService struct {
	v1.UnimplementedRateServer
	uc *biz.RateUsecase
}

func NewRateService(uc *biz.RateUsecase) *RateService {
	return &RateService{uc: uc}
}

func (s *RateService) SupportedCurrencies(ctx context.Context, req *v1.RateRequest) (*v1.RateReply, error) {
	result, err := s.uc.ListAll(ctx, req)
	if err != nil {
		return nil, err
	}
	return result, nil
}
