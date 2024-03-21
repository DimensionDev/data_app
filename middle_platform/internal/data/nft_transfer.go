package data

import (
	"context"
	"encoding/json"

	"errors"

	"fmt"

	// "github.com/bytehouse-cloud/driver-go/sdk"

	"github.com/go-kratos/kratos/v2/log"

	pb "middle_platform/api/nft_transfer/v1"

	"middle_platform/internal/biz"

	resty "github.com/go-resty/resty/v2"

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

type transaction struct {
	chain            string
	transaction_hash string
	owner            string
	event_type       string
	block_timestamp  time.Time
}

const ZERO_ADDRESS string = "0x0000000000000000000000000000000000000000"

// NewNftTransferRepo .

func NewNftTransferRepo(data *Data, logger log.Logger) biz.NftTransferRepo {

	return &NftTransferRepo{

		data: data,

		log: log.NewHelper(logger),
	}

}

func (r *NftTransferRepo) GetHandleNftinfo(ctx context.Context, req *pb.GetNftTransferRequest) (*pb.GetNftTransferReply, error) {

	handle_start_time := time.Now().UnixMilli()
	handles, action_num, err := r.GetHandleNftinfoFromDB(req)
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
	fmt.Println("action_num:", action_num)
	if req.Limit == action_num {

		str := strconv.FormatUint(uint64(req.Cursor+req.Limit), 10)
		data.Cursor = &str

	}

	handle_end_time := time.Now().UnixMilli()
	fmt.Println("handle time:", handle_end_time-handle_start_time)
	return &pb.GetNftTransferReply{

		Code: 200,

		Reason: "SUCCESS",

		Message: "SUCCESS",

		Data: &data,
	}, err

}

func containsString(collection []string, target string) bool {
	for _, s := range collection {
		if s == target {
			return true // 如果找到目标字符串，返回true
		}
	}
	return false // 如果遍历完仍未找到目标字符串，返回false
}

func (r *NftTransferRepo) PostSpamReport(ctx context.Context, req *pb.PostReportSpamRequest) (*pb.PostReportSpamReply, error) {
	//判断状态
	collection_id := req.CollectionId
	next_status := req.Status
	req_source := req.Source
	var source string

	// 查找先前的 report 记录
	query_str := fmt.Sprintf("select status,create_at,source from spam_report where collection_id = '%s'", collection_id)
	res, err := r.data.data_query_single(query_str)
	if err != nil {
		fmt.Println("post spam report fail", collection_id, next_status)
		return nil, err
	}

	type report struct {
		status    string
		create_at time.Time
		source    string
	}
	var rt report
	if res != nil {
		if err := res.Scan(&rt.status, &rt.create_at, &rt.source); err != nil {
			fmt.Println(err.Error())
			if err.Error() == "sql: no rows in result set" {
				res = nil
			} else {
				fmt.Printf("error = %v", err)
				return nil, err
			}
		}
	}

	fmt.Println("rt:", rt)
	// row, ok := res.NextRow()
	// if !ok {
	// 	row = nil
	// }

	const targetLayout = "2006-01-02T15:04:05Z"
	if next_status == "reporting" {
		if req_source == nil {
			source = "firefly"
		} else {
			source = *req_source
			sources := []string{"firefly", "mask-network"}
			if !containsString(sources, source) {
				fmt.Println("source:", source)
				// return nil, errors.New(fmt.Sprintf("value of source field should be in %s", sources))
				return nil, fmt.Errorf("value of source field should be in %s", sources)
			}
		}
		// 检查 collection 是否已经被report
		if res != nil {
			if rt.status == next_status || rt.status == "rejected" {
				create_at := rt.create_at.Format(targetLayout)
				update_at := time.Now().UTC().Format(targetLayout)
				insert_str := fmt.Sprintf("insert into spam_report values ('%s','%s','%s','%s','%s')", collection_id, next_status, create_at, update_at, rt.source)
				insert_err := InsertIntoSpamReportTable(r, insert_str)
				if insert_err != nil {
					return nil, insert_err
				}
				data := pb.SpamReport{
					CollectionId: collection_id,
					Status:       next_status,
					CreateAt:     &create_at,
					UpdateAt:     &update_at,
					Source:       &source,
				}
				return &pb.PostReportSpamReply{
					Code:    200,
					Message: "",
					Data:    &data,
				}, nil
			} else {
				return &pb.PostReportSpamReply{
					Code:    400,
					Message: fmt.Sprintf("this report of %s is already %s", collection_id, rt.status),
					Data:    nil,
				}, nil
			}
		} else {
			create_at := time.Now().UTC().Format(targetLayout)
			update_at := create_at
			insert_str := fmt.Sprintf("insert into spam_report values ('%s','%s','%s','%s','%s')", collection_id, next_status, create_at, update_at, source)
			insert_err := InsertIntoSpamReportTable(r, insert_str)
			if insert_err != nil {
				return nil, insert_err
			}
			data := pb.SpamReport{
				CollectionId: collection_id,
				Status:       next_status,
				CreateAt:     &create_at,
				UpdateAt:     &update_at,
				Source:       &source,
			}
			return &pb.PostReportSpamReply{
				Code:    200,
				Message: "",
				Data:    &data,
			}, nil
		}
	} else if next_status == "approved" || next_status == "rejected" {
		if res != nil {
			fmt.Println(res)
			if rt.status == "reporting" {
				update_at := time.Now().UTC().Format(targetLayout)
				create_at := rt.create_at.Format(targetLayout)

				// 更新 collection 的 spam_score
				if next_status == "approved" {
					report_err := reportSpamToSimpleHash(collection_id)
					if report_err != nil {
						return &pb.PostReportSpamReply{
							Code:    500,
							Message: report_err.Error(),
							Data:    nil,
						}, nil
					}
					err := UpdataCollectionSpamScore(r, collection_id)
					if err != nil {
						return &pb.PostReportSpamReply{
							Code:    500,
							Message: err.Error(),
							Data:    nil,
						}, nil
					}
				}

				reply_source := rt.source
				insert_str := fmt.Sprintf("insert into spam_report values ('%s','%s','%s','%s','%s')", collection_id, next_status, create_at, update_at, reply_source)
				insert_err := InsertIntoSpamReportTable(r, insert_str)
				if insert_err != nil {
					return &pb.PostReportSpamReply{
						Code:    500,
						Message: insert_err.Error(),
						Data:    nil,
					}, nil
				}
				// 返回数据

				data := pb.SpamReport{
					CollectionId: collection_id,
					Status:       next_status,
					CreateAt:     &create_at,
					UpdateAt:     &update_at,
					Source:       &reply_source,
				}
				return &pb.PostReportSpamReply{
					Code:    200,
					Message: "",
					Data:    &data,
				}, nil
			} else {
				return &pb.PostReportSpamReply{
					Code:    400,
					Message: fmt.Sprintf("this report of %s is already %s", collection_id, rt.status),
					Data:    nil,
				}, nil
				// return nil, fmt.Errorf("the report of colleciont %s have beed %s", collection_id, row[0])
			}
		} else {
			return &pb.PostReportSpamReply{
				Code:    400,
				Message: fmt.Sprintf("no collection of %s be reported before", collection_id),
				Data:    nil,
			}, nil
			// return nil, fmt.Errorf("no collection of %s be reported before", collection_id)
		}
	} else {
		return &pb.PostReportSpamReply{
			Code:    400,
			Message: "value of status should be in ['reporting','approved','rejected']",
			Data:    nil,
		}, nil
		// return nil, fmt.Errorf("value of status should be in ['reporting','approved','rejected']")
	}
}

func InsertIntoSpamReportTable(r *NftTransferRepo, insert_str string) error {
	fmt.Println("insert str:", insert_str)
	_, err := r.data.data_query(insert_str)
	if err != nil {
		return fmt.Errorf("writing data into bytehouse error:%s", err)
	}
	return nil
}

func UpdataCollectionSpamScore(r *NftTransferRepo, collection_id string) error {
	update_str := fmt.Sprintf("insert into spam_collections_with_bucket values ('%s',100,'')", collection_id)
	// update_str := fmt.Sprintf("update collections_new_test set spam_score=100 where collection_id='%s'", collection_id)
	// update_str := fmt.Sprintf("insert into collections (collection_id,spam_score) values ('%s', 100)", "collection_id")
	fmt.Println("update_str:", update_str)
	_, err := r.data.data_query(update_str)
	if err != nil {
		return fmt.Errorf("writing data into bytehouse error:%s", err)
	}
	return nil
}

func reportSpamToSimpleHash(collection_id string) error {
	// curl --location 'https://api.simplehash.com/api/v0/nfts/report/spam' \ --header 'X-API-KEY: mask_sk_Wv1uXGWUVWHx7LAPOvKWHmSot0' \ --header 'accept: application/json' \ --header 'content-type: application/json' \ --data '{ "collection_id":"bf60f01b784a501dded3dd73f5347832", "event_type":"mark_as_spam" }'
	http_client := resty.New()

	resp, err := http_client.R().
		SetHeader("X-API-KEY", "mask_sk_Wv1uXGWUVWHx7LAPOvKWHmSot0").
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetBody(`{"collection_id": "` + collection_id + `", "event_type":"mark_as_spam"}`).
		Post("https://api.simplehash.com/api/v0/nfts/report/spam")
	if err != nil {
		fmt.Println("report to simplehash error:", err)
		return fmt.Errorf("the collection of %s report to simplehash fail", collection_id)
	}
	fmt.Println("report resp:", resp)
	return nil
}

func (r *NftTransferRepo) GetSpamReport(ctx context.Context, req *pb.GetReportSpamRequest) (*pb.GetReportSpamReply, error) {
	where_str := "where "
	condition_str := ""

	collection_id_str := ""
	if req.CollectionId != "" {
		collection_id_str = combineAndRemoveDuplicates("collection_id", strings.Split(req.CollectionId, ","))
		// collection_id_str = fmt.Sprintf(" collection_id in ('%s')", req.CollectionId)
	}
	status_str := ""
	if req.Status != "" {
		status_str = combineAndRemoveDuplicates("status", strings.Split(req.Status, ","))
	}

	source_str := ""
	if req.Source != "" {
		source_str = combineAndRemoveDuplicates("source", strings.Split(req.Source, ","))
	}

	if collection_id_str != "" && status_str != "" {
		condition_str = where_str + collection_id_str + " and " + status_str
	} else if collection_id_str != "" || status_str != "" {
		condition_str = where_str + collection_id_str + status_str
	}

	if condition_str != "" && source_str != "" {
		condition_str = condition_str + " and " + source_str
	} else if condition_str == "" && source_str != "" {
		condition_str = where_str + source_str
	}

	var page uint32
	if req.Page == 0 {
		page = uint32(1)
	} else {
		page = req.Page
	}

	// cursor_str := strconv.FormatUint(uint64(req.Cursor), 10)
	var limit uint32
	if req.Limit == uint32(0) {
		limit = uint32(100)
	} else {
		limit = req.Limit
	}

	cursor := (page - 1) * limit
	cursor_str := strconv.FormatUint(uint64(cursor), 10)
	limit_str := strconv.FormatUint(uint64(limit), 10)
	cursor_limit_str := "limit " + cursor_str + "," + limit_str

	order_str := "order by update_at desc"
	query_str := fmt.Sprintf(
		"select collection_id,status,create_at,update_at,source from spam_report %s %s %s",
		condition_str, order_str, cursor_limit_str)

	total_query_str := fmt.Sprintf("select count(1) from spam_report %s ", condition_str)
	fmt.Println("total_query_str:", total_query_str)
	total_count, err := r.GetTotalNumberOfSpamReport(total_query_str)
	if err != nil {
		return &pb.GetReportSpamReply{
			Code: 500,
			Data: nil,
		}, err
	}
	fmt.Println("total count:", total_count)
	// current_page := cursor / limit + 1

	fmt.Println("query str:", query_str)
	res, err := r.data.data_query(query_str)

	if err != nil {
		fmt.Println("query from bytehouse:", query_str)
		return &pb.GetReportSpamReply{
			Code: 500,
			Data: nil,
		}, err
	}

	var report_list []*pb.SpamReport

	// row := make([]interface{}, 0)
	const targetLayout = "2006-01-02T15:04:05Z"
	type report struct {
		collection_id string
		status        string
		created_at    time.Time
		update_at     time.Time
		source        string
	}
	var rt report
	for res.Next() {
		if err := res.Scan(&rt.collection_id, &rt.status, &rt.created_at, &rt.update_at, &rt.source); err != nil {
			log.Error("failed to scan row err = %v", err)
			return nil, err
		}
		var spam_report pb.SpamReport
		var create_at string
		var update_at string
		spam_report.CollectionId = rt.collection_id
		spam_report.Status = rt.status
		create_at = rt.created_at.Format(targetLayout)
		spam_report.CreateAt = &create_at
		update_at = rt.update_at.Format(targetLayout)
		spam_report.UpdateAt = &update_at
		reply_source := rt.source
		spam_report.Source = &reply_source
		report_list = append(report_list, &spam_report)
	}
	// reportSpamReply.Data = report_list
	var result pb.GetReportSpamReply
	result.Code = 200
	result.Limit = limit
	result.Data = report_list
	result.Page = page
	result.Total = total_count

	// var next_cursor uint32
	// if uint32(len(report_list)) == limit {
	// 	next_cursor = req.Cursor + limit
	// 	result.Cursor = &next_cursor
	// }

	return &result, nil
	// return &reportSpamReply, err
	// return nil, nil
}

func (r *NftTransferRepo) GetTotalNumberOfSpamReport(query_str string) (uint64, error) {
	res, err := r.data.data_query_single(query_str)
	if err != nil {
		fmt.Println("query total count for spam reports err:", err)
		return uint64(0), err
	}
	if res == nil {
		// if !ok {
		fmt.Println("query fail from bytehouse")
		return uint64(0), errors.New("query fail from bytehouse")
	} else {
		var count uint64
		if err := res.Scan(&count); err != nil {
			log.Error("failed to scan row err = %v", err)
			return 0, err
		}
		return count, nil
	}
}

func (r *NftTransferRepo) GetHandleNftinfoFromDB(req *pb.GetNftTransferRequest) (map[string]NftTransfertmpSt, uint64, error) {

	//nftlist := make([]*pb.PnftTransferSt, 5, 5)
	if req.Address == "" {

		return nil, 0, errors.New("input address is empty")

	}

	owners := strings.Split(req.Address, ",")

	//str_where := "where batch_transfer_index = 0 and owner in ('"
	str_where := "where owner in ('"

	for i, owner := range owners {

		str_where += owner

		if i == len(owners)-1 {

			break

		}

		str_where += "','"

	}

	str_where += "')"

	// if req.Network != "" {

	// 	if req.Network == "binance_smart_chain" {

	// 		req.Network = "bsc"

	// 	}
	if req.Network != "" {
		if !strings.Contains(strings.ToLower(req.Network), "all") {
			networks := strings.Split(req.Network, ",")
			networkCondition := combineAndRemoveDuplicates("chain", networks)
			str_where = str_where + " and " + networkCondition
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

	//re_filter_str := " match(name, '(^(([1-9][0-9]{3}\\\\$)|(\\\\$[1-9][0-9]{3})) [a-zA-Z]+)|(.*lens-Follower$)') "
	//collection_sub_query := " (select collection_id from collections_new_test where spam_score>=50 or " + re_filter_str + ") "
	//spam_filter_condition := " and collection_id not in  " + collection_sub_query
	// first_q := "select chain,transaction_hash,owner,event_type,block_timestamp from transfer_nft_filter_index " + str_where + spam_filter_condition + group_by + str_order + str_limit

	spam_filter_condition := " and collection_id not in (select collection_id from spam_collections_with_bucket ) "
	//first_q := "select chain,transaction_hash,owner,event_type,block_timestamp from transfer_nft_filter_new " + str_where + spam_filter_condition + group_by + str_order + str_limit
	first_q := "select chain,transaction_hash,owner,event_type,block_timestamp from nft_transfer_summary_selected_chains " + str_where + spam_filter_condition + str_order + str_limit
	fmt.Println("first_q:", first_q)
	first_res, err := r.data.data_query(first_q)
	if err != nil {
		return nil, 0, err
	}

	// release connection
	defer first_res.Close()

	var ts transaction
	var chains []string
	var hashs []string
	var _owners []string
	var event_types []string
	var dup_keys map[string]bool
	dup_keys = make(map[string]bool)
	var action_num uint64 = 0
	var timestamps []time.Time

	for first_res.Next() {
		if err := first_res.Scan(
			&ts.chain,
			&ts.transaction_hash,
			&ts.owner,
			&ts.event_type,
			&ts.block_timestamp); err != nil {
			log.Error("failed to scan row err = %v", err)
			return nil, 0, err
		}
		action_num += 1
		chains = append(chains, ts.chain)
		hashs = append(hashs, ts.transaction_hash)
		_owners = append(_owners, ts.owner)
		event_types = append(event_types, ts.event_type)
		timestamps = append(timestamps, ts.block_timestamp)
		dups := []string{ts.chain, ts.transaction_hash, ts.owner, ts.event_type}
		dup_string := strings.Join(dups, "_")
		dup_keys[dup_string] = true
	}

	if len(chains) == 0 && len(_owners) == 0 {
		return nil, 0, nil
	}

	maxTime := timestamps[0]
	minTime := timestamps[0]

	// 遍历数组，比较找到最大和最小时间
	for _, t := range timestamps {
		if t.After(maxTime) {
			maxTime = t
		}
		if t.Before(minTime) {
			minTime = t
		}
	}

	chain_condition := combineAndRemoveDuplicates("chain", chains)
	hash_condition := combineAndRemoveDuplicates("transaction_hash", hashs)
	owner_condition := combineAndRemoveDuplicates("owner", _owners)
	event_type_condition := combineAndRemoveDuplicates("event_type", event_types)

	conditions := []string{owner_condition, hash_condition, chain_condition, event_type_condition}
	combine_in_condition := strings.Join(conditions, " and ")
	maxTime_str := maxTime.Format("2006-01-02 15:04:05")
	minTime_str := minTime.Format("2006-01-02 15:04:05")
	time_condition := " block_timestamp >= '" + minTime_str + "' and block_timestamp <= '" + maxTime_str + "' and "
	combine_in_condition = time_condition + combine_in_condition
	// where_str := " prewhere " + owner_condition + " where" + combine_in_condition + " and batch_transfer_index=0"
	where_str := " prewhere " + owner_condition + " where" + combine_in_condition + " and batch_transfer_index=0"

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
		"from transfer_nft_filter_index_selected_chains"

	str_sql_p += where_str

	fmt.Println("str_sql:", str_sql_p)

	log_rows, err := r.data.data_query(str_sql_p)

	if err != nil {
		return nil, 0, err

	}

	type transaction_log struct {
		chain                 string
		transaction_initiator string
		transaction_hash      string
		block_timestamp       time.Time
		event_type            string
		log_index             uint32
		contract_address      string
		token_id              string
		address_from          *string
		address_to            *string
		owner                 string
		sale_details          *string
	}

	var data_nodes map[string]NftTransfertmpSt

	data_nodes = make(map[string]NftTransfertmpSt)

	for log_rows.Next() {
		var ts_log transaction_log
		if err := log_rows.Scan(
			&ts_log.chain,
			&ts_log.transaction_initiator,
			&ts_log.transaction_hash,
			&ts_log.block_timestamp,
			&ts_log.event_type,
			&ts_log.log_index,
			&ts_log.contract_address,
			&ts_log.token_id,
			&ts_log.address_from,
			&ts_log.address_to,
			&ts_log.owner,
			&ts_log.sale_details); err != nil {
			log.Error("failed to scan row err = %v", err)
			return nil, 0, err
		}

		// 这里要去除查出来的多余的记录
		dups := []string{ts_log.chain, ts_log.transaction_hash, ts_log.owner, ts_log.event_type}
		dup_string := strings.Join(dups, "_")
		_, ok := dup_keys[dup_string]
		if !ok {
			continue
		}

		var node NftTransfertmpSt

		node.network = ts_log.chain
		if node.network == "bsc" {
			node.network = "binance_smart_chain"
		}
		node.init_address = ts_log.transaction_initiator
		node.hash = ts_log.transaction_hash

		const targetLayout = "2006-01-02T15:04:05Z"
		node.timestamp = ts_log.block_timestamp.Format(targetLayout)
		node.event_type = ts_log.event_type
		node.contract_address = ts_log.contract_address
		node.owner = ts_log.owner

		if ts_log.sale_details != nil {
			node.sale_details = *ts_log.sale_details
		} else {
			node.sale_details = ""
		}

		node.tag = "collectible"

		node.actios = make(map[string]DataActionST)

		node_ukey := node.network + node.hash + node.owner + node.event_type

		var action DataActionST

		if ts_log.address_from != nil && *ts_log.address_from != "" {
			action.address_from = *ts_log.address_from
		} else {
			action.address_from = ZERO_ADDRESS
		}

		if ts_log.address_to != nil && *ts_log.address_to != "" {
			action.address_to = *ts_log.address_to
		} else {
			action.address_to = ZERO_ADDRESS
		}

		action.tag = "collectible"
		action.event_type = ts_log.event_type
		action.index = ts_log.log_index
		action.token_id = ts_log.token_id

		if action.event_type == "sale" {
			action.event_type = "trade"
		}

		action_ukey := strconv.FormatUint(uint64(action.index), 10)

		if _, ok := data_nodes[node_ukey]; ok {
			if data_nodes[node_ukey].timestamp != node.timestamp {
				actions := data_nodes[node_ukey].actios
				for old_action_ukey, exist_action := range actions {
					if exist_action.token_id == action.token_id && data_nodes[node_ukey].timestamp < node.timestamp {
						data_nodes[node_ukey].actios[old_action_ukey] = action
					}
				}

			} else {
				if _, ok := data_nodes[node_ukey].actios[action_ukey]; ok {
				} else {
					data_nodes[node_ukey].actios[action_ukey] = action
				}
			}
		} else {
			node.actios[action_ukey] = action
			data_nodes[node_ukey] = node
		}
	}

	if len(data_nodes) == 0 {
		return nil, action_num, nil
	}
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
func (r *NftTransferRepo) GetTransferNft(ctx context.Context, req *pb.GetTransferNftRequest) (*pb.GetTransferNftReply, error) {
	whereCondition := "where `block_timestamp` >='2023-07-31 00:00:00' "
	if req.ContractAddress != "" {
		whereCondition = whereCondition + "and `contract_address`='" + req.ContractAddress + "'"
	}
	eventTypeCondition := ""
	if req.EventType != "" {
		eventTypeCondition = combineAndRemoveDuplicates("`event_type`", strings.Split(req.EventType, ","))
		if eventTypeCondition != "" {
			whereCondition = whereCondition + " and " + eventTypeCondition
		}
	}
	ownerCondition := ""
	if req.Owners != "" {
		ownerCondition = combineAndRemoveDuplicates("`owner`", strings.Split(req.Owners, ","))
		if ownerCondition != "" {
			whereCondition = whereCondition + " and " + ownerCondition
		}
	}

	var page uint32
	if req.Page <= 0 {
		page = uint32(1)
	} else {
		page = req.Page
	}
	var limit uint32
	if req.Limit == uint32(0) {
		limit = uint32(100)
	} else {
		limit = req.Limit
	}
	query := fmt.Sprintf(
		"select nft_id,chain,contract_address,token_id,collection_id,event_type,address_from,address_to,block_timestamp,owner from transfer_nft %s order by create_time desc limit %d,%d",
		whereCondition, (page-1)*limit, limit)
	fmt.Println("query:", query)
	totalQuery := fmt.Sprintf("select count(1) from transfer_nft %s", whereCondition)
	fmt.Println("totalQuery:", totalQuery)
	totalCount, err := r.GetTotalNumberOfSpamReport(totalQuery)
	if err != nil {
		return &pb.GetTransferNftReply{
			Code: 500,
			Data: nil,
		}, err
	}
	res, err := r.data.data_query(query)
	if err != nil {
		return &pb.GetTransferNftReply{
			Code: 500,
			Data: nil,
		}, err
	}

	var transferNftList []*pb.TransferNft
	type transfer struct {
		nft_id           string
		chain            string
		contract_address string
		token_id         string
		collection_id    string
		event_type       string
		address_from     *string
		address_to       *string
		block_timestamp  time.Time
		owner            string
	}
	var tf transfer
	for res.Next() {
		// row, ok := res.NextRow()
		// if !ok {
		// 	break
		// }

		var transferNft pb.TransferNft
		if err := res.Scan(&tf.nft_id, &tf.chain, &tf.contract_address, &tf.token_id, &tf.collection_id, &tf.event_type, &tf.address_from, &tf.address_to, &tf.block_timestamp, &tf.owner); err != nil {
			log.Error("failed to scan row err = %v", err)
			return &pb.GetTransferNftReply{
				Code: 500,
				Data: nil,
			}, err
		}

		transferNft.NftId = tf.nft_id
		transferNft.Chain = tf.chain
		transferNft.ContractAddress = tf.contract_address
		transferNft.TokenId = tf.token_id
		transferNft.CollectionId = tf.collection_id
		transferNft.EventType = tf.event_type
		if tf.address_from == nil {
			transferNft.AddressFrom = ""
		} else {
			transferNft.AddressFrom = *tf.address_from
		}
		if tf.address_to == nil {
			transferNft.AddressTo = ""
		} else {
			transferNft.AddressTo = *tf.address_to
		}

		transferNft.BlockTimestamp = tf.block_timestamp.Format("2006-01-02T15:04:05Z")
		transferNft.Owner = tf.owner
		transferNftList = append(transferNftList, &transferNft)
	}

	var result pb.GetTransferNftReply
	result.Code = 200
	result.Limit = limit
	result.Data = transferNftList
	result.Page = page
	result.Total = totalCount

	return &result, nil
}
