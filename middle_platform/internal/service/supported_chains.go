package service

import (
	"context"

	api "middle_platform/api/nft_transfer/v1"
)

// GetSupportedChains implements the NftTransferServer interface method.
// It returns a fixed list of supported chains.
func (s *NftTransferService) GetSupportedChains(ctx context.Context, req *api.GetSupportedChainsRequest) (*api.GetSupportedChainsReply, error) {
	supportedChains := []string{
		"ethereum",
		"polygon",
		"arbitrum",
		"arbitrum-nova",
		"avalanche",
		"base",
		"bsc",
		"linea",
		"optimism",
		"polygon-zkevm",
		"scroll",
		"zksync-era",
		"zora",
		"gnosis",
	}

	chains := make([]*api.ChainInfo, len(supportedChains))
	for i, chainName := range supportedChains {
		chains[i] = &api.ChainInfo{
			Name: chainName,
		}
	}

	reply := &api.GetSupportedChainsReply{
		Code:    0, // 0 indicates success
		Message: "success",
		Data:    chains,
	}

	return reply, nil
}
