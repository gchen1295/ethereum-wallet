package eclient

import (
	"context"
	"math/big"
	"nft-engine/internal/request"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Client is used as an interface for the Go Eth-Client
type Client struct {
	request *request.Client
	geth    *ethclient.Client
	context.Context
}

// ClientOptions define options to configure a new Client
type ClientOptions struct {
	// RPC relay link
	RelayLink string
}

// NewClient creates a new ethereum client.
func NewClient(etherscanToken, relayLink string) (*Client, error) {
	client, err := ethclient.Dial(relayLink)
	if err != nil {
		return nil, err
	}

	rqClient, err := request.NewClient(&request.Options{
		DisableDecompression: false,
		UserAgent:            "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.61 Safari/537.36",
	}, "")
	if err != nil {
		return nil, err
	}

	return &Client{
		request: rqClient,
		geth:    client,
		Context: context.Background(),
	}, nil
}

// CheckBalance checks the balance of the given address.
func (c *Client) CheckBalance(address string) (*big.Int, error) {
	account := common.HexToAddress(address)
	balance, err := c.geth.BalanceAt(c.Context, account, nil)

	if err != nil {
		return nil, err
	}

	return balance, nil
}

