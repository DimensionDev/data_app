package data

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"

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

type report struct {
	status    string
	create_at []uint8
	source    string
	create_by *string
	update_by *string
}

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

func reportSpamToNftScan(collection_id string) error {
	split_list := strings.Split(collection_id, "_")
	contract_address_list := split_list[len(split_list)-1:]
	nftscan_api_key := "YQ3S6KXZ"
	base_url := "https://restapi.nftscan.com/api/v2/submit/spam/contracts"
	client := &http.Client{}
	body := map[string][]string{
		"contract_address_list": contract_address_list,
	}
	json_data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal body: %v", err)
	}
	req, err := http.NewRequest("POST", base_url, bytes.NewBuffer(json_data))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", nftscan_api_key)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		return fmt.Errorf("request failed with status %d: %s", resp.StatusCode, bodyString)
	}

	return nil

}

func (r *NftTransferRepo) reportSpamV2(collection_id string, req_source *string, create_by *string, update_by *string, status string, collectionInfo *string) (*pb.PostReportSpamReply, error) {
	if status != "reporting" && status != "approved" && status != "rejected" {
		return &pb.PostReportSpamReply{
			Code:    400,
			Message: "value of status should be in ['reporting','approved','rejected']",
			Data:    nil,
		}, nil
	}

	final_status := status
	var final_source string
	if req_source != nil && *req_source != "" {
		final_source = *req_source + ":nftscan"
	} else {
		final_source = "nftscan"
	}

	nowUTC := time.Now().UTC()
	update_at := nowUTC.Format(time.RFC3339)
	var createByStr, updateByStr string
	if create_by != nil {
		createByStr = *create_by
	}
	if update_by != nil {
		updateByStr = *update_by
	}

	query_str := fmt.Sprintf("SELECT collection_id, status, create_at, update_at, source, create_by, update_by FROM spam_report WHERE collection_id = '%s' ORDER BY update_at DESC LIMIT 1", collection_id)
	existing_res, query_err := r.data.data_query_single(query_str)

	var insert_str string
	var create_at_str string
	var prev_status string
	var prev_create_by string
	var prev_create_at string
	var hasRecord bool

	if query_err == nil && existing_res != nil {
		type existingReport struct {
			collection_id string
			status        string
			create_at     []uint8
			update_at     []uint8
			source        sql.NullString
			create_by     sql.NullString
			update_by     sql.NullString
		}
		var exRt existingReport
		if err := existing_res.Scan(&exRt.collection_id, &exRt.status, &exRt.create_at, &exRt.update_at, &exRt.source, &exRt.create_by, &exRt.update_by); err == nil {
			hasRecord = true
			prev_status = exRt.status
			prev_create_at = string(exRt.create_at)
			if exRt.create_by.Valid {
				prev_create_by = exRt.create_by.String
			}
		} else {
			hasRecord = false
		}
	}

	if hasRecord {
		if prev_status == "reporting" {
			if final_status == "reporting" {
				// 只更新 update_at 和 update_by
				insert_str = fmt.Sprintf("INSERT INTO spam_report (collection_id, status, create_at, update_at, source, create_by, update_by) VALUES ('%s','%s','%s','%s','%s','%s','%s')",
					collection_id, prev_status, prev_create_at, update_at, final_source, prev_create_by, updateByStr)
			} else if final_status == "approved" {
				// 只有 approved 时才调用 reportSpamToNftScan
				report_err := reportSpamToNftScan(collection_id)
				if report_err != nil {
					log.Errorf("Failed to report spam to NFTScan for collection %s: %v", collection_id, report_err)
					return &pb.PostReportSpamReply{
						Code:    500,
						Message: fmt.Sprintf("Failed to report to NFTScan: %v", report_err),
						Data:    nil,
					}, nil
				}
				create_at_str = prev_create_at
				insert_str = fmt.Sprintf("INSERT INTO spam_report (collection_id, status, create_at, update_at, source, create_by, update_by) VALUES ('%s','%s','%s','%s','%s','%s','%s')",
					collection_id, final_status, create_at_str, update_at, final_source, prev_create_by, updateByStr)
			} else if final_status == "rejected" {
				create_at_str = prev_create_at
				insert_str = fmt.Sprintf("INSERT INTO spam_report (collection_id, status, create_at, update_at, source, create_by, update_by) VALUES ('%s','%s','%s','%s','%s','%s','%s')",
					collection_id, final_status, create_at_str, update_at, final_source, prev_create_by, updateByStr)
			}
		} else if prev_status == "rejected" {
			if final_status == "reporting" {
				// rejected 可以再次 reporting，插入新 reporting 记录
				create_at_str = update_at
				insert_str = fmt.Sprintf("INSERT INTO spam_report (collection_id, status, create_at, update_at, source, create_by, update_by) VALUES ('%s','%s','%s','%s','%s','%s','%s')",
					collection_id, final_status, create_at_str, update_at, final_source, createByStr, updateByStr)
			} else {
				// rejected 状态不能流转到非 reporting
				return &pb.PostReportSpamReply{
					Code:    400,
					Message: "Cannot change status from rejected to non-reporting state.",
					Data:    nil,
				}, nil
			}
		} else if prev_status == "approved" {
			if final_status == "reporting" {
				return &pb.PostReportSpamReply{
					Code:    400,
					Message: "Cannot reporting, already in approved state.",
					Data:    nil,
				}, nil
			}
			// approved -> approved/rejected 不允许
			return &pb.PostReportSpamReply{
				Code:    400,
				Message: "Already in approved state.",
				Data:    nil,
			}, nil
		}
	} else {
		// 无记录，只能 reporting
		if final_status != "reporting" {
			return &pb.PostReportSpamReply{
				Code:    400,
				Message: "First report must be reporting.",
				Data:    nil,
			}, nil
		}
		create_at_str = update_at
		insert_str = fmt.Sprintf("INSERT INTO spam_report (collection_id, status, create_at, update_at, source, create_by, update_by) VALUES ('%s','%s','%s','%s','%s','%s','%s')",
			collection_id, final_status, create_at_str, update_at, final_source, createByStr, updateByStr)
	}

	if insert_str != "" {
		insert_err := InsertIntoSpamReportTable(r, insert_str)
		if insert_err != nil {
			return &pb.PostReportSpamReply{
				Code:    500,
				Message: fmt.Sprintf("DB operation failed: %v", insert_err),
				Data:    nil,
			}, nil
		}
	}

	final_query_str := fmt.Sprintf("SELECT collection_id, status, create_at, update_at, source, create_by, update_by FROM spam_report WHERE collection_id = '%s' ORDER BY update_at DESC LIMIT 1", collection_id)
	final_res, final_err := r.data.data_query_single(final_query_str)
	if final_err != nil {
		log.Errorf("Failed to fetch final spam report for %s: %v", collection_id, final_err)
		return &pb.PostReportSpamReply{
			Code:    200,
			Message: "Reported and DB updated, but failed to fetch final state.",
			Data:    nil,
		}, nil
	}

	var fetched_report pb.SpamReport
	type dbReport struct {
		collection_id string
		status        string
		created_at    []uint8
		updated_at    []uint8
		source        sql.NullString
		create_by     sql.NullString
		update_by     sql.NullString
	}
	var dbRt dbReport
	if err := final_res.Scan(&dbRt.collection_id, &dbRt.status, &dbRt.created_at, &dbRt.updated_at, &dbRt.source, &dbRt.create_by, &dbRt.update_by); err != nil {
		log.Errorf("Failed to scan final report for %s: %v", collection_id, err)
		return &pb.PostReportSpamReply{
			Code:    500,
			Message: "DB scan failed for final result.",
			Data:    nil,
		}, nil
	}

	fetched_report.CollectionId = dbRt.collection_id
	fetched_report.Status = dbRt.status
	createAtStr := formatDbTimestamp(dbRt.created_at)
	fetched_report.CreateAt = &createAtStr
	updateAtStr := formatDbTimestamp(dbRt.updated_at)
	fetched_report.UpdateAt = &updateAtStr
	if dbRt.source.Valid {
		srcStr := dbRt.source.String
		fetched_report.Source = &srcStr
	}
	if dbRt.create_by.Valid {
		cbStr := dbRt.create_by.String
		fetched_report.CreateBy = &cbStr
	}
	if dbRt.update_by.Valid {
		ubStr := dbRt.update_by.String
		fetched_report.UpdateBy = &ubStr
	}

	// 在 status==reporting 时解析 collection_info 并写入 spam_collection_info
	insertSpamCollectionInfoIfPresent(r, collection_id, status, collectionInfo)

	return &pb.PostReportSpamReply{
		Code:    200,
		Message: "Reported and database updated.",
		Data:    &fetched_report,
	}, nil
}

