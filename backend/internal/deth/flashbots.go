package deth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
)

type FlashbotMethod string

var (
	EthSendBundle  = FlashbotMethod("eth_sendBundle")
	EthCallBundle  = FlashbotMethod("eth_callBundle")
	GetUserStats   = FlashbotMethod("flashbots_getUserStats")
	GetBundleStats = FlashbotMethod("flashbots_getBundleStats")
)


// WriteBundleFlashbots sends a bundle to flashbots relay.
func (b *Bot) WriteBundleFlashbots(txns []*types.Transaction) (*FlashbotSendBundleResponse, error) {
	block, err := b.HeaderByNumber(b.ctx, nil)
	if err != nil {
		return nil, err
	}

	txs := []string{}
	for _, v := range txns {
		txnBin, err := v.MarshalBinary()
		if err != nil {
			return nil, errors.New("invalid tx in bundle")
		}

		txs = append(txs, hexutil.Encode(txnBin))
	}
	fresponses := []*FlashbotSendBundleResponse{}
	for i := 0; i < 10; i++ {
		blocknum := block.Number.Add(block.Number, big.NewInt(1))
		res, err := b.SendFlashbotBundle(FlashbotRelayRequest{
			RPC:    "2.0",
			ID:     1,
			Method: EthSendBundle,
			Params: []EthSendBundleParams{
				{
					FlashbotEthParams: FlashbotEthParams{
						FlashbotRelayParams: FlashbotRelayParams{
							BlockNumber: hexutil.EncodeBig(blocknum),
						},
						Txns: txs,
					},
				},
			},
		}, b.privKey)

		if err != nil {
			return nil, err
		}
		defer res.Body.Close()
		response, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		flashbotResponse := &FlashbotSendBundleResponse{}
		err = json.Unmarshal(response, flashbotResponse)
		if err != nil {
			return nil, err
		}

		fresponses = append(fresponses, flashbotResponse)
	}

	return fresponses[0], nil
}

// CallBundleFlashbots simulates a bundle.
func (b *Bot) CallBundleFlashbots(txns []*types.Transaction) (*FlashbotSendBundleResponse, error) {
	block, err := b.BlockNumber(b.ctx)
	if err != nil {
		return nil, err
	}
	block += 1

	txs := []string{}
	for _, v := range txns {
		txnBin, err := v.MarshalBinary()
		if err != nil {
			return nil, errors.New("invalid tx in bundle")
		}

		txs = append(txs, hexutil.Encode(txnBin))
	}

	res, err := b.SendFlashbotBundle(FlashbotRelayRequest{
		RPC:    "2.0",
		ID:     1,
		Method: EthCallBundle,
		Params: []EthCallBundleParams{
			{
				FlashbotEthParams: FlashbotEthParams{
					FlashbotRelayParams: FlashbotRelayParams{
						BlockNumber: hexutil.EncodeUint64(block),
					},
					Txns: txs,
				},
				StateBlockNumber: "latest",
				TimeStamp:        int(time.Now().Unix()),
			},
		},
	}, b.privKey)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	response, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	flashbotResponse := &FlashbotSendBundleResponse{}
	err = json.Unmarshal(response, flashbotResponse)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 && res.StatusCode != 201 && res.StatusCode != 202 {
		eRes := &FlashbotErrorResponse{}
		json.Unmarshal(response, eRes)
		return nil, fmt.Errorf("failed with error: %s", eRes.Error.Message)
	}

	return flashbotResponse, nil
}

func (b *Bot) CheckBundleFlashbots(bundleHash string) (*FlashbotBundleStats, error) {
	block, err := b.BlockNumber(b.ctx)
	if err != nil {
		return nil, err
	}
	block += 1

	res, err := b.SendFlashbotBundle(FlashbotRelayRequest{
		RPC:    "2.0",
		ID:     1,
		Method: GetBundleStats,
		Params: []GetBundleStatsParams{
			{
				FlashbotRelayParams: FlashbotRelayParams{
					BlockNumber: hexutil.EncodeUint64(block),
				},
				BundleHash: bundleHash,
			},
		},
	}, b.privKey)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	response, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	flashbotResponse := &FlashbotBundleStats{}
	err = json.Unmarshal(response, flashbotResponse)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 && res.StatusCode != 201 && res.StatusCode != 202 {
		eRes := &FlashbotErrorResponse{}
		json.Unmarshal(response, eRes)
		return nil, fmt.Errorf("failed with error: %s", eRes.Error.Message)
	}
	return flashbotResponse, nil
}
