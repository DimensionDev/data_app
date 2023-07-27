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
const BASECURRENCY_KEY_PREFIX = "base_currency_"
const CURRENCIES_KEY = "currencies"

type Currencies struct {
	Result          string     `json:"result"`
	Documentation   string     `json:"documentation"`
	Terms_of_use    string     `json:"terms_of_use"`
	Supported_codes [][]string `json:"supported_codes"`
}

type BaseCurrency struct {
	Result                string             `json:"result"`
	Documentation         string             `json:"documentation"`
	Terms_of_use          string             `json:"terms_of_use"`
	Time_last_update_unix int32              `json:"time_last_update_unix"`
	Time_last_update_utc  string             `json:"time_last_update_utc"`
	Time_next_update_unix int32              `json:"time_next_update_unix"`
	Time_next_update_utc  string             `json:"time_next_update_utc"`
	Base_code             string             `json:"base_code"`
	Conversion_rates      map[string]float64 `json:"conversion_rates"`
}

func (uc *rateRepo) SupportedCurrencies(ctx context.Context, req *v1.RateRequest) (*v1.RateReply, error) {
	redis_res, _ := uc.GetCurrenciesfromRedis(ctx, req)
	if redis_res != nil {
		return redis_res, nil
	}

	// 缓存中没有数据，先从第三方查，再写入缓存
	codes_url := "https://v6.exchangerate-api.com/v6/382c9a9a547162edbc97f2fb/codes"
	client := resty.New()
	var result Currencies
	_, err := client.R().EnableTrace().SetResult(&result).Get(codes_url)
	if err != nil {
		uc.log.Errorf("get SupportedCurrencies fail from third party.")
		return nil, err
	}
	fmt.Println("http error", err)

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
	redis_key := EXCHANGE_RATE_KEY_PREFIX + CURRENCIES_KEY
	new_ctx := context.Background()
	res := r.data.RedisCli.Set(new_ctx, redis_key, jsonData, 24*time.Hour)
	if res.Err() != nil {
		r.log.Errorf("SetCurrenciesToRedis Error[%v]", res.Err())
		return false
	}
	return true
}

func (uc *rateRepo) BaseCurrency(ctx context.Context, req *v1.BaseCurrencyRequest) (*v1.BaseCurrencyReply, error) {

	// 从缓存中取
	redis_res, _ := uc.GetExchangeRatefromRedis(ctx, req)
	if redis_res != nil {
		return redis_res, nil
	}

	// 没有则从第三方取，然后写入缓存
	codes_url := "https://v6.exchangerate-api.com/v6/382c9a9a547162edbc97f2fb/latest/" + req.Base
	client := resty.New()
	var result BaseCurrency
	_, err := client.R().EnableTrace().SetResult(&result).Get(codes_url)
	fmt.Println("http error", err, req.Base)

	if result.Result != "success" {
		uc.log.Errorf("get SupportedCurrencies fail from third party. param:", req.Base)
		return nil, err
	}
	var res v1.BaseCurrencyReply

	res.Timestamp = result.Time_last_update_unix
	res.Base = result.Base_code
	res.Rates = result.Conversion_rates
	// for _, value := range result.Conversion_rates {
	// 	res.Rates = append(res.Rates, value)
	// }
	// for i := 0; i < len(result.Conversion_rates); i++ {
	// 	code := result.Conversion_rates[i][0]
	// 	res.Currencies = append(res.Currencies, code)
	// }

	is_success := uc.SetExchangeRateToRedis(ctx, req, res)
	if !is_success {
		uc.log.Errorf("SetExchangeRateToRedis fail")
	}

	return &res, nil
}

func (r *rateRepo) GetExchangeRatefromRedis(ctx context.Context, req *v1.BaseCurrencyRequest) (*v1.BaseCurrencyReply, error) {
	redis_key := BASECURRENCY_KEY_PREFIX + req.Base

	res := r.data.RedisCli.Get(ctx, redis_key)
	if res.Err() == nil {
		var result v1.BaseCurrencyReply
		err := json.Unmarshal([]byte(res.Val()), &result)
		if err != nil {
			r.log.Errorf("parse json Error[%v]", err)
			return nil, nil
		}
		return &result, nil
	}
	return nil, res.Err()
}

func (r *rateRepo) SetExchangeRateToRedis(ctx context.Context, req *v1.BaseCurrencyRequest, data v1.BaseCurrencyReply) bool {
	jsonData, err := json.Marshal(data)
	if err != nil {
		r.log.Errorf("SetExchangeRateToRedis Json Marshal Error[%v]", err)
	}
	redis_key := BASECURRENCY_KEY_PREFIX + req.Base
	new_ctx := context.Background()
	res := r.data.RedisCli.Set(new_ctx, redis_key, jsonData, 24*time.Hour)
	if res.Err() != nil {
		r.log.Errorf("SetExchangeRateToRedis Error[%v]", res.Err())
		return false
	}
	return true
}
