package data

import (
	"context"
	"errors"
	"fmt"
	"github.com/bytehouse-cloud/driver-go/sdk"
	pb "nft_transfer/api/nft_transfer/v1"
	"time"

	"nft_transfer/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type NftTransferRepo struct {
	data *Data
	log  *log.Helper
}

type NftTransferSt struct {
	Nft_id                string `json:"nft_id"`
	Chain                 string `json:"chain"`
	Collection_id         string `json:"collection_idstring"`
	Transaction_initiator string `json:"block_number"`
	Block_number          int64  `json:"block_number"`
	Block_hash            string `json:"block_hash"`
	Transaction_hash      string `json:"transaction_hash"`
	Block_timestamp       string `json:"block_timestamp"`
	Event_type            string `json:"bvent_type"`
	Log_index             int32  `json:"log_index"`
	Batch_transfer_index  int32  `json:"batch_transfer_index"`
	Contract_address      string `json:"contract_address"`
	Token_id              string `json:"token_id"`
	Address_from          string `json:"address_from"`
	Address_to            string `json:"address_to"`
	Quantity              int32  `json:"quantity"`
	Sale_details          string `json:"sale_details"`
}

// NewNftTransferRepo .
func NewNftTransferRepo(data *Data, logger log.Logger) biz.NftTransferRepo {
	return &NftTransferRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *NftTransferRepo) GetHandleNftinfo(ctx context.Context, req *pb.GetNftTransferRequest) (*pb.GetNftTransferReply, error) {

	var sources *[]string = nil

	handles, err := GetHandleNftinfoFromDB(r.data.DataBaseCli, req, sources)

	fmt.Println(handles)
	fmt.Print("hhhhh:", handles[1].AddressTo)

	return &pb.GetNftTransferReply{
		Code:    200,
		Reason:  "SUCCESS",
		Message: "SUCCESS",
		Data:    handles,
	}, err

}

func GetHandleNftinfoFromDB(db *sdk.Gateway, req *pb.GetNftTransferRequest, sources *[]string) ([]*pb.PnftTransferSt, error) {

	//nftlist := make([]*pb.PnftTransferSt, 5, 5)

	res, err := db.Query("select " +
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
		"sale_details " +
		"from transfer_nft_filter limit 3")

	if err != nil {
		return nil, err
	}

	var nftlist []*pb.PnftTransferSt
	for {
		row, ok := res.NextRow()
		if !ok {
			break
		}

		fmt.Println(row)
		var nts pb.PnftTransferSt
		nts.NftId = row[0].(string)
		nts.Chain = row[1].(string)
		nts.CollectionId = row[2].(string)
		nts.TransactionInitiator = row[3].(string)
		nts.BlockNumber = row[4].(uint64)
		nts.BlockHash = row[5].(string)
		nts.TransactionHash = row[6].(string)
		nts.BlockTimestamp = row[7].(time.Time).String()
		nts.EventType = row[8].(string)
		//nts.log_index = row[9].(int32)
		nts.BatchTransferIndex = row[10].(uint32)

		if row[11] != nil {
			nts.AddressFrom = row[11].(string)
		}
		nts.AddressTo = row[12].(string)
		nts.Quantity = row[13].(uint64)
		if row[14] != nil {
			nts.SaleDetails = row[14].(string)
		}

		/*
			nts.Nft_id = row[0].(string)
			nts.Chain = row[1].(string)
			nts.Collection_id = row[2].(string)
			nts.Transaction_initiator = row[3].(string)
			nts.Block_number = row[4].(int64)
			nts.Block_hash = row[5].(string)
			nts.Transaction_hash = row[6].(string)
			nts.Block_timestamp = row[7].(string)
			nts.Event_type = row[8].(string)
			nts.Log_index = row[9].(int32)
			nts.Batch_transfer_index = row[10].(int32)
			nts.Address_from = row[11].(string)
			nts.Address_to = row[12].(string)
			nts.Quantity = row[13].(int32)
			nts.Sale_details = row[14].(string)
		*/
		fmt.Println(nts)
		nftlist = append(nftlist, &nts)
		fmt.Println("fffffff:", nftlist)
	}

	// Return an error if no data is found
	if len(nftlist) == 0 {
		return nil, errors.New("no data")
	}

	return nftlist, errors.New("success")
}
