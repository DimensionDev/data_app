package data

import (
	"context"
	"errors"
	"fmt"
	"github.com/bytehouse-cloud/driver-go/sdk"
	pb "nft_transfer/api/nft_transfer/v1"
	"nft_transfer/internal/biz"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
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

	handles, err := GetHandleNftinfoFromDB(r.data.DataBaseCli, req)

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
			node.Actions = append(node.Actions, &action)
		}
		data.Nodes = append(data.Nodes, &node)
	}
	data.Total = int32(len(data.Nodes))
	data.Cursor = req.Cursor + req.Limit
	fmt.Println(data)

	return &pb.GetNftTransferReply{
		Code:    200,
		Reason:  "SUCCESS",
		Message: "SUCCESS",
		Data:    &data,
	}, err

}

func GetHandleNftinfoFromDB(db *sdk.Gateway, req *pb.GetNftTransferRequest) (map[string]NftTransfertmpSt, error) {

	//nftlist := make([]*pb.PnftTransferSt, 5, 5)
	if req.Address == "" {
		return nil, errors.New("input address is empty")
	}

	owners := strings.Split(req.Address, ",")

	str_where := "where owner in ('"
	for i, owner := range owners {
		str_where += owner
		if i == len(owners)-1 {
			break
		}
		str_where += "','"
	}
	str_where += "')"

	if req.Network != "" {
		str_where += " and chain='" + req.Network + "'"
	}

	if req.Type != "" {
		str_where += " and event_type='" + req.Type + "'"
	}

	//fmt.Print("order by:", req.OrderBy, req.OrderDirection)
	if req.OrderBy != "" {
		str_where += " order by " + req.OrderBy
		if req.OrderDirection != "" {
			str_where += " " + req.OrderDirection
		}
	}

	if req.Limit > 0 {
		if req.Cursor >= 0 {
			str_where += fmt.Sprintf(" limit  %d,%d", req.Cursor, req.Limit+req.Cursor)
			//str_where += " limit " + req.Limit
		}
	}

	str_sql_p := "select distinct  " +
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
	str_sql_p += str_where

	fmt.Print("sql:", str_sql_p, "\n")

	res, err := db.Query(str_sql_p)

	//fmt.Println("sql eeor:", res, err)

	if err != nil {
		return nil, err
	}

	var data_nodes map[string]NftTransfertmpSt
	data_nodes = make(map[string]NftTransfertmpSt)
	for {
		row, ok := res.NextRow()
		if !ok {
			break
		}
		fmt.Println(row)

		var node NftTransfertmpSt
		if row[0] != nil {
			node.network = row[0].(string)
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

		node_ukey := node.network + node.init_address + node.hash + node.timestamp + node.event_type + node.owner
		var action DataActionST
		if row[8] != nil {
			action.address_from = row[8].(string)
		} else {
			action.address_from = ""
		}
		if row[9] != nil {
			action.address_to = row[9].(string)
		} else {
			action.address_to = ""
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

		action_ukey := node.contract_address + action.token_id

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
		return nil, errors.New("no data in database ")
	}

	return data_nodes, errors.New("success")
}
