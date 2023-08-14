package data

import (
	"context"

	"encoding/json"

	"errors"

	"fmt"

	"github.com/bytehouse-cloud/driver-go/sdk"

	"github.com/go-kratos/kratos/v2/log"

	pb "middle_platform/api/nft_transfer/v1"

	"middle_platform/internal/biz"

	"strings"

	"sort"
	"strconv"
	"time"
)

type NftTransferRepo struct {
	data *Data

	log *log.Helper
}

type NftTransfertmpSt struct {
	contract_address string

	network string

	init_address string

	event_type string

	hash string

	owner string

	sale_details string

	tag string

	timestamp string

	actios map[string]DataActionST
}

type DataActionST struct {
	address_to string

	event_type string

	tag string

	address_from string

	index uint32

	token_id string
}

type PaymentToken struct {
	Payment_token_id string `json:"payment_token_id"`

	Name string `json:"name"`

	Symbol string `json:"symbol"`

	Address string `json:"address"`

	Decimals uint32 `json:"decimals"`
}

type SaleInfo struct {
	Marketplace_id string `json:"marketplace_id"`

	Marketplace_name string `json:"marketplace_name"`

	Is_bundle_sale bool `json:"is_bundle_sale"`

	Payment_token PaymentToken `json:"payment_token"`

	Unit_price json.Number `json:"unit_price"`

	Total_price json.Number `json:"total_price"`
}

// NewNftTransferRepo .

func NewNftTransferRepo(data *Data, logger log.Logger) biz.NftTransferRepo {

	return &NftTransferRepo{

		data: data,

		log: log.NewHelper(logger),
	}

}

func (r *NftTransferRepo) GetHandleNftinfo(ctx context.Context, req *pb.GetNftTransferRequest) (*pb.GetNftTransferReply, error) {

	handles, action_num, err := r.GetHandleNftinfoFromDB(r.data.DataBaseCli, req)
	if err != nil {
		return &pb.GetNftTransferReply{

			Code: 500,

			Reason: err.Error(),

			Message: err.Error(),

			Data: nil,
		}, err
	}

	if handles == nil {

		return &pb.GetNftTransferReply{

			Code: 200,

			Reason: "",

			Message: "",

			Data: nil,
		}, err

	}

	var data pb.PnftTransferSt

	for _, nvalue := range handles {
		var node pb.NodeStArr

		node.AddressFrom = nvalue.init_address

		node.Network = nvalue.network

		node.AddressTo = nvalue.contract_address

		node.Hash = nvalue.hash

		node.Owner = nvalue.owner

		node.Tag = nvalue.tag

		node.Timestamp = nvalue.timestamp

		node.Type = nvalue.event_type
		if node.Type == "sale" {
			node.Type = "trade"
		}

		if node.Network == "ethereum" || node.Network == "gnosis" {
			if node.AddressTo == "0x22c1f6050e56d2876009903609a2cc3fef83b415" {
				node.Type = "poap"
			}
		}

		for _, cvalue := range nvalue.actios {

			var action pb.ActionStArr
			var sale_info SaleInfo
			if &nvalue.sale_details != nil && nvalue.sale_details != "" {
				err := json.Unmarshal([]byte(nvalue.sale_details), &sale_info)

				if err != nil {
					fmt.Println("sale details parsing error:", err)

				} else {
					var cost pb.CostSt
					cost.Symbol = sale_info.Payment_token.Symbol
					// blur pool token translate to ETH for price searing in front-end
					if sale_info.Payment_token.Payment_token_id == "ethereum.0x0000000000a39bb272e79075ade125fd351887ac" {
						cost.Symbol = "ETH"
					}
					cost.Value = sale_info.Total_price.String()
					cost.Decimals = sale_info.Payment_token.Decimals
					action.Cost = &cost
				}
			}

			action.ContractAddress = nvalue.contract_address
			action.TokenId = cvalue.token_id
			action.Tag = cvalue.tag
			action.AddressTo = cvalue.address_to
			action.Type = node.Type
			action.Index = cvalue.index
			action.AddressFrom = cvalue.address_from
			node.Actions = append(node.Actions, &action)
		}

		data.Result = append(data.Result, &node)
		sort.Slice(data.Result, func(i int, j int) bool {
			return data.Result[i].Timestamp > data.Result[j].Timestamp
		})

	}

	if req.Limit == action_num {

		str := strconv.FormatUint(uint64(req.Cursor+req.Limit), 10)
		data.Cursor = &str

	}

	return &pb.GetNftTransferReply{

		Code: 200,

		Reason: "SUCCESS",

		Message: "SUCCESS",

		Data: &data,
	}, err

}

