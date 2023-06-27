// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"nft_transfer/internal/biz"
	"nft_transfer/internal/conf"
	"nft_transfer/internal/data"
	"nft_transfer/internal/server"
	"nft_transfer/internal/service"
)

import (
	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(confServer *conf.Server, confData *conf.Data, logger log.Logger) (*kratos.App, func(), error) {
	/*dataData, cleanup, err := data.NewData(confData, logger)
	if err != nil {
		return nil, nil, err
	}*/

	db, err := data.NewDataBase(confData, logger)
	if err != nil {
		return nil, nil, err
	}
	client, cleanup, err := data.NewRedis(confData, logger)
	if err != nil {
		return nil, nil, err
	}
	dataData, cleanup2, err := data.NewData(confData, logger, db, client)
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
	httpServer := server.NewHTTPServer(confServer, greeterService, nftTransferService, logger)
	app := newApp(logger, grpcServer, httpServer)
	return app, func() {
		cleanup2()
		cleanup()
	}, nil
}
