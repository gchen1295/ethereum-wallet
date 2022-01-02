package eclient

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/url"
	"nft-engine/pkg/utils"
	"regexp"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/plzn0/go-http-1.17.1"
)

var (
	ContractAddressRegex = regexp.MustCompile(`^0x[a-z0-9A-Z]{40}$`)
	EtherscanURL         = &url.URL{
		Scheme: "https",
		Host:   "api.etherscan.io",
		Path:   url.PathEscape("api"),
	}
)

type Contract struct {
	Address *common.Address
	ABI     *abi.ABI
	RawABI  string
}

// PullContractABI pulls the ABI of a verified contract from etherscan.
func (c *Client) PullContractABI(address *common.Address) (*Contract, error) {
	if !ContractAddressRegex.MatchString(address.String()) {
		return nil, errors.New("invalid contract address")
	}
	contractURL := *EtherscanURL
	//?module=contract&action=getabi&address=0x8a90cab2b38dba80c64b7734e58ee1db38b8992e&format=raw

	form := url.Values{}
	form.Add("module", "contract")
	form.Add("action", "getabi")
	form.Add("address", address.String())
	form.Add("format", "raw")
	// form.Add("apikey", "c.key")
	contractURL.RawQuery = form.Encode()
	origin, err := utils.GenerateRandomString(32)
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

	if res.StatusCode != 200 {
		return nil, errors.New("failed to fetch contract abi")
	}

	rawAbi, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if string(rawAbi) == "Contract source code not verified" {
		return nil, errors.New("contract source code not verified")
	}
	// log.Println(string(rawAbi))
	j, err := abi.JSON(bytes.NewBuffer(rawAbi))
	if err != nil {
		return nil, err
	}

	return &Contract{
		Address: address,
		ABI:     &j,
		RawABI:  string(rawAbi),
	}, nil
}

// QueryContract queries a contract for transactions.
func (c *Client) QueryContract(contract *Contract, method *abi.Method, from *common.Address, payable *big.Int, args ...interface{}) ([]interface{}, error) {
	if contract.ABI == nil {
		return nil, errors.New("contract abi not set")
	}

	if method == nil {
		return nil, errors.New("no method provided")
	}
	log.Println(fmt.Sprintf("%+v\n", args...))
	methodID, err := contract.ABI.Pack(method.Name, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to pack method\nerror: %q", err)
	}

	fee, err := c.geth.SuggestGasPrice(c.Context)
	if err != nil {
		fee = big.NewInt(params.GWei * 21000)
	}

	tip, err := c.geth.SuggestGasTipCap(c.Context)
	if err != nil {
		tip = big.NewInt(params.GWei * 2)
	}

	block, err := c.geth.HeaderByNumber(c.Context, nil)
	if err != nil {
		return nil, err
	}
	gas := big.NewInt(2000000)
	estimateGas, err := c.GetGasEstimate()
	if err == nil {
		gas = big.NewInt(estimateGas.Rapid)
	}

	res, err := c.geth.CallContract(c.Context, ethereum.CallMsg{
		From:      *from,
		To:        contract.Address,
		Gas:       gas.Uint64(),
		GasFeeCap: fee,
		GasTipCap: tip,
		Value:     payable,
		Data:      methodID,
	}, block.Number)
	if err != nil {
		return nil, err
	}

	results, err := contract.ABI.Unpack(method.Name, res)
	if err != nil {
		return nil, err
	}

	return results, nil
}

// Read is used to read a value from contract from the latest known block.
func (c *Client) Read(contract *Contract, methodName string, from *common.Address) ([]interface{}, error) {
	if contract.ABI == nil {
		return nil, errors.New("contract abi not set")
	}

	var targetMethod *abi.Method
	for i, v := range contract.ABI.Methods {
		if strings.Contains(i, methodName) {
			targetMethod = &v
			break
		}
	}

	if targetMethod == nil {
		return nil, errors.New("failed to find method")
	}

	methodID, err := contract.ABI.Pack(targetMethod.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to pack method\nerror: %q", err)
	}

	fee, err := c.geth.SuggestGasPrice(c.Context)
	if err != nil {
		fee = big.NewInt(params.GWei * 21000)
	}

	tip, err := c.geth.SuggestGasTipCap(c.Context)
	if err != nil {
		tip = big.NewInt(params.GWei * 2)
	}

	block, err := c.geth.HeaderByNumber(c.Context, nil)
	if err != nil {
		return nil, err
	}
	gas := big.NewInt(2000000)
	estimateGas, err := c.GetGasEstimate()
	if err == nil {
		gas = big.NewInt(estimateGas.Rapid)
	}

	res, err := c.geth.CallContract(c.Context, ethereum.CallMsg{
		From:      *from,
		To:        contract.Address,
		Gas:       gas.Uint64(),
		GasFeeCap: fee,
		GasTipCap: tip,
		Value:     big.NewInt(0),
		Data:      methodID,
	}, block.Number)
	if err != nil {
		return nil, err
	}

	results, err := contract.ABI.Unpack(methodName, res)
	if err != nil {
		return nil, err
	}

	return results, nil
}

// MonitorNewBlocks starts a monitor for new block heads.
func (c *Client) MonitorNewBlocks(address common.Address, contractABI *abi.ABI) {
	eventChan := make(chan *types.Header)
	// query := ethereum.FilterQuery{
	// 	Addresses: []common.Address{address},
	// }

	sub, err := c.geth.SubscribeNewHead(c.Context, eventChan)
	if err != nil {
		log.Fatal(err)
	}

	defer sub.Unsubscribe()

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLog := <-eventChan:
			go c.SendWebhook("https://discord.com/api/webhooks/926977129978150912/5XPRCCpRLyCWT8WMf5hWjXyx3xKp6XfSBUKVzg0tB_oEqziP7TxmOtejJw8djxDYNbD_", vLog)
			fmt.Println("New Block!")
			fmt.Println("Block Nonce: ", vLog.Nonce.Uint64())
			fmt.Println("Block Number: ", vLog.Number.String())
		}
	}
}
