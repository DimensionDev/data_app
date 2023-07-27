package biz

import (
	"context"
	v1 "middle_platform/api/exchange_rate/v1"

	"github.com/go-kratos/kratos/v2/log"
)

type RateRepo interface {
	SupportedCurrencies(ctx context.Context, req *v1.RateRequest) (*v1.RateReply, error)
	BaseCurrency(ctx context.Context, req *v1.BaseCurrencyRequest) (*v1.BaseCurrencyReply, error)
}

// RateUsecase is a Rate usecase.
type RateUsecase struct {
	repo RateRepo
	log  *log.Helper
}

func NewRateUsecase(repo RateRepo, logger log.Logger) *RateUsecase {
	return &RateUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *RateUsecase) SupportedCurrencies(ctx context.Context, req *v1.RateRequest) (*v1.RateReply, error) {
	res, err := uc.repo.SupportedCurrencies(ctx, req)
	return res, err
}

func (uc *RateUsecase) BaseCurrency(ctx context.Context, req *v1.BaseCurrencyRequest) (*v1.BaseCurrencyReply, error) {
	res, err := uc.repo.BaseCurrency(ctx, req)
	return res, err
}
