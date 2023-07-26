package main

import (
	"context"
	"fmt"

	"github.com/bytehouse-cloud/driver-go/sdk"
)

func main() {

	//dsn := fmt.Sprintf("tcp://%s?region=%s&account=%s&user=%s&password=%s&secure=true&database=%s", host, region, account, user, password, dbname)

	dsn := "tcp://gateway.aws-us-east-1.bytehouse.cloud:19000?account=AWS6R1TJ&user=jewen.han&password=Hwj521138=&secure=true&database=data_warehouse"
	g, err := sdk.Open(context.Background(), dsn)
	if err != nil {
		panic(err)
	}

	if err := g.Ping(); err != nil {
		panic(err)
	}

	str_sql := "select  chain, transaction_initiator,transaction_hash,block_timestamp,event_type,log_index,contract_address,token_id,address_from,address_to,owner "
	str_sql += "from transfer_nft_filter where  batch_transfer_index = 0  and  owner in ('0x63c8e1155e2be1e10041ff56625805abcf1fbf9b') and chain='polygon' and event_type='mint' order by block_timestamp desc limit  1,4 "

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
