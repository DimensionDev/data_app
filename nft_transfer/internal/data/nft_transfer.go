package data

import (
	"context"
	"errors"
	"fmt"
	"github.com/bytehouse-cloud/driver-go/sdk"
	"github.com/go-kratos/kratos/v2/log"
	pb "nft_transfer/api/nft_transfer/v1"
	"nft_transfer/internal/biz"
	"strings"
	"time"
)

type NftTransferRepo struct {
	data *Data
	log  *log.Helper
}

type NftTransfertmpSt struct {
	contract_address string
	network          string
	init_address     string
	event_type       string
	hash             string
	owner            string
	tag              string
	timestamp        string
	actios           map[string]DataActionST
}

type DataActionST struct {
	address_to   string
	event_type   string
	tag          string
	address_from string
	index        uint32
	token_id     string
}

// NewNftTransferRepo .
func NewNftTransferRepo(data *Data, logger log.Logger) biz.NftTransferRepo {
	return &NftTransferRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *NftTransferRepo) GetHandleNftinfo(ctx context.Context, req *pb.GetNftTransferRequest) (*pb.GetNftTransferReply, error) {

	handles, total, err := r.GetHandleNftinfoFromDB(r.data.DataBaseCli, req)

	if handles == nil {
		return &pb.GetNftTransferReply{
			Code:    500,
			Reason:  "EROR",
			Message: "query data fail",
			Data:    nil,
		}, err
	}

	var data pb.PnftTransferSt
	for _, nvalue := range handles {
		var node pb.NodeStArr
		node.AddressFrom = nvalue.init_address
		node.Network = nvalue.network
		node.AddressTo = nvalue.contract_address
		node.Type = nvalue.event_type
		node.Hash = nvalue.hash
		node.Owner = nvalue.owner
		node.Tag = nvalue.tag
		node.Timestamp = nvalue.timestamp

		for _, cvalue := range nvalue.actios {
			var action pb.ActionStArr
			action.Tag = cvalue.tag
			action.AddressTo = cvalue.address_to
			action.Type = cvalue.event_type
			action.Index = cvalue.index
			action.AddressFrom = cvalue.address_from

			if node.Network == "ethereum" || node.Network == "gnosis" {
				if node.AddressTo == "0x22c1f6050e56d2876009903609a2cc3fef83b415" {
					node.Type = "poap"
					action.Type = "poap"
				}
			}
			node.Actions = append(node.Actions, &action)
		}
		data.Result = append(data.Result, &node)
	}
	data.Total = total
	data.Cursor = req.Cursor + req.Limit

	if data.Cursor >= data.Total {
		data.Cursor = data.Total
	}
	//fmt.Println(data)

	return &pb.GetNftTransferReply{
		Code:    200,
		Reason:  "SUCCESS",
		Message: "SUCCESS",
		Data:    &data,
	}, err

}

func (r *NftTransferRepo) GetTotalFromDB(db *sdk.Gateway, total_sql string, ch *chan uint64) error {
	res1, err1 := db.Query(total_sql)
	var total uint64
	total = 0
	if err1 == nil {
		row1, ok1 := res1.NextRow()
		if ok1 {
			//u_total, _ := strconv.ParseUint(row1[0].(string), 10, 64)
			total = row1[0].(uint64)
		}
	}
	fmt.Print("go func in err ")
	*ch <- total
	fmt.Print("go func in err 1111 ")
	return nil
}