// Helper function to format DB timestamp bytes to RFC3339 string
func formatDbTimestamp(dbTime []uint8) string {
	const targetLayout = "2006-01-02T15:04:05Z"
	t, err := time.Parse(time.DateTime, string(dbTime)) // Assuming DB stores in 'YYYY-MM-DD HH:MM:SS'
	if err != nil {
		log.Warnf("Error parsing timestamp from DB: %v, raw: %s", err, string(dbTime))
		return time.Now().UTC().Format(targetLayout) // Fallback or return empty?
	}
	return t.UTC().Format(targetLayout)
}

func (r *NftTransferRepo) PostSpamReport(ctx context.Context, req *pb.PostReportSpamRequest) (*pb.PostReportSpamReply, error) {
	//判断状态
	collection_id := req.CollectionId
	next_status := req.Status
	req_source := req.Source
	var source string
	if req_source == nil || *req_source == "" {
		source = "firefly"
	} else {
		source = *req_source
		sources := []string{"firefly", "mask-network", "web3bio"}
		if !containsString(sources, source) {
			fmt.Println("source:", source)
			return nil, fmt.Errorf("value of source field should be in %s", sources)
		}
	}

	return r.reportSpamV2(collection_id, &source, req.CreateBy, req.UpdateBy, next_status, req.CollectionInfo)
}

