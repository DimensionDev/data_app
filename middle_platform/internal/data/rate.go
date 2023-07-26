package data

import (
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2/log"
	v1 "middle_platform/api/exchange_rate/v1"
	"middle_platform/internal/biz"
)

type rateRepo struct {
	data *Data
	log  *log.Helper
}

func NewRateRepo(data *Data, logger log.Logger) biz.RateRepo {
	return &rateRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *rateRepo) ListAll(context.Context, *v1.RateRequest) ([]*biz.Rate, error) {
	fmt.Print("test2")
	return nil, nil
}
