package data

import (
	"context"
	"fmt"
	"github.com/bytehouse-cloud/driver-go/sdk"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"time"

	//"gorm.io/gorm"
	"middle_platform/internal/conf"
	"os"
	"strconv"
)

// ProviderSet is data providers.
// var ProviderSet = wire.NewSet(NewData, NewGreeterRepo, NewNftTransferRepo)
var ProviderSet = wire.NewSet(NewData, NewDataBase, NewRedis, NewGreeterRepo, NewNftTransferRepo, NewRateRepo)

// Data .
type Data struct {
	RedisCli *redis.Client
	//DataBaseCli driver.Conn
	// DataBaseCli *sdk.Gateway
	dc  *conf.Data
	DSN string
}

// NewData .
func NewData(c *conf.Data, logger log.Logger, dsn string, redisCli *redis.Client) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{
		// DataBaseCli: dataBaseCli,
		DSN:      dsn,
		RedisCli: redisCli,
		dc:       c,
	}, cleanup, nil
}

func NewDataBase(c *conf.Data, logger log.Logger) (string, error) {
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

	// fmt.Println("dsn: ", dsn)
	return dsn, nil
}

// func (r *Data) Conn() (*sdk.Gateway, error) {
// 	//dsn := fmt.Sprintf("tcp://%s?region=%s&account=%s&user=%s&password=%s&secure=true&database=%s", host, region, account, user, password, dbname)
// 	dsn := fmt.Sprintf("tcp://%s?account=%s&user=%s&password=%s&secure=true&database=%s", r.dc.Database.Host, r.dc.Database.Account, r.dc.Database.User, r.dc.Database.Password, r.dc.Database.Dbname)

// 	db, err := sdk.Open(context.Background(), dsn)
// 	if err != nil {
// 		fmt.Printf("error = %v", err)
// 		return nil, nil
// 	}
// 	if err != nil {
// 		//log.NewHelper(logger).Errorf("Failed to connect to database", err)
// 		fmt.Print("Failed to connect to database", err)
// 		return nil, err
// 	}
// 	//log.NewHelper(logger).Info("Connected to DataBase!")
// 	fmt.Print("Connected to DataBase!\n")
// 	return db, nil
// }

func (r *Data) data_query(str_sql string) (*sdk.QueryResult, error) {

	// if err := r.DataBaseCli.Ping(); err != nil {
	// 	// re-establish connection
	// 	db, connErr := r.Conn()
	// 	if db == nil {
	// 		return nil, connErr
	// 	}
	// 	_ = r.DataBaseCli.Close()
	// 	r.DataBaseCli = db
	// }
	fmt.Println("before connect to db:", time.Now())
	db, err := sdk.Open(context.Background(), r.DSN)
	if err != nil {
		fmt.Println("connect to db fail:", err)
		return nil, err
	}
	fmt.Println("after connected to db:", time.Now())

	start_time := time.Now().UnixMilli()
	fmt.Println("before query:", start_time)
	res, qerr := db.Query(str_sql)
	end_time := time.Now().UnixMilli()
	fmt.Println("after query:", end_time)
	use_time := fmt.Sprintf("query duration: %d(ms)", end_time-start_time)
	fmt.Println(use_time)

	return res, qerr
}

func NewRedis(c *conf.Data, logger log.Logger) (*redis.Client, func(), error) {

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
	// err := client.Ping(client.Context()).Err()
	// if err != nil {
	// 	log.NewHelper(logger).Errorf("Failed to connect to Redis", err)
	// 	return nil, nil, err
	// }

	log.NewHelper(logger).Info("Connected to Redis!")
	return client, cleanup, nil
}
