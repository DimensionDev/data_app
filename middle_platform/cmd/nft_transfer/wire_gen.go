// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"middle_platform/internal/biz"
	"middle_platform/internal/conf"
	"middle_platform/internal/data"
	"middle_platform/internal/server"
	"middle_platform/internal/service"
)

import (
	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(confServer *conf.Server, confData *conf.Data, logger log.Logger) (*kratos.App, func(), error) {
	gateway, err := data.NewDataBase(confData, logger)
	if err != nil {
		return nil, nil, err
	}
	client, cleanup, err := data.NewRedis(confData, logger)
	if err != nil {
		return nil, nil, err
	}
	dataData, cleanup2, err := data.NewData(confData, logger, gateway, client)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	greeterRepo := data.NewGreeterRepo(dataData, logger)
	greeterUsecase := biz.NewGreeterUsecase(greeterRepo, logger)
	greeterService := service.NewGreeterService(greeterUsecase)
	nftTransferRepo := data.NewNftTransferRepo(dataData, logger)
	nftTransferUsecase := biz.NewNftTransferUsecase(nftTransferRepo, logger)
	nftTransferService := service.NewNftTransferService(nftTransferUsecase, logger)
	grpcServer := server.NewGRPCServer(confServer, greeterService, nftTransferService, logger)
	rateRepo := data.NewRateRepo(dataData, logger)
	rateUsecase := biz.NewRateUsecase(rateRepo, logger)
	exchangeRateService := service.NewExchangeRateService(rateUsecase, logger)
	httpServer := server.NewHTTPServer(confServer, greeterService, nftTransferService, exchangeRateService, logger)
	app := newApp(logger, grpcServer, httpServer)
	return app, func() {
		cleanup2()
		cleanup()
	}, nil
}