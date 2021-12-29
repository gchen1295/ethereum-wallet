package deth

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"io/ioutil"
	"math/big"
	"net/url"
	"nft-engine/internal/request"
	"nft-engine/pkg/utils"
	"regexp"

	"github.com/plzn0/go-http-1.17.1"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	ContractAddressRegex = regexp.MustCompile(`^0x[a-z0-9A-Z]{40}$`)

	flashbotsRpc = &url.URL{
		Scheme: "https",
		Host:   "relay.flashbots.net",
	}
)

type IOTypes struct {
	Indexed      string `json:"indexed,omitempty"`
	InternalType string `json:"internalType"`
	Name         string `json:"name"`
	Type         string `json:"type"`
}

type ABI []struct {
	Inputs          []IOTypes `json:"inputs"`
	StateMutability string    `json:"stateMutability,omitempty"`
	Type            string    `json:"type"`
	Anonymous       bool      `json:"anonymous,omitempty"`
	Name            string    `json:"name,omitempty"`
	Outputs         []IOTypes `json:"outputs,omitempty"`
}

type Contract struct {
	Address *common.Address
	ABI     *abi.ABI
	RawABI  string
}

// Client is used as an interface for the Go Eth-Client
type Client struct {
	request *request.Client
	*ethclient.Client

	hostURL *url.URL
	key     string
}

// ClientOptions define options to configure a new Client
type ClientOptions struct {
	// etherscan api key
	EtherscanToken string
	// RPC relay link
	RelayLink string
}

func NewClient(etherscanToken, relayLink string) (*Client, error) {
	client, err := ethclient.Dial(relayLink)
	if err != nil {
		return nil, err
	}

	return &Client{
		request: request.NewClient(&request.Options{
			DisableDecompression: false,
			UserAgent:            "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.61 Safari/537.36",
		}, nil),
		hostURL: &url.URL{
			Scheme: "https",
			Host:   "api.etherscan.io",
			Path:   url.PathEscape("api"),
		},
		key:    etherscanToken,
		Client: client,
	}, nil

}

type GasOracleResponse struct {
	Code int        `json:"code"`
	Data OracleData `json:"data" `
}

type OracleData struct {
	Rapid     int64   `json:"rapid"`
	Fast      int64   `json:"fast"`
	Standard  int64   `json:"standard"`
	Slow      int64   `json:"slow"`
	TimeStamp int64   `json:"timestamp"`
	USDPrice  float64 `json:"priceUSD"`
}

// GetGasEstimate fetches a gas estimate from etherchain
func (c *Client) GetGasEstimate() (*OracleData, error) {
	// https://etherchain.org/api/gasnow
	ep, _ := url.Parse("https://etherchain.org/api/gasnow")

	headers := http.Header{}
	headers.Add("Accept", "application/json")

	req, err := http.NewRequest("GET", ep.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header = headers
	res, err := c.request.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	result := &GasOracleResponse{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return nil, err
	}

	if result.Code == 200 {
		return &result.Data, nil
	}

	return nil, errors.New("failed to fetch gas price")
}

// SendFlashbotBundle sends a bundle of txns to the flashbots relay.
func (c *Client) SendFlashbotBundle(query FlashbotRelayRequest, privKey *ecdsa.PrivateKey) (*http.Response, error) {
	body, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	signature, err := utils.SignPayload(body, privKey)
	if err != nil {
		return nil, err
	}

	headers := http.Header{}
	headers.Add("accept", "application/json")
	headers.Add("content-type", "application/json")
	headers.Add("X-Flashbots-Signature", signature)

	req, err := http.NewRequest("POST", flashbotsRpc.String(), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header = headers

	return c.request.Do(req)
}

// PullContractABI pulls the ABI of a verified contract from etherscan.
func (c *Client) PullContractABI(address string) (*abi.ABI, error) {
	if !ContractAddressRegex.MatchString(address) {
		return nil, errors.New("invalid contract address")
	}
	contractURL := *c.hostURL
	//?module=contract&action=getabi&address=0x8a90cab2b38dba80c64b7734e58ee1db38b8992e&format=raw

	form := url.Values{}
	form.Add("module", "contract")
	form.Add("action", "getabi")
	form.Add("address", address)
	form.Add("format", "raw")
	// form.Add("apikey", "c.key")
	contractURL.RawQuery = form.Encode()
	origin, err := GenerateRandomString(32)
	if err != nil {
		origin = ""
	}
	headers := http.Header{}
	headers.Add("accept-encoding", "gzip, deflate, br")
	headers.Add("origin", origin)
	headers.Add("sec-ch-ua", `"Chromium";v="94", "Google Chrome";v="94", ";Not A Brand";v="99"`)
	headers.Add("sec-ch-ua-mobile", "?0")
	headers.Add("sec-fetch-site", "same-origin")
	headers.Add("sec-fetch-mode", "cors")
	headers.Add("sec-fetch-dest", "empty")
	headers.Add("accept-language", "application/json")
	headers.Add("content-type", "application/x-www-form-urlencoded")
	reqUrl, err := url.Parse(contractURL.String())
	if err != nil {
		return nil, err
	}

	req := &http.Request{
		Method: http.MethodGet,
		URL:    reqUrl,
		Header: headers,
	}

	res, err := c.request.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	j, err := abi.JSON(res.Body)
	if err != nil {
		return nil, err
	}

	return &j, nil
}

func (c *Client) CreateInstance() error {

	return nil
}

// CheckBalance checks the balance of the given address.
func (c *Client) CheckBalance(ctx context.Context, address string) (*big.Int, error) {
	account := common.HexToAddress(address)
	balance, err := c.BalanceAt(ctx, account, nil)

	if err != nil {
		return nil, err
	}

	return balance, nil
}
