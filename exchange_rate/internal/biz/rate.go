package biz

import (
	"context"
	v1 "exchange_rate/api/rate/v1"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	resty "github.com/go-resty/resty/v2"
)

var (
	// ErrUserNotFound is user not found.
	ErrUserNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
)

// data model.
type Rate struct {
	Message string
}

type RateRepo interface {
	ListAll(ctx context.Context, req *v1.RateRequest) ([]*Rate, error)
}

// RateUsecase is a Rate usecase.
type RateUsecase struct {
	repo RateRepo
	log  *log.Helper
}

// NewRateUsecase new a Rate usecase.
func NewRateUsecase(repo RateRepo, logger log.Logger) *RateUsecase {
	return &RateUsecase{repo: repo, log: log.NewHelper(logger)}
}

type Currencies struct {
	Result          string     `json:"result"`
	Documentation   string     `json:"documentation"`
	Terms_of_use    string     `json:"terms_of_use"`
	Supported_codes [][]string `json:"supported_codes"`
}

func (uc *RateUsecase) ListAll(ctx context.Context, req *v1.RateRequest) (*v1.RateReply, error) {
	// url := "https://api.freecurrencyapi.com/v1/latest?apikey=fca_live_WCVTOTiGKXSLGF1BsI46ycAhsrkl42LCuF8jL5zF"
	// url := "https://v6.exchangerate-api.com/v6/382c9a9a547162edbc97f2fb/latest/USD"
	codes_url := "https://v6.exchangerate-api.com/v6/382c9a9a547162edbc97f2fb/codes"
	client := resty.New()
	var result Currencies
	client.R().EnableTrace().SetResult(&result).Get(codes_url)

	var res v1.RateReply

	for i := 0; i < len(result.Supported_codes); i++ {
		code := result.Supported_codes[i][0]
		res.Currencies = append(res.Currencies, code)
	}

	return &res, nil
}
