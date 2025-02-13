package data

import (
	"context"
	"database/sql"
	"encoding/json"

	"errors"

	"fmt"

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

	contract_address string
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
	block_timestamp  []uint8
}

type transaction_log struct {
	chain                 string
	transaction_initiator string
	transaction_hash      string
	block_timestamp       []uint8
	event_type            string
	log_index             uint32
	contract_address      string
	token_id              string
	address_from          *string
	address_to            *string
	owner                 string
	sale_details          *string
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
			if nvalue.sale_details != "" {
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

			action.ContractAddress = cvalue.contract_address
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

func (r *NftTransferRepo) PostNftMute(ctx context.Context, req *pb.PostReportAccountMuteRequest) (*pb.PostReportAccountMuteReply, error) {
	collection_id := req.CollectionId
	account_id := req.AccountId

	// const targetLayout = "2006-01-02T15:04:05Z"

	create_at := time.Now().UTC().Format(time.DateTime)
	insert_str := fmt.Sprintf("insert into account_collection_mute values ('%s','%s','%s', NULL)", account_id, collection_id, create_at)
	insert_err := InsertIntoAccountCollectionMuteTable(r, insert_str)
	if insert_err != nil {
		return nil, insert_err
	}
	data := pb.AccountMuteReport{
		AccountId:    account_id,
		CollectionId: collection_id,
		CreatedAt:    &create_at,
	}
	return &pb.PostReportAccountMuteReply{
		Code:    200,
		Message: "",
		Data:    &data,
	}, nil

}

func (r *NftTransferRepo) PostSpamReport(ctx context.Context, req *pb.PostReportSpamRequest) (*pb.PostReportSpamReply, error) {
	//判断状态
	collection_id := req.CollectionId
	next_status := req.Status
	req_source := req.Source
	var source string

	// 查找先前的 report 记录
	query_str := fmt.Sprintf("select status,create_at,source,create_by,update_by from spam_report where collection_id = '%s'", collection_id)
	res, err := r.data.data_query_single(query_str)
	if err != nil {
		fmt.Println("post spam report fail", collection_id, next_status)
		return nil, err
	}

	type report struct {
		status    string
		create_at []uint8
		source    string
		create_by *string
		update_by *string
	}
	var rt report
	if res != nil {
		if err := res.Scan(&rt.status, &rt.create_at, &rt.source, &rt.create_by, &rt.update_by); err != nil {
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

	const targetLayout = "2006-01-02T15:04:05Z"
	if next_status == "reporting" {
		if req_source == nil {
			source = "firefly"
		} else {
			source = *req_source
			sources := []string{"firefly", "mask-network", "web3bio"}
			if !containsString(sources, source) {
				fmt.Println("source:", source)
				return nil, fmt.Errorf("value of source field should be in %s", sources)
			}
		}
		// 检查 collection 是否已经被report
		if res != nil {
			create_at, err := time.Parse(time.DateTime, string(rt.create_at))
			if err != nil {
				fmt.Println("解析时间时出错:", err)
				return nil, fmt.Errorf("解析时间时出错: %w", err)
			}
			if rt.status == next_status || rt.status == "rejected" {
				create_at_str := create_at.Format(targetLayout)
				update_at := time.Now().UTC().Format(targetLayout)
				update_by := req.UpdateBy
				var createByStr string
				if rt.create_by != nil {
					createByStr = *rt.create_by
				}
				var updateByStr string
				if update_by != nil {
					updateByStr = *update_by
				}
				insert_str := fmt.Sprintf("insert into spam_report values ('%s','%s','%s','%s','%s','%s','%s')",
					collection_id, next_status, create_at_str, update_at, rt.source,
					createByStr, updateByStr)
				insert_err := InsertIntoSpamReportTable(r, insert_str)
				if insert_err != nil {
					return nil, insert_err
				}
				data := pb.SpamReport{
					CollectionId: collection_id,
					Status:       next_status,
					CreateAt:     &create_at_str,
					UpdateAt:     &update_at,
					Source:       &source,
					CreateBy:     &createByStr,
					UpdateBy:     &updateByStr,
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
			create_by := req.CreateBy
			update_by := req.UpdateBy
			var createByStr string
			if create_by != nil {
				createByStr = *create_by
			}
			var updateByStr string
			if update_by != nil {
				updateByStr = *update_by
			}
			insert_str := fmt.Sprintf("insert into spam_report values ('%s','%s','%s','%s','%s','%s','%s')",
				collection_id, next_status, create_at, update_at, source, createByStr, updateByStr)
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
				CreateBy:     create_by,
				UpdateBy:     update_by,
			}
			return &pb.PostReportSpamReply{
				Code:    200,
				Message: "",
				Data:    &data,
			}, nil
		}
	} else if next_status == "approved" || next_status == "rejected" {
		if res != nil {
			create_at, err := time.Parse(time.DateTime, string(rt.create_at))
			if err != nil {
				fmt.Println("解析时间时出错:", err)
				return nil, fmt.Errorf("解析时间时出错: %w", err)
			}
			fmt.Println(res)
			if rt.status == "reporting" {
				update_at := time.Now().UTC().Format(targetLayout)
				create_at_str := create_at.Format(targetLayout)

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
				var createBy, updateBy string
				if rt.create_by != nil {
					createBy = *rt.create_by
				} else {
					createBy = ""
				}
				if req.UpdateBy != nil {
					updateBy = *req.UpdateBy
				} else {
					updateBy = ""
				}
				insert_str := fmt.Sprintf("insert into spam_report values ('%s','%s','%s','%s','%s','%s','%s')", collection_id, next_status, create_at_str, update_at, reply_source, createBy, updateBy)
				insert_err := InsertIntoSpamReportTable(r, insert_str)
				if insert_err != nil {
					return &pb.PostReportSpamReply{
						Code:    500,
						Message: insert_err.Error(),
						Data:    nil,
					}, nil
				}
				data := pb.SpamReport{
					CollectionId: collection_id,
					Status:       next_status,
					CreateAt:     &create_at_str,
					UpdateAt:     &update_at,
					Source:       &reply_source,
					CreateBy:     rt.create_by,
					UpdateBy:     req.UpdateBy,
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
			return &pb.PostReportSpamReply{
				Code:    400,
				Message: fmt.Sprintf("no collection of %s be reported before", collection_id),
				Data:    nil,
			}, nil
		}
	} else {
		return &pb.PostReportSpamReply{
			Code:    400,
			Message: "value of status should be in ['reporting','approved','rejected']",
			Data:    nil,
		}, nil
	}
}

func InsertIntoSpamReportTable(r *NftTransferRepo, insert_str string) error {
	fmt.Println("insert str:", insert_str)
	res, err := r.data.data_query(insert_str)
	if err != nil {
		return fmt.Errorf("writing data into db error:%s", err)
	}
	defer res.Close()
	return nil
}

func InsertIntoAccountCollectionMuteTable(r *NftTransferRepo, insert_str string) error {
	fmt.Println("insert str:", insert_str)
	res, err := r.data.data_query(insert_str)
	if err != nil {
		return fmt.Errorf("writing data into starRocks error:%s", err)
	}
	defer res.Close()
	return nil
}

func UpdataCollectionSpamScore(r *NftTransferRepo, collection_id string) error {
	update_str := fmt.Sprintf("insert into spam_collections_with_bucket values ('%s',100,'')", collection_id)
	// update_str := fmt.Sprintf("update collections_new_test set spam_score=100 where collection_id='%s'", collection_id)
	// update_str := fmt.Sprintf("insert into collections (collection_id,spam_score) values ('%s', 100)", "collection_id")
	fmt.Println("update_str:", update_str)
	res, err := r.data.data_query(update_str)
	if err != nil {
		return fmt.Errorf("writing data into db error:%s", err)
	}
	defer res.Close()
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
		"select collection_id,status,create_at,update_at,source,create_by,update_by from spam_report %s %s %s",
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
		fmt.Println("query from db error:", query_str)
		return &pb.GetReportSpamReply{
			Code: 500,
			Data: nil,
		}, err
	}
	defer res.Close()
	var report_list []*pb.SpamReport

	// row := make([]interface{}, 0)
	const targetLayout = "2006-01-02T15:04:05Z"
	type report struct {
		collection_id string
		status        string
		created_at    []uint8
		update_at     []uint8
		source        string
		create_by     *string
		update_by     *string
	}
	var rt report
	for res.Next() {
		if err := res.Scan(&rt.collection_id, &rt.status, &rt.created_at, &rt.update_at, &rt.source, &rt.create_by, &rt.update_by); err != nil {
			log.Error("failed to scan row err = ", err)
			return nil, err
		}
		var spam_report pb.SpamReport
		var create_at string
		var update_at string
		spam_report.CollectionId = rt.collection_id
		spam_report.Status = rt.status
		parsedTime, err := time.Parse(time.DateTime, string(rt.created_at))
		if err != nil {
			log.Error("无法解析创建时间:", err)
			create_at = ""
		} else {
			create_at = parsedTime.Format(targetLayout)
		}
		spam_report.CreateAt = &create_at
		parsedUpdateTime, err := time.Parse(time.DateTime, string(rt.update_at))
		if err != nil {
			log.Error("无法解析更新时间:", err)
			update_at = ""
		} else {
			update_at = parsedUpdateTime.Format(targetLayout)
		}
		spam_report.UpdateAt = &update_at
		reply_source := rt.source
		spam_report.Source = &reply_source
		spam_report.CreateBy = rt.create_by
		spam_report.UpdateBy = rt.update_by
		report_list = append(report_list, &spam_report)
	}

	var result pb.GetReportSpamReply
	result.Code = 200
	result.Limit = limit
	result.Data = report_list
	result.Page = page
	result.Total = total_count

	return &result, nil
}

func (r *NftTransferRepo) GetTotalNumberOfSpamReport(query_str string) (uint64, error) {
	res, err := r.data.data_query_single(query_str)
	if err != nil {
		fmt.Println("query total count for spam reports err:", err)
		return uint64(0), err
	}
	if res == nil {
		// if !ok {
		fmt.Println("query fail from db")
		return uint64(0), errors.New("query fail from db")
	} else {
		var count uint64
		if err := res.Scan(&count); err != nil {
			log.Error("failed to scan row err = ", err)
			return 0, err
		}
		return count, nil
	}
}

func (r *NftTransferRepo) GetHandleNftinfoFromDB(req *pb.GetNftTransferRequest) (map[string]NftTransfertmpSt, uint64, error) {
	// 创建一个切片来存储每个操作的耗时
	timings := make([]struct {
		operation string
		duration  time.Duration
	}, 0)

	// 定义一个辅助函数来记录时间
	timeTrack := func(start time.Time, name string) {
		elapsed := time.Since(start)
		timings = append(timings, struct {
			operation string
			duration  time.Duration
		}{name, elapsed})
	}

	// 主要逻辑开始
	startTime := time.Now()
	defer timeTrack(startTime, "Total execution")

	//nftlist := make([]*pb.PnftTransferSt, 5, 5)
	if req.Address == "" {

		return nil, 0, errors.New("input address is empty")

	}

	owners := strings.Split(req.Address, ",")

	if len(owners) == 1 {
		owner := owners[0]
		processSecondResultStart := time.Now()
		isTagged, err := r.IsAddressTagged(context.Background(), owner)
		timeTrack(processSecondResultStart, "IsAddressTagged")
		if err != nil {
			// 处理错误
			return nil, 0, err
		}

		if !isTagged {
			return r.getHandleNftinfoForSingleOwner(owner, req)
		}
	}

	str_where := "where owner in ('"
	for i, owner := range owners {
		str_where += owner
		if i == len(owners)-1 {
			break
		}
		str_where += "','"
	}

	str_where += "')"

	if req.AccountId != nil {
		account_id := *req.AccountId
		str_where += fmt.Sprintf(" and collection_id not in (select collection_id from account_collection_mute where account_id='%s' and deleted_at is NULL) ", account_id)
	}

	if req.Network != "" {
		if !strings.Contains(strings.ToLower(req.Network), "all") {
			networks := strings.Split(req.Network, ",")
			networkCondition := combineAndRemoveDuplicates("chain", networks)
			str_where = str_where + " and " + networkCondition
		} else {
			str_where = makeAllNetworksCondition(str_where)
		}
	} else {
		str_where = makeAllNetworksCondition(str_where)
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

		limit_n = 20

		if req.Cursor <= 0 {

			cursor_n = 0

		}

	}

	if limit_n > 100 {

		limit_n = 100

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

	// 修改spam过滤条件，排除白名单中的collection
	spam_filter_condition := " AND (NOT EXISTS (SELECT 1 FROM spam_collections_with_bucket sp WHERE nft_transfer_summary_selected_chains.collection_id = sp.collection_id) OR EXISTS (SELECT 1 FROM nft_whitelist_collections wl WHERE nft_transfer_summary_selected_chains.collection_id = wl.collection_id))"

	first_q := "SELECT " +
		"chain, " +
		"transaction_hash, " +
		"owner, " +
		"event_type, " +
		"block_timestamp " +
		"FROM nft_transfer_summary_selected_chains " +
		str_where +
		spam_filter_condition +
		str_order +
		str_limit
	query_first_start := time.Now()
	first_res, err := r.data.data_query(first_q)
	timeTrack(query_first_start, "Query first result")
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
	var action_num uint64 = 0
	var timestamps []time.Time

	for first_res.Next() {
		if err := first_res.Scan(
			&ts.chain,
			&ts.transaction_hash,
			&ts.owner,
			&ts.event_type,
			&ts.block_timestamp); err != nil {
			log.Error("failed to scan row err = ", err)
			return nil, 0, err
		}
		action_num += 1
		chains = append(chains, ts.chain)
		hashs = append(hashs, ts.transaction_hash)
		_owners = append(_owners, ts.owner)
		event_types = append(event_types, ts.event_type)
		t, err := time.Parse(time.DateTime, string(ts.block_timestamp))
		if err != nil {
			fmt.Println("Error parsing time:", err)
		}
		timestamps = append(timestamps, t)
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

	var sb strings.Builder
	sb.WriteString(" where block_timestamp >= '")
	sb.WriteString(minTime.Format("2006-01-02 15:04:05"))
	sb.WriteString("' and block_timestamp <= '")
	sb.WriteString(maxTime.Format("2006-01-02 15:04:05"))
	sb.WriteString("' and ")

	conditions := []string{owner_condition, hash_condition, chain_condition, event_type_condition}
	sb.WriteString(strings.Join(conditions, " and "))

	sb.WriteString(" and batch_transfer_index=0")

	where_str := sb.String()

	str_sql_p := `select 
		chain, 
		transaction_initiator,
		transaction_hash,
		block_timestamp,
		event_type,
		log_index,
		contract_address,
		token_id,
		address_from,
		address_to,
		owner,
		sale_details 
		from transfer_nft_filter_index_selected_chains`

	str_sql_p += where_str

	// fmt.Println("str_sql:", str_sql_p)

	processSecondResultStart := time.Now()
	log_rows, err := r.data.data_query(str_sql_p)
	if err != nil {
		return nil, 0, err
	}
	defer log_rows.Close()
	timeTrack(processSecondResultStart, "Process second query result")

	data_nodes, err := r.processLogRows(log_rows, uint32(limit_n))

	// 在函数返回之前，打印并排序耗时信息
	sort.Slice(timings, func(i, j int) bool {
		return timings[i].duration > timings[j].duration
	})

	fmt.Println("Operation timings (sorted by duration):")
	for _, timing := range timings {
		fmt.Printf("%s: %v\n", timing.operation, timing.duration)
	}

	return data_nodes, action_num, nil
}

func makeAllNetworksCondition(str_where string) string {
	networks := strings.Split("ethereum,polygon,arbitrum,arbitrum-nova,avalanche,base,bsc,linea,optimism,polygon-zkevm,scroll,zksync-era,zora,gnosis", ",")
	networkCondition := combineAndRemoveDuplicates("chain", networks)
	str_where = str_where + " and " + networkCondition
	return str_where
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
		"select distinct nft_id,chain,contract_address,token_id,collection_id,event_type,address_from,address_to,block_timestamp,owner from transfer_nft %s order by block_timestamp desc limit %d,%d",
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
		block_timestamp  []uint8
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
			log.Error("failed to scan row err = ", err)
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
		block_timestamp, err := time.Parse(time.DateTime, string(tf.block_timestamp))
		if err != nil {
			fmt.Println("解析时间时出错:", err)
			return nil, fmt.Errorf("解析时间时出错: %w", err)
		}
		transferNft.BlockTimestamp = block_timestamp.Format("2006-01-02T15:04:05Z")
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

func (r *NftTransferRepo) getHandleNftinfoForSingleOwner(owner string, req *pb.GetNftTransferRequest) (map[string]NftTransfertmpSt, uint64, error) {
	networkCondition := makeAllNetworksCondition("")

	cursor_n := req.Cursor
	limit_n := req.Limit

	if limit_n <= 0 {
		limit_n = 20
	}
	if limit_n > 100 {
		limit_n = 100
	}
	fullQuery := fmt.Sprintf(`
	WITH grouped_records AS (
		SELECT
			chain,
			transaction_hash,
			event_type,
			block_timestamp
		FROM
			transfer_nft
		WHERE
			owner = LOWER('%s')
			AND batch_transfer_index = 0
			%s
			AND NOT EXISTS (
				SELECT 1
				FROM spam_collections_with_bucket sp
				WHERE transfer_nft.collection_id = sp.collection_id
			)
		GROUP BY
			transaction_hash,
			block_timestamp,
			chain,
			event_type
		ORDER BY
			block_timestamp DESC
		LIMIT %d, %d
	)
	SELECT
		t.chain, 
		t.transaction_initiator,
		t.transaction_hash,
		t.block_timestamp,
		t.event_type,
		t.log_index,
		t.contract_address,
		t.token_id,
		t.address_from,
		t.address_to,
		t.owner,
		t.sale_details 
	FROM
		transfer_nft t
	JOIN grouped_records g ON t.transaction_hash = g.transaction_hash
		AND t.block_timestamp = g.block_timestamp
		AND t.chain = g.chain
		AND t.event_type = g.event_type
	WHERE
		t.owner = LOWER('%s')`, owner, networkCondition, cursor_n, limit_n, owner)

	log_rows, err := r.data.data_query(fullQuery)
	if err != nil {
		return nil, 0, err
	}
	defer log_rows.Close()

	data_nodes, err := r.processLogRows(log_rows, uint32(limit_n))
	return data_nodes, uint64(limit_n), nil
}

func (r *NftTransferRepo) IsAddressTagged(ctx context.Context, address string) (bool, error) {
	query := fmt.Sprintf("SELECT COUNT(1) FROM jdbcEX.full_tagged_address_string WHERE address = LOWER('%s')", address)

	res, err := r.data.data_query_single(query)
	if err != nil {
		return false, fmt.Errorf("查询地址标记时出错: %w", err)
	}

	var count int
	if err := res.Scan(&count); err != nil {
		return false, fmt.Errorf("扫描结果时出错: %w", err)
	}

	return count > 0, nil
}

// 新增的内部方法
func (r *NftTransferRepo) processLogRows(log_rows *sql.Rows, limit_n uint32) (map[string]NftTransfertmpSt, error) {

	data_nodes := make(map[string]NftTransfertmpSt, limit_n)

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
			&ts_log.sale_details,
		); err != nil {
			return nil, fmt.Errorf("扫描行失败: %w", err)
		}

		node := r.createNftTransfertmpSt(ts_log)
		action := r.createDataActionST(ts_log)

		node_ukey := node.network + node.hash + node.owner + node.event_type + ts_log.contract_address
		action_ukey := strconv.FormatUint(uint64(action.index), 10)

		r.updateDataNodes(data_nodes, node, action, node_ukey, action_ukey)
	}

	return data_nodes, nil
}

func (r *NftTransferRepo) createNftTransfertmpSt(ts_log transaction_log) NftTransfertmpSt {
	node := NftTransfertmpSt{
		network:          ts_log.chain,
		init_address:     ts_log.transaction_initiator,
		hash:             ts_log.transaction_hash,
		event_type:       ts_log.event_type,
		contract_address: ts_log.contract_address,
		owner:            ts_log.owner,
		tag:              "collectible",
		actios:           make(map[string]DataActionST),
	}

	if node.network == "bsc" {
		node.network = "binance_smart_chain"
	}

	const targetLayout = "2006-01-02T15:04:05Z"
	t, err := time.Parse(time.DateTime, string(ts_log.block_timestamp))
	if err != nil {
		log.Error("Error parsing time:", err)
	}
	node.timestamp = t.Format(targetLayout)

	if ts_log.sale_details != nil {
		node.sale_details = *ts_log.sale_details
	}

	return node
}

func (r *NftTransferRepo) createDataActionST(ts_log transaction_log) DataActionST {
	action := DataActionST{
		tag:              "collectible",
		event_type:       ts_log.event_type,
		index:            ts_log.log_index,
		token_id:         ts_log.token_id,
		contract_address: ts_log.contract_address,
	}

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

	if action.event_type == "sale" {
		action.event_type = "trade"
	}

	return action
}

func (r *NftTransferRepo) updateDataNodes(data_nodes map[string]NftTransfertmpSt, node NftTransfertmpSt, action DataActionST, node_ukey, action_ukey string) {
	if existing_node, ok := data_nodes[node_ukey]; ok {
		if existing_node.timestamp != node.timestamp {
			actions := existing_node.actios
			for old_action_ukey, exist_action := range actions {
				if exist_action.token_id == action.token_id && existing_node.timestamp < node.timestamp {
					data_nodes[node_ukey].actios[old_action_ukey] = action
				}
			}
		} else {
			if _, ok := existing_node.actios[action_ukey]; ok {
			} else {
				data_nodes[node_ukey].actios[action_ukey] = action
			}
		}
	} else {
		node.actios[action_ukey] = action
		data_nodes[node_ukey] = node
	}
}

// WhitelistCollection represents a whitelisted collection in the database
type WhitelistCollection struct {
	CollectionID string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Description  sql.NullString
}

// AddWhitelistCollection implements the whitelist collection addition
func (r *NftTransferRepo) AddWhitelistCollection(ctx context.Context, req *pb.AddWhitelistCollectionRequest) (*pb.AddWhitelistCollectionReply, error) {
	if r == nil {
		return nil, fmt.Errorf("NftTransferRepo is nil")
	}

	now := time.Now()
	var description sql.NullString
	if req.Description != nil {
		description = sql.NullString{String: *req.Description, Valid: *req.Description != ""}
	} else {
		description = sql.NullString{Valid: false}
	}

	query := fmt.Sprintf(
		"INSERT INTO nft_whitelist_collections (collection_id, created_at, updated_at, description) VALUES ('%s', '%s', '%s', '%s')",
		req.CollectionId, now.Format(time.RFC3339), now.Format(time.RFC3339), description.String,
	)

	// // Create a new context with timeout
	// newCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	// defer cancel()

	_, err := r.data.data_query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to add whitelist collection: %v", err)
	}

	return &pb.AddWhitelistCollectionReply{
		Code:    0,
		Message: "Success",
		Data: &pb.WhitelistCollection{
			CollectionId: req.CollectionId,
			CreatedAt:    now.Format(time.RFC3339),
			UpdatedAt:    now.Format(time.RFC3339),
			Description:  req.Description,
		},
	}, nil
}

// DeleteWhitelistCollection implements the whitelist collection deletion
func (r *NftTransferRepo) DeleteWhitelistCollection(ctx context.Context, req *pb.DeleteWhitelistCollectionRequest) (*pb.DeleteWhitelistCollectionReply, error) {
	query := fmt.Sprintf("DELETE FROM nft_whitelist_collections WHERE collection_id = '%s'", req.CollectionId)

	result, err := r.data.DataBaseCli.ExecContext(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("failed to delete whitelist collection: %v", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get affected rows: %v", err)
	}

	if affected == 0 {
		return &pb.DeleteWhitelistCollectionReply{
			Code:    404,
			Message: "Collection not found in whitelist",
		}, nil
	}

	return &pb.DeleteWhitelistCollectionReply{
		Code:    0,
		Message: "Success",
	}, nil
}

// ListWhitelistCollections implements the whitelist collection listing
func (r *NftTransferRepo) ListWhitelistCollections(ctx context.Context, req *pb.ListWhitelistCollectionsRequest) (*pb.ListWhitelistCollectionsReply, error) {
	page := req.Page
	if page == 0 {
		page = 1
	}
	limit := req.Limit
	if limit == 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	// Get total count
	var total uint64
	countQuery := "SELECT COUNT(*) FROM nft_whitelist_collections"
	countRes, err := r.data.data_query_single(countQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %v", err)
	}
	if err := countRes.Scan(&total); err != nil {
		return nil, fmt.Errorf("failed to scan total count: %v", err)
	}
	fmt.Printf("Total count: %d\n", total)

	// Get paginated results
	query := fmt.Sprintf(
		"SELECT collection_id, created_at, updated_at, description FROM nft_whitelist_collections ORDER BY created_at DESC LIMIT %d OFFSET %d",
		limit, offset,
	)
	fmt.Printf("Query: %s\n", query)

	rows, err := r.data.data_query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list whitelist collections: %v", err)
	}
	defer rows.Close()

	var collections []*pb.WhitelistCollection
	for rows.Next() {
		var collectionID string
		var createdAt []uint8
		var updatedAt []uint8
		var description sql.NullString

		err := rows.Scan(&collectionID, &createdAt, &updatedAt, &description)
		if err != nil {
			fmt.Printf("Error scanning row: %v\n", err)
			return nil, fmt.Errorf("failed to scan whitelist collection: %v", err)
		}

		// Parse timestamps
		createdTime, err := time.Parse(time.DateTime, string(createdAt))
		if err != nil {
			fmt.Printf("Error parsing created_at: %v\n", err)
			continue
		}

		updatedTime, err := time.Parse(time.DateTime, string(updatedAt))
		if err != nil {
			fmt.Printf("Error parsing updated_at: %v\n", err)
			continue
		}

		var desc *string
		if description.Valid {
			desc = &description.String
		}

		collection := &pb.WhitelistCollection{
			CollectionId: collectionID,
			CreatedAt:    createdTime.Format(time.RFC3339),
			UpdatedAt:    updatedTime.Format(time.RFC3339),
			Description:  desc,
		}

		fmt.Printf("Adding collection to result: %+v\n", collection)
		collections = append(collections, collection)
	}

	if err = rows.Err(); err != nil {
		fmt.Printf("Error after scanning rows: %v\n", err)
		return nil, fmt.Errorf("error after scanning rows: %v", err)
	}

	fmt.Printf("Found %d collections\n", len(collections))

	return &pb.ListWhitelistCollectionsReply{
		Code:  0,
		Page:  req.Page,
		Limit: req.Limit,
		Total: total,
		Data:  collections,
	}, nil
}

// IsCollectionWhitelisted checks if a collection is in the whitelist
func (r *NftTransferRepo) IsCollectionWhitelisted(ctx context.Context, collectionID string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM nft_whitelist_collections WHERE collection_id = ?)`

	err := r.data.DataBaseCli.QueryRowContext(context.Background(), query, collectionID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check whitelist status: %v", err)
	}

	return exists, nil
}

// WhitelistAddress represents a whitelisted address in the database
type WhitelistAddress struct {
	Address     string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Description sql.NullString
	Chain       sql.NullString
}

// AddWhitelistAddress implements the whitelist address addition
func (r *NftTransferRepo) AddWhitelistAddress(ctx context.Context, req *pb.AddWhitelistAddressRequest) (*pb.AddWhitelistAddressReply, error) {
	if r == nil {
		return nil, fmt.Errorf("NftTransferRepo is nil")
	}

	now := time.Now().UTC()
	var description string
	if req.Description != nil {
		description = *req.Description
	}

	var chain string
	if req.Chain != nil {
		chain = *req.Chain
	}

	// For OLAP database, we use single quotes for string values and NULL for null values
	query := fmt.Sprintf(
		"INSERT INTO nft_whitelist_addresses (address, created_at, updated_at, description, chain) VALUES ('%s', '%s', '%s', %s, %s)",
		req.Address,
		now.Format(time.DateTime),
		now.Format(time.DateTime),
		formatNullableString(description),
		formatNullableString(chain),
	)

	_, err := r.data.data_query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to add whitelist address: %v", err)
	}

	return &pb.AddWhitelistAddressReply{
		Code:    0,
		Message: "Success",
		Data: &pb.WhitelistAddress{
			Address:     req.Address,
			CreatedAt:   now.Format(time.RFC3339),
			UpdatedAt:   now.Format(time.RFC3339),
			Description: req.Description,
			Chain:       req.Chain,
		},
	}, nil
}

// DeleteWhitelistAddress implements the whitelist address deletion
func (r *NftTransferRepo) DeleteWhitelistAddress(ctx context.Context, req *pb.DeleteWhitelistAddressRequest) (*pb.DeleteWhitelistAddressReply, error) {
	// First check if the address exists
	exists, err := r.IsAddressWhitelisted(ctx, req.Address, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to check address existence: %v", err)
	}

	if !exists {
		return &pb.DeleteWhitelistAddressReply{
			Code:    404,
			Message: "Address not found in whitelist",
		}, nil
	}

	query := fmt.Sprintf("DELETE FROM nft_whitelist_addresses WHERE address = '%s'", req.Address)
	_, err = r.data.data_query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to delete whitelist address: %v", err)
	}

	return &pb.DeleteWhitelistAddressReply{
		Code:    0,
		Message: "Success",
	}, nil
}

// ListWhitelistAddresses implements the whitelist address listing
func (r *NftTransferRepo) ListWhitelistAddresses(ctx context.Context, req *pb.ListWhitelistAddressesRequest) (*pb.ListWhitelistAddressesReply, error) {
	page := req.Page
	if page == 0 {
		page = 1
	}
	limit := req.Limit
	if limit == 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	// Get total count
	var total uint64
	countQuery := "SELECT COUNT(*) FROM nft_whitelist_addresses"
	if req.Chain != nil && *req.Chain != "" {
		countQuery += fmt.Sprintf(" WHERE chain = '%s'", *req.Chain)
	}
	countRes, err := r.data.data_query_single(countQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %v", err)
	}
	if err := countRes.Scan(&total); err != nil {
		return nil, fmt.Errorf("failed to scan total count: %v", err)
	}

	// Get paginated results
	query := "SELECT address, created_at, updated_at, description, chain FROM nft_whitelist_addresses"
	if req.Chain != nil && *req.Chain != "" {
		query += fmt.Sprintf(" WHERE chain = '%s'", *req.Chain)
	}
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT %d OFFSET %d", limit, offset)

	rows, err := r.data.data_query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list whitelist addresses: %v", err)
	}
	defer rows.Close()

	var addresses []*pb.WhitelistAddress
	for rows.Next() {
		var address string
		var createdAt, updatedAt []uint8
		var description, chain sql.NullString

		err := rows.Scan(&address, &createdAt, &updatedAt, &description, &chain)
		if err != nil {
			return nil, fmt.Errorf("failed to scan whitelist address: %v", err)
		}

		// Parse timestamps
		createdTime, err := time.Parse(time.DateTime, string(createdAt))
		if err != nil {
			fmt.Printf("Error parsing created_at: %v\n", err)
			continue
		}

		updatedTime, err := time.Parse(time.DateTime, string(updatedAt))
		if err != nil {
			fmt.Printf("Error parsing updated_at: %v\n", err)
			continue
		}

		var desc *string
		if description.Valid {
			desc = &description.String
		}

		var ch *string
		if chain.Valid {
			ch = &chain.String
		}

		whitelistAddress := &pb.WhitelistAddress{
			Address:     address,
			CreatedAt:   createdTime.Format(time.RFC3339),
			UpdatedAt:   updatedTime.Format(time.RFC3339),
			Description: desc,
			Chain:       ch,
		}

		addresses = append(addresses, whitelistAddress)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after scanning rows: %v", err)
	}

	return &pb.ListWhitelistAddressesReply{
		Code:  0,
		Page:  req.Page,
		Limit: req.Limit,
		Total: total,
		Data:  addresses,
	}, nil
}

// IsAddressWhitelisted checks if an address is in the whitelist
func (r *NftTransferRepo) IsAddressWhitelisted(ctx context.Context, address string, chain *string) (bool, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM nft_whitelist_addresses WHERE address = '%s'", address)
	if chain != nil && *chain != "" {
		query += fmt.Sprintf(" AND chain = '%s'", *chain)
	}

	var count int
	row, err := r.data.data_query_single(query)
	if err != nil {
		return false, fmt.Errorf("failed to execute query: %v", err)
	}

	if err := row.Scan(&count); err != nil {
		return false, fmt.Errorf("failed to check whitelist status: %v", err)
	}

	return count > 0, nil
}

// Helper function to format nullable string for SQL
func formatNullableString(s string) string {
	if s == "" {
		return "NULL"
	}
	return fmt.Sprintf("'%s'", s)
}
