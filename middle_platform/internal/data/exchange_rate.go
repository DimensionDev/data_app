package data

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	v1 "middle_platform/api/exchange_rate/v1"
	"middle_platform/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	resty "github.com/go-resty/resty/v2"
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

const EXCHANGE_RATE_KEY_PREFIX = "exchange_rate_"
const CURRENCIES_KEY = "currencies"

type Currencies struct {
	Result          string     `json:"result"`
	Documentation   string     `json:"documentation"`
	Terms_of_use    string     `json:"terms_of_use"`
	Supported_codes [][]string `json:"supported_codes"`
}

func (uc *rateRepo) SupportedCurrencies(ctx context.Context, req *v1.RateRequest) (*v1.RateReply, error) {
	// url := "https://api.freecurrencyapi.com/v1/latest?apikey=fca_live_WCVTOTiGKXSLGF1BsI46ycAhsrkl42LCuF8jL5zF"
	// url := "https://v6.exchangerate-api.com/v6/382c9a9a547162edbc97f2fb/latest/USD"
	redis_res, _ := uc.GetCurrenciesfromRedis(ctx, req)
	if redis_res != nil {
		return redis_res, nil
	}

	// 缓存中没有数据，先从第三方查，再写入缓存
	codes_url := "https://v6.exchangerate-api.com/v6/382c9a9a547162edbc97f2fb/codes"
	client := resty.New()
	var result Currencies
	client.R().EnableTrace().SetResult(&result).Get(codes_url)

	var res v1.RateReply

	for i := 0; i < len(result.Supported_codes); i++ {
		code := result.Supported_codes[i][0]
		res.Currencies = append(res.Currencies, code)
	}

	is_success := uc.SetCurrenciesToRedis(ctx, res)
	if !is_success {
		uc.log.Errorf("SetCurrenciesToRedis fail")
	}

	return &res, nil
}

func (r *rateRepo) GetCurrenciesfromRedis(ctx context.Context, req *v1.RateRequest) (*v1.RateReply, error) {
	redis_key := EXCHANGE_RATE_KEY_PREFIX + CURRENCIES_KEY

	res := r.data.RedisCli.Get(ctx, redis_key)
	if res.Err() == nil {
		var result v1.RateReply
		err := json.Unmarshal([]byte(res.Val()), &result)
		if err != nil {
			r.log.Errorf("parse json Error[%v]", err)
			return nil, nil
		}
		return &result, nil
	}
	return nil, res.Err()
}

func (r *rateRepo) SetCurrenciesToRedis(ctx context.Context, data v1.RateReply) bool {
	jsonData, err := json.Marshal(data)
	if err != nil {
		r.log.Errorf("SetCurrenciesToRedis Json Marshal Error[%v]", err)
	}
	fmt.Println("json: ", jsonData)
	redis_key := EXCHANGE_RATE_KEY_PREFIX + CURRENCIES_KEY
	res := r.data.RedisCli.Set(ctx, redis_key, jsonData, 24*time.Hour)
	if res.Err() != nil {
		r.log.Errorf("SetCurrenciesToRedis Error[%v]", res.Err())
		return false
	}
	return true
}
