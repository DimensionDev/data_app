package server

import (
	exchange_rate "middle_platform/api/exchange_rate/v1"
	v1 "middle_platform/api/helloworld/v1"
	nv1 "middle_platform/api/nft_transfer/v1"
	"middle_platform/internal/conf"
	"middle_platform/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/ratelimit"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/handlers"
)

// Create a rate with the given limit (number of requests) for the given
// period (a time.Duration of your choice).

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, greeter *service.GreeterService, nft *service.NftTransferService, rater *service.ExchangeRateService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			NewCustomRateLimiter(),
			ratelimit.Server(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	opts = append(opts, http.Filter(handlers.CORS(
		// 域名配置
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"}),
	)))
	srv := http.NewServer(opts...)
	v1.RegisterGreeterHTTPServer(srv, greeter)
	nv1.RegisterNftTransferHTTPServer(srv, nft)
	exchange_rate.RegisterExchangeRateHTTPServer(srv, rater)
	return srv
}