func (r *NftTransferRepo) GetHandleNftinfoFromDB(db *sdk.Gateway, req *pb.GetNftTransferRequest) (map[string]NftTransfertmpSt, uint64, error) {

	//nftlist := make([]*pb.PnftTransferSt, 5, 5)
	if req.Address == "" {
		return nil, 0, errors.New("input address is empty")
	}

	owners := strings.Split(req.Address, ",")

	str_where := "where  batch_transfer_index = 0  and  owner in ('"
	for i, owner := range owners {
		str_where += owner
		if i == len(owners)-1 {
			break
		}
		str_where += "','"
	}
	str_where += "')"

	if req.Network != "" {

		if req.Network == "binance_smart_chain" {
			req.Network = "bsc"
		}
		str_where += " and chain='" + req.Network + "'"
	}

	if req.Type != "" {
		if req.Type == "trade" {
			req.Type = "sale"
		}
		str_where += " and event_type='" + req.Type + "'"
	}

	//fmt.Print("order by:", req.OrderBy, req.OrderDirection)
	str_order := ""
	if req.OrderBy != "" {
		str_order += " order by " + req.OrderBy
		if req.OrderDirection != "" {
			str_order += " " + req.OrderDirection
		}
	} else {
		str_order += " order by block_timestamp "
	}

	str_limit := ""

	limit_n := req.Limit
	cursor_n := req.Cursor
	if req.Limit <= 0 {
		limit_n = 100
		if req.Cursor <= 0 {
			cursor_n = 0
		}
	}

	if limit_n > 1000 {
		limit_n = 1000
	}

	req.Limit = limit_n
	req.Cursor = cursor_n

	str_limit += fmt.Sprintf(" limit  %d,%d", cursor_n, limit_n+cursor_n)

	str_sql_p := "select  " +
		"chain, " +
		"transaction_initiator," +
		"transaction_hash," +
		"block_timestamp," +
		"event_type," +
		"log_index," +
		"contract_address," +
		"token_id," +
		"address_from," +
		"address_to," +
		"owner " +
		"from transfer_nft_filter "
	str_sql_p += str_where + str_order + str_limit
	fmt.Print("str_sql:", str_sql_p, "\n")
	var total uint64
	total = 0
	/*
		total_sql := "select count(distinct(chain, transaction_hash, log_index) ) from transfer_nft_filter  " + str_where

		fmt.Print("totalsql:", total_sql, "\n")

		res1, err1 := r.data.data_query(total_sql)
		var total uint64
		total = 0
		if err1 != nil {
			log.Errorf("query aaatotal  error", err1)
			return nil, total, err1
		}

		row1, ok1 := res1.NextRow()
		if ok1 {
			total = row1[0].(uint64)
		} else {
			log.Errorf("query total  error", row1)
			return nil, total, nil
		}*/

	//ch := make(chan uint64)

	//defer close(ch)

	//go r.GetTotalFromDB(db, total_sql, &ch)

	res, err := r.data.data_query(str_sql_p)

	if err != nil {
		//fmt.Print("err fail 666666666666666")
		//ttt := <-ch
		//fmt.Print("err fail 77777777777", ttt)
		return nil, 0, err
	}

	var data_nodes map[string]NftTransfertmpSt
	data_nodes = make(map[string]NftTransfertmpSt)
	for {
		row, ok := res.NextRow()
		if !ok {
			break
		}
		//fmt.Println(row)

		var node NftTransfertmpSt
		if row[0] != nil {
			node.network = row[0].(string)
			if node.network == "bsc" {
				node.network = "binance_smart_chain"
			}
		} else {
			node.network = ""
		}

		if row[1] != nil {
			node.init_address = row[1].(string)
		} else {
			node.init_address = ""
		}

		if row[2] != nil {
			node.hash = row[2].(string)
		} else {
			node.hash = ""
		}

		if row[3] != nil {
			node.timestamp = row[3].(time.Time).String()
		} else {
			node.timestamp = ""
		}
		if row[4] != nil {
			node.event_type = row[4].(string)
		} else {
			node.event_type = ""
		}

		if row[6] != nil {
			node.contract_address = row[6].(string)
		} else {
			node.contract_address = ""
		}

		if row[10] != nil {
			node.owner = row[10].(string)
		}

		node.tag = "collectible"
		node.actios = make(map[string]DataActionST)

		node_ukey := node.network + node.init_address + node.hash + node.owner
		var action DataActionST
		if row[8] != nil {
			action.address_from = row[8].(string)
		} else {
			action.address_from = "0x0000000000000000000000000000000000000000"
		}
		if row[9] != nil {
			action.address_to = row[9].(string)
		} else {
			action.address_to = "0x0000000000000000000000000000000000000000"
		}

		action.tag = "collectible"
		if row[4] != nil {
			action.event_type = row[4].(string)
		} else {
			action.event_type = ""
		}

		if row[5] != nil {
			action.index = row[5].(uint32)
		} else {
			action.index = 0
		}

		if row[7] != nil {
			action.token_id = row[7].(string)
		} else {
			action.token_id = ""
		}

		if action.event_type == "burn" {
			action.address_to = "0x0000000000000000000000000000000000000000"
		}

		if action.event_type == "sale" {
			action.event_type = "trade"
		}

		action_ukey := node.contract_address + action.token_id + action.event_type

		action_ukey += fmt.Sprintf("%d", action.index)

		if _, ok := data_nodes[node_ukey]; ok {

			if _, ok := data_nodes[node_ukey].actios[action_ukey]; ok {

			} else {
				data_nodes[node_ukey].actios[action_ukey] = action
			}

		} else {
			node.actios[action_ukey] = action
			data_nodes[node_ukey] = node
		}

	}

	// Return an error if no data is found
	if len(data_nodes) == 0 {
		//fmt.Print("chanbnel aaaaaaaaaaaaaa\n")
		//bbbb := <-ch
		//fmt.Print("chanbnel aaaaaaaaaaaaaa", bbbb)
		return nil, 0, errors.New("no data in database ")
	}

	//fmt.Print("chanbnel 11111111111111111111111111111\n")
	//utotal := <-ch
	//fmt.Print("chan value", utotal)

	return data_nodes, total, errors.New("success")
}