// 在 status==reporting 时解析 collection_info 并写入 spam_collection_info
func insertSpamCollectionInfoIfPresent(r *NftTransferRepo, collectionID string, status string, collectionInfo *string) {
	if status != "reporting" || collectionInfo == nil || *collectionInfo == "" {
		return
	}
	var infoMap map[string]interface{}
	err := json.Unmarshal([]byte(*collectionInfo), &infoMap)
	if err != nil {
		return
	}
	var name, collectionURL, detail *string
	if v, ok := infoMap["name"]; ok {
		if s, ok := v.(string); ok {
			name = &s
		}
	}
	if v, ok := infoMap["collection_url"]; ok {
		if s, ok := v.(string); ok {
			collectionURL = &s
		}
	}
	if v, ok := infoMap["detail"]; ok {
		if s, ok := v.(string); ok {
			detail = &s
		}
	}
	InsertIntoSpamCollectionInfoTable(r, collectionID, name, collectionURL, detail)
}

// Insert collection info into spam_collection_info table
func InsertIntoSpamCollectionInfoTable(r *NftTransferRepo, collectionID string, name, collectionURL, detail *string) error {
	// 先查是否已存在
	checkQuery := fmt.Sprintf("SELECT count(*) AS cnt FROM spam_collection_info WHERE collection_id = '%s'", collectionID)
	res, err := r.data.data_query_single(checkQuery)
	if err != nil {
		return fmt.Errorf("check spam_collection_info existence error: %s", err)
	}
	var cnt int
	if err := res.Scan(&cnt); err != nil {
		return fmt.Errorf("scan count error: %s", err)
	}
	if cnt > 0 {
		fmt.Println("spam_collection_info already exists, skip insert.")
		return nil
	}
	// 不存在则插入
	insertStr := "INSERT INTO spam_collection_info (collection_id, name, collection_url, detail) VALUES ('%s', %s, %s, %s)"
	nameVal := "NULL"
	if name != nil && *name != "" {
		nameVal = fmt.Sprintf("'%s'", *name)
	}
	urlVal := "NULL"
	if collectionURL != nil && *collectionURL != "" {
		urlVal = fmt.Sprintf("'%s'", *collectionURL)
	}
	detailVal := "NULL"
	if detail != nil && *detail != "" {
		detailVal = fmt.Sprintf("'%s'", *detail)
	}
	query := fmt.Sprintf(insertStr, collectionID, nameVal, urlVal, detailVal)
	fmt.Println("insert spam_collection_info:", query)
	insertRes, err := r.data.data_query(query)
	if err != nil {
		return fmt.Errorf("writing data into spam_collection_info error:%s", err)
	}
	defer insertRes.Close()
	return nil
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

	// 1. 收集所有collection_id
	collectionIDSet := make(map[string]struct{})
	var collectionIDs []string
	for res.Next() {
		var tmpID string
		if err := res.Scan(&tmpID, new(string), new([]uint8), new([]uint8), new(sql.NullString), new(sql.NullString), new(sql.NullString)); err == nil {
			if _, exists := collectionIDSet[tmpID]; !exists {
				collectionIDSet[tmpID] = struct{}{}
				collectionIDs = append(collectionIDs, tmpID)
			}
		}
	}
	// 重置游标
	res.Close()
	res, err = r.data.data_query(query_str)
	if err != nil {
		return &pb.GetReportSpamReply{
			Code: 500,
			Data: nil,
		}, err
	}
	defer res.Close()

	// 2. 批量查询spam_collection_info
	infoMap := make(map[string]struct{ name, url, detail *string })
	if len(collectionIDs) > 0 {
		var sb strings.Builder
		sb.WriteString("'")
		sb.WriteString(strings.Join(collectionIDs, "','"))
		sb.WriteString("'")
		infoQuery := fmt.Sprintf("SELECT collection_id, name, collection_url, detail FROM spam_collection_info WHERE collection_id IN (%s)", sb.String())
		rows, err := r.data.data_query(infoQuery)
		if err == nil {
			for rows.Next() {
				var cid string
				var n, u, d sql.NullString
				if err := rows.Scan(&cid, &n, &u, &d); err == nil {
					var namePtr, urlPtr, detailPtr *string
					if n.Valid {
						namePtr = &n.String
					}
					if u.Valid {
						urlPtr = &u.String
					}
					if d.Valid {
						detailPtr = &d.String
					}
					infoMap[cid] = struct{ name, url, detail *string }{namePtr, urlPtr, detailPtr}
				}
			}
			rows.Close()
		}
	}

	// 3. 正式组装结果
	const targetLayout = "2006-01-02T15:04:05Z"
	for res.Next() {
		var rt struct {
			collection_id string
			status        string
			create_at     []uint8
			update_at     []uint8
			source        sql.NullString
			create_by     sql.NullString
			update_by     sql.NullString
		}
		if err := res.Scan(&rt.collection_id, &rt.status, &rt.create_at, &rt.update_at, &rt.source, &rt.create_by, &rt.update_by); err != nil {
			log.Error("failed to scan row err = ", err)
			return nil, err
		}
		var spam_report pb.SpamReport
		var create_at string
		var update_at string
		spam_report.CollectionId = rt.collection_id
		spam_report.Status = rt.status
		parsedTime, err := time.Parse(time.DateTime, string(rt.create_at))
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
		if rt.source.Valid {
			srcStr := rt.source.String
			spam_report.Source = &srcStr
		}
		if rt.create_by.Valid {
			cbStr := rt.create_by.String
			spam_report.CreateBy = &cbStr
		}
		if rt.update_by.Valid {
			ubStr := rt.update_by.String
			spam_report.UpdateBy = &ubStr
		}
		// 批量map取值
		if info, ok := infoMap[rt.collection_id]; ok {
			spam_report.Name = info.name
			spam_report.CollectionUrl = info.url
			spam_report.Detail = info.detail
		}
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

	// Prepare optional fields
	var description, chain, address, createBy sql.NullString

	if req.Description != nil {
		description = sql.NullString{String: *req.Description, Valid: true}
	}
	if req.Chain != nil {
		chain = sql.NullString{String: *req.Chain, Valid: true}
	}
	if req.Address != nil {
		address = sql.NullString{String: *req.Address, Valid: true}
	}
	if req.CreateBy != nil {
		createBy = sql.NullString{String: *req.CreateBy, Valid: true}
	}

	query := fmt.Sprintf(
		"INSERT INTO nft_whitelist_collections (collection_id, created_at, updated_at, description, chain, address, create_by) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s')",
		req.CollectionId,
		now.Format(time.RFC3339),
		now.Format(time.RFC3339),
		description.String,
		chain.String,
		address.String,
		createBy.String,
	)

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
			Chain:        req.Chain,
			Address:      req.Address,
			CreateBy:     req.CreateBy,
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

	// Get paginated results with new fields
	query := fmt.Sprintf(
		"SELECT collection_id, created_at, updated_at, description, chain, address, create_by FROM nft_whitelist_collections ORDER BY created_at DESC LIMIT %d OFFSET %d",
		limit, offset,
	)

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
		var description, chain, address, createBy sql.NullString

		err := rows.Scan(&collectionID, &createdAt, &updatedAt, &description, &chain, &address, &createBy)
		if err != nil {
			return nil, fmt.Errorf("failed to scan whitelist collection: %v", err)
		}

		// Parse timestamps
		createdTime, err := time.Parse(time.DateTime, string(createdAt))
		if err != nil {
			continue
		}

		updatedTime, err := time.Parse(time.DateTime, string(updatedAt))
		if err != nil {
			continue
		}

		// Handle optional fields
		var desc, chainPtr, addressPtr, createByPtr *string
		if description.Valid {
			desc = &description.String
		}
		if chain.Valid {
			chainPtr = &chain.String
		}
		if address.Valid {
			addressPtr = &address.String
		}
		if createBy.Valid {
			createByPtr = &createBy.String
		}

		collection := &pb.WhitelistCollection{
			CollectionId: collectionID,
			CreatedAt:    createdTime.Format(time.RFC3339),
			UpdatedAt:    updatedTime.Format(time.RFC3339),
			Description:  desc,
			Chain:        chainPtr,
			Address:      addressPtr,
			CreateBy:     createByPtr,
		}

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

// 获取 spam_collection_info 的 name, collection_url, detail 字段
func (r *NftTransferRepo) getSpamCollectionInfo(collectionID string) (name, collectionURL, detail *string, err error) {
	query := fmt.Sprintf("SELECT name, collection_url, detail FROM spam_collection_info WHERE collection_id = '%s' LIMIT 1", collectionID)
	row, err := r.data.data_query_single(query)
	if err != nil {
		return nil, nil, nil, err
	}
	var n, u, d sql.NullString
	if err := row.Scan(&n, &u, &d); err != nil {
		return nil, nil, nil, err
	}
	if n.Valid {
		name = &n.String
	}
	if u.Valid {
		collectionURL = &u.String
	}
	if d.Valid {
		detail = &d.String
	}
	return name, collectionURL, detail, nil
}
