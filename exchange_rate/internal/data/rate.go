package data

import (
	"context"
	"fmt"

	v1 "exchange_rate/api/rate/v1"
	"exchange_rate/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
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
