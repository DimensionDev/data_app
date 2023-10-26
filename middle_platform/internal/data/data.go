package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"

	bytehouse "github.com/bytehouse-cloud/driver-go"
	// "gorm.io/gorm/logger"

	// "github.com/bytehouse-cloud/driver-go/sdk"
	"time"

	_ "github.com/bytehouse-cloud/driver-go/sql"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"

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
	dc          *conf.Data
	DataBaseCli *sql.DB
}

// NewData .
func NewData(c *conf.Data, logger log.Logger, db_pool *sql.DB, redisCli *redis.Client) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
		db_pool.Close()
	}
	return &Data{
		DataBaseCli: db_pool,
		RedisCli:    redisCli,
		dc:          c,
	}, cleanup, nil
}

func NewDataBase(c *conf.Data, logger log.Logger) (*sql.DB, error) {
	host := os.Getenv("BYTEHOUSE_DB_HOST")
	if host == "" {
		host = c.Database.Host
	}

	port := os.Getenv("BYTEHOUSE_DB_PORT")
	if port == "" {
		port = strconv.Itoa(int(c.Database.Port))
	}

	apiToken := os.Getenv("BYTEHOUSE_DB_API_TOKEN")
	if apiToken == "" {
		apiToken = c.Database.ApiToken
	}

	dbname := os.Getenv("BYTEHOUSE_DB_NAME")
	if dbname == "" {
		dbname = c.Database.Dbname
	}

	dsn := fmt.Sprintf("tcp://%s:%s?secure=true&user=bytehouse&password=%s&database=%s", host, port, apiToken, dbname)
	//dsn := fmt.Sprintf("tcp://%s?region=%s&account=%s&user=%s&password=%s&secure=true&database=%s", host, region, account, user, password, dbname)
	// dsn := fmt.Sprintf("tcp://%s?account=%s&user=%s&password=%s&secure=true&database=%s", host, account, user, password, dbname)

	// fmt.Println("dsn: ", dsn)
	pool, err := sql.Open("bytehouse", dsn)
	if err != nil {
		log.NewHelper(logger).Error("create bytehouse connection pool failed:", err)
		return nil, errors.New("create bytehouse connection pool failed")
	}
	pool.SetMaxOpenConns(200)
	pool.SetConnMaxIdleTime(time.Minute)
	pool.SetConnMaxLifetime(time.Minute * 5)
	pool.SetMaxIdleConns(10)
	return pool, nil
}

func (r *Data) data_query(str_sql string) (*sql.Rows, error) {
	queryCtx := bytehouse.NewQueryContext(context.Background())
	//set the query ID here, duplicate query IDs will be rejected
	query_id := fmt.Sprintf("firefly-%v", uuid.New().String())
	fmt.Println("query ID:", query_id)
	queryCtx.SetQueryID(query_id)
	if err := queryCtx.AddQuerySetting("max_block_size", "2000"); err != nil {
		log.Error("query_id %v failed to add query setting err = %v", query_id, err)
		return nil, err
	}

	if err := r.DataBaseCli.Ping(); err != nil {
		log.Error("failed to ping err = %v", err)
		return nil, err
	}

	start_time := time.Now().UnixMilli()
	rows, qerr := r.DataBaseCli.QueryContext(queryCtx, str_sql)
	end_time := time.Now().UnixMilli()
	use_time := fmt.Sprintf("query duration: %d(ms)", end_time-start_time)
	fmt.Println(use_time)

	return rows, qerr
}

func (r *Data) data_query_single(str_sql string) (*sql.Row, error) {
	queryCtx := bytehouse.NewQueryContext(context.Background())
	//set the query ID here, duplicate query IDs will be rejected
	query_id := fmt.Sprintf("firefly-%v", uuid.New().String())
	fmt.Println("query ID:", query_id)
	queryCtx.SetQueryID(query_id)
	if err := queryCtx.AddQuerySetting("max_block_size", "2000"); err != nil {
		log.Error("query_id %v failed to add query setting err = %v", query_id, err)
		return nil, err
	}

	if err := r.DataBaseCli.Ping(); err != nil {
		log.Error("failed to ping err = %v", err)
		return nil, err
	}

	start_time := time.Now().UnixMilli()
	row := r.DataBaseCli.QueryRowContext(queryCtx, str_sql)
	end_time := time.Now().UnixMilli()
	use_time := fmt.Sprintf("query duration: %d(ms)", end_time-start_time)
	fmt.Println(use_time)

	return row, nil
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
