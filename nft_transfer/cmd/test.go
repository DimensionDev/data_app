package main

import (
	"context"
	"fmt"

	"github.com/bytehouse-cloud/driver-go/sdk"
)

func main() {

	//dsn := fmt.Sprintf("tcp://%s?region=%s&account=%s&user=%s&password=%s&secure=true&database=%s", host, region, account, user, password, dbname)

	dsn := "tcp://gateway.aws-us-east-1.bytehouse.cloud:19000?account=AWSCOPZD&user=ksam&password=P@55word&secure=true&database=data_warehouse"
	g, err := sdk.Open(context.Background(), dsn)
	if err != nil {
		panic(err)
	}

	if err := g.Ping(); err != nil {
		panic(err)
	}

	str_sql := "select " +
		"nft_id," +
		"chain, " +
		"collection_id," +
		"transaction_initiator," +
		"block_number," +
		"block_hash," +
		"transaction_hash," +
		"block_timestamp," +
		"event_type," +
		"log_index," +
		"batch_transfer_index," +
		"address_from," +
		"address_to," +
		"quantity," +
		"sale_details" +
		" from transfer_nft_filter limit 10 "

	fmt.Println(str_sql)
	res, err := g.Query(str_sql)

	fmt.Println(err)

	fmt.Println(res)

	for {
		row, ok := res.NextRow()
		if !ok {
			break
		}
		fmt.Println(row)
	}
}