func (r *NftTransferRepo) GetHandleNftinfoFromDB(db *sdk.Gateway, req *pb.GetNftTransferRequest) (map[string]NftTransfertmpSt, uint64, error) {

	//nftlist := make([]*pb.PnftTransferSt, 5, 5)

	if req.Address == "" {

		return nil, 0, errors.New("input address is empty")

	}

	owners := strings.Split(req.Address, ",")

	str_where := "where batch_transfer_index = 0 and owner in ('"

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

		if !strings.Contains(strings.ToLower(req.Network), "all") {
			str_where += " and chain='" + req.Network + "'"
		}

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

		if req.OrderDirection == "" {

			req.OrderDirection = "desc"
		}

		str_order += " " + req.OrderDirection

	} else {

		str_order += " order by block_timestamp desc"

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

	str_limit += fmt.Sprintf(" limit %d,%d", req.Cursor, req.Limit)

	group_by := ""
	if req.OrderBy != "" {
		group_by += " group by chain,transaction_hash,owner,event_type," + req.OrderBy
	} else {
		group_by += " group by chain,transaction_hash,owner,event_type,block_timestamp"
	}

	re_filter_str := " match(name, '^(([1-9][0-9]{3}\\$)|(\\$[1-9][0-9]{3})) [a-zA-Z]+') "
	collection_sub_query := " (select collection_id from collections where spam_score>=80 or name like '%.lens-Follower' or " + re_filter_str + ") "
	spam_filter_condition := " and contract_address not in (select contract from spam_contracts) and collection_id not in  " + collection_sub_query
	first_q := "select chain,transaction_hash,owner,event_type,block_timestamp from transfer_nft_filter_index " + str_where + spam_filter_condition + group_by + str_order + str_limit
	fmt.Println("first_q:", first_q)
	first_res, err := r.data.data_query(first_q)
	if err != nil {
		return nil, 0, err
	}

	var chains []string
	var hashs []string
	var _owners []string
	var event_types []string
	var dup_keys map[string]bool
	dup_keys = make(map[string]bool)
	var action_num uint64 = 0

	for {
		row, ok := first_res.NextRow()
		if !ok {
			break
		}
		action_num += 1
		chains = append(chains, row[0].(string))
		hashs = append(hashs, row[1].(string))
		_owners = append(_owners, row[2].(string))
		event_types = append(event_types, row[3].(string))
		dups := []string{row[0].(string), row[1].(string), row[2].(string), row[3].(string)}
		dup_string := strings.Join(dups, "_")
		dup_keys[dup_string] = true
	}

	chain_condition := combineAndRemoveDuplicates("chain", chains)
	hash_condition := combineAndRemoveDuplicates("transaction_hash", hashs)
	owner_condition := combineAndRemoveDuplicates("owner", _owners)
	event_type_condition := combineAndRemoveDuplicates("event_type", event_types)

	conditions := []string{chain_condition, hash_condition, owner_condition, event_type_condition}
	combine_in_condition := strings.Join(conditions, " and ")
	where_str := "where" + combine_in_condition

	str_sql_p := "select " +
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
		"owner," +
		"sale_details " +
		"from transfer_nft_filter_index "

	str_sql_p += where_str

	fmt.Println("str_sql:", str_sql_p)

	res, err := r.data.data_query(str_sql_p)

	if err != nil {
		return nil, 0, err

	}

	var data_nodes map[string]NftTransfertmpSt

	data_nodes = make(map[string]NftTransfertmpSt)

	for {

		row, ok := res.NextRow()
		if !ok {
			break
		}

		// 这里要去除查出来的多余的记录
		dups := []string{row[0].(string), row[2].(string), row[10].(string), row[4].(string)}
		dup_string := strings.Join(dups, "_")
		_, ok = dup_keys[dup_string]
		if !ok {
			continue
		}

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

		const targetLayout = "2006-01-02T15:04:05Z"
		if row[3] != nil {

			node.timestamp = row[3].(time.Time).Format(targetLayout)

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

		if row[11] != nil {
			node.sale_details = row[11].(string)
		}

		node.tag = "collectible"

		node.actios = make(map[string]DataActionST)

		node_ukey := node.network + node.hash + node.owner + node.event_type

		var action DataActionST

		if row[8] != nil && row[8] != "" {

			action.address_from = row[8].(string)

		} else {

			action.address_from = "0x0000000000000000000000000000000000000000"

		}

		if row[9] != nil && row[9] != "" {

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

		if action.event_type == "sale" {

			action.event_type = "trade"

		}

		action_ukey := strconv.FormatUint(uint64(action.index), 10)

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

	if len(data_nodes) == 0 {

		//fmt.Print("chanbnel aaaaaaaaaaaaaa\n")

		//bbbb := <-ch

		//fmt.Print("chanbnel aaaaaaaaaaaaaa", bbbb)

		return nil, action_num, nil

	}

	//fmt.Print("chanbnel 11111111111111111111111111111\n")

	//utotal := <-ch

	//fmt.Print("chan value", utotal)

	return data_nodes, action_num, nil

}

func combineAndRemoveDuplicates(field string, strArr []string) string {
	elements := make(map[string]bool)
	var result strings.Builder

	for _, str := range strArr {
		if !elements[str] {
			elements[str] = true
			if result.Len() > 0 {
				result.WriteString(",")
			}
			result.WriteString("'" + str + "'")
		}
	}
	left_str := " " + field + " in ("
	combine_str := left_str + result.String() + ")"
	return combine_str
}
