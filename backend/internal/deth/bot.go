package deth

import (
	"bufio"
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
)

// Bot manages minting from a contract and wallets
type Bot struct {
	*Client
	Contract *Contract
	PubKey   *ecdsa.PublicKey
	Workers  map[string]*Worker

	// private key for signing flashbot bundles
	// DO NOT USE A REAL WALLET ADDRESS FOR THIS
	privKey *ecdsa.PrivateKey
	ctx     context.Context
	mintFn  *abi.Method
}

type BotOptions struct {
	ClientOptions
	Contract *common.Address
	Key      *ecdsa.PrivateKey
}

func NewBot(opts BotOptions, ctx context.Context) (*Bot, error) {
	log.Println(opts)
	client, err := NewClient(opts.EtherscanToken, opts.RelayLink)
	if err != nil {
		return nil, err
	}

	bot := &Bot{
		Client: client,
		Contract: &Contract{
			Address: opts.Contract,
		},
		ctx: ctx,
	}
	var publicKeyECDSA *ecdsa.PublicKey
	var ok bool
	if opts.Key != nil {
		publicKey := opts.Key.Public()
		publicKeyECDSA, ok = publicKey.(*ecdsa.PublicKey)
		if !ok {
			return nil, errors.New("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
		}

		bot.privKey = opts.Key
		bot.PubKey = publicKeyECDSA
	}

	return bot, nil
}

// SetSigner sets the private key for signing flashbot bundles.
// THIS IS NOT YOUR PRIVATE KEY FOR YOUR WALLET. GENERATE A NEW PK!!!
func (b *Bot) SetSigner(key *ecdsa.PrivateKey) {
	b.privKey = key
}

// Read is used to read a value from contract from the latest known block.
func (b *Bot) Read(methodName string) ([]interface{}, error) {
	if b.Contract.ABI == nil {
		return nil, errors.New("contract abi not set")
	}

	var targetMethod *abi.Method
	for i, v := range b.Contract.ABI.Methods {
		if strings.Contains(i, methodName) {
			targetMethod = &v
			break
		}
	}

	if targetMethod == nil {
		return nil, errors.New("failed to find method")
	}

	methodID, err := b.Contract.ABI.Pack(targetMethod.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to pack method\nerror: %q", err)
	}

	fee, err := b.SuggestGasPrice(b.ctx)
	if err != nil {
		fee = big.NewInt(params.GWei * 21000)
	}

	tip, err := b.SuggestGasTipCap(b.ctx)
	if err != nil {
		tip = big.NewInt(params.GWei * 2)
	}

	block, err := b.HeaderByNumber(b.ctx, nil)
	if err != nil {
		return nil, err
	}
	gas := big.NewInt(2000000)
	estimateGas, err := b.GetGasEstimate()
	if err == nil {
		gas = big.NewInt(estimateGas.Rapid)
	}

	fromAddress := crypto.PubkeyToAddress(*b.PubKey)

	res, err := b.CallContract(b.ctx, ethereum.CallMsg{
		From:      fromAddress,
		To:        b.Contract.Address,
		Gas:       gas.Uint64(),
		GasFeeCap: fee,
		GasTipCap: tip,
		Value:     big.NewInt(0),
		Data:      methodID,
	}, block.Number)
	if err != nil {
		return nil, err
	}

	results, err := b.Contract.ABI.Unpack(methodName, res)
	if err != nil {
		return nil, err
	}

	return results, nil
}

// Write is used to send a txn to be executed on the chain.
func (b *Bot) Write(gas, tip, amount *big.Int, nonce uint64, methodName string, args ...interface{}) (*types.Transaction, error) {

	signedTxn, err := b.NewWriteTxn(gas, tip, amount, nonce, methodName, args...)
	if err != nil {
		return nil, err
	}

	err = b.SendTransaction(b.ctx, signedTxn)
	if err != nil {
		return nil, err
	}

	return signedTxn, nil
}

// NewWriteTxn generates a new txn.
func (b *Bot) NewWriteTxn(fee, tip, amount *big.Int, nonce uint64, methodName string, args ...interface{}) (*types.Transaction, error) {
	var targetMethod *abi.Method
	for i, v := range b.Contract.ABI.Methods {
		if strings.Contains(i, methodName) {
			targetMethod = &v
			break
		}
	}
	if targetMethod == nil {
		return nil, errors.New("failed to find method")
	}

	methodID, err := b.Contract.ABI.Pack(targetMethod.Name, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to pack method\nerror: %q", err)
	}

	if fee.Int64() < 1 {
		fee, err = b.SuggestGasPrice(b.ctx)
		if err != nil {
			fee = big.NewInt(params.GWei * 2)
		}
	}

	if tip.Int64() < 1 {
		tip, err = b.SuggestGasTipCap(b.ctx)
		if err != nil {
			tip = big.NewInt(params.GWei * 2)
		}
	}

	block, err := b.HeaderByNumber(b.ctx, nil)
	if err != nil {
		return nil, err
	}
	log.Println(fee, tip, block.GasUsed, block.GasLimit)
	txn := types.NewTx(&types.DynamicFeeTx{
		To:        b.Contract.Address,
		Gas:       block.GasUsed + 1,
		GasFeeCap: fee,
		GasTipCap: tip,
		Value:     amount,
		Data:      methodID,
		Nonce:     nonce,
	})

	chainID, err := b.ChainID(b.ctx)
	if err != nil {
		return nil, err
	}

	signedTxn, err := types.SignTx(txn, types.LatestSignerForChainID(chainID), b.privKey)
	if err != nil {
		return nil, err
	}

	return signedTxn, nil
}

// GetMethodNames retrieves the method names of the contract.
func (b *Bot) GetMethodNames() []string {
	names := []string{}
	for i, _ := range b.Contract.ABI.Methods {
		names = append(names, i)
	}

	return names
}

func (b *Bot) GetTxn(txnHash common.Hash) (*types.Transaction, bool, error) {
	return b.TransactionByHash(b.ctx, txnHash)
}

// GetMethodByName retrieves a method by name or nil if none found.
func (b *Bot) GetMethodByName(methodName string) *abi.Method {
	var target *abi.Method

	for i, v := range b.Contract.ABI.Methods {
		if i == methodName {
			target = &v
			break
		}
	}

	return target
}

// SetMintFunction sets the mint function for later use.
func (b *Bot) SetMintFunction(fnName string) {
	for i, v := range b.Contract.ABI.Methods {
		if i == fnName {
			b.mintFn = &v
			break
		}
	}
}

// Task runs a mint task on the contract.
func (b *Bot) Task() error {
	consoleReader := bufio.NewReader(os.Stdin)
	var method *abi.Method
	for method == nil {
		log.Println("Input mint function name:")
		input, _ := consoleReader.ReadString('\n')
		input = strings.TrimSpace(input)
		method = b.GetMethodByName(input)
		if method == nil {
			log.Printf("Mint function %s not found!\n", input)
		}
	}

	log.Println(method)

	return nil
}
