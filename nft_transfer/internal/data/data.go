package data

import (
	"context"
	"fmt"
	"github.com/bytehouse-cloud/driver-go/sdk"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis"
	"github.com/google/wire"
	//"gorm.io/gorm"
	"nft_transfer/internal/conf"
	"os"
	"strconv"
)

// ProviderSet is data providers.
// var ProviderSet = wire.NewSet(NewData, NewGreeterRepo, NewNftTransferRepo)
var ProviderSet = wire.NewSet(NewData, NewDataBase, NewGreeterRepo, NewNftTransferRepo)

// Data .
type Data struct {
	RedisCli *redis.Client
	//DataBaseCli driver.Conn
	DataBaseCli *sdk.Gateway
}

// NewData .
func NewData(c *conf.Data, logger log.Logger, dataBaseCli *sdk.Gateway, redisCli *redis.Client) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{
		DataBaseCli: dataBaseCli,
		RedisCli:    redisCli,
	}, cleanup, nil
}
func NewDataBase(c *conf.Data, logger log.Logger) (*sdk.Gateway, error) {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = c.Database.Host
	}

	region := os.Getenv("DB_REGION")
	if region == "" {
		region = c.Database.Region
	}

	account := os.Getenv("DB_ACCOUNT")
	if account == "" {
		account = c.Database.Account
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = c.Database.User
	}
	password := os.Getenv("DB_PWD")
	if password == "" {
		password = c.Database.Password
	}
	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = c.Database.Dbname
	}

	//dsn := fmt.Sprintf("tcp://%s?region=%s&account=%s&user=%s&password=%s&secure=true&database=%s", host, region, account, user, password, dbname)
	dsn := fmt.Sprintf("tcp://%s?account=%s&user=%s&password=%s&secure=true&database=%s", host, account, user, password, dbname)

	db, err := sdk.Open(context.Background(), dsn)
	if err != nil {
		fmt.Printf("error = %v", err)
		return nil, nil
	}

	// Note first return value is sql.Result, which can be discarded since it is not implemented in the driver
	/*
		if _, err := db.Query("set default warehouse yiko_test_xl;"); err != nil {
			log.NewHelper(logger).Errorf("set default warehouse error", err)
		}*/

	defer db.Close()

	if err != nil {
		log.NewHelper(logger).Errorf("Failed to connect to database", err)
		return nil, err
	}
	log.NewHelper(logger).Info("Connected to DataBase!")
	return db, nil
}

func NewRedis(c *conf.Data, logger log.Logger) (*redis.Client, func(), error) {
	return nil, nil, nil

	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = c.Redis.Addr
	}
	db, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
	if db == 0 {
		db = int(c.Redis.Db)
	}
	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     "", // no password set
		DB:           db, // use default DB
		PoolSize:     15,
		MinIdleConns: 10,
	})

	// Check the connection
	err := client.Ping().Err()
	if err != nil {
		log.NewHelper(logger).Errorf("Failed to connect to Redis", err)
		return nil, nil, err
	}

	log.NewHelper(logger).Info("Connected to Redis!")
	return client, cleanup, nil
}
