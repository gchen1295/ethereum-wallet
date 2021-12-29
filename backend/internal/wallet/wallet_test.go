package wallet

import (
	"bytes"
	"context"
	"log"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
)

type MockChainStateReader struct {
	// number of accounts with previous txns to mock
	numAccounts int
}

func newMockChainStateReader(accounts int) *MockChainStateReader {
	return &MockChainStateReader{
		numAccounts: accounts,
	}
}

func (c *MockChainStateReader) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	panic("not implemented")
}
func (c *MockChainStateReader) StorageAt(ctx context.Context, account common.Address, key common.Hash, blockNumber *big.Int) ([]byte, error) {
	panic("not implemented")
}
func (c *MockChainStateReader) CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error) {
	panic("not implemented")
}
func (c *MockChainStateReader) NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error) {
	if c.numAccounts == 0 {
		return 0, nil
	}

	c.numAccounts--
	return uint64(c.numAccounts), nil
}

func TestWallet(t *testing.T) {
	// should create a valid wallet
	wallet, err := NewFromSeed(bytes.NewBufferString(`some random deterministic seed`).Bytes())
	if err != nil {
		log.Fatal(err)
	}

	if wallet == nil {
		log.Fatal("failed to create a wallet")
	}

	if len(wallet.accounts) != 1 {
		log.Fatalf("create account failed")
	}

	// should match Metamask derivation scheme
	path, err := accounts.ParseDerivationPath("m/44'/60'/0'/0/0")
	if err != nil {
		log.Fatal("accounts.ParseDerivationPath failed with error: ", err)
	}

	if wallet.path.String() != path.String() {
		log.Fatalf("mismatched derivation paths.\nExpected: %s\nActual: %s", path.String(), wallet.path.String())
	}

	// should create an account
	wallet.NewAccount()

	if len(wallet.accounts) != 2 {
		log.Fatalf("create account failed. expected (%d) actual (%d)", 2, len(wallet.accounts))
	}

	if len(wallet.paths) != 2 {
		log.Fatalf("create account failed. expected (%d) actual (%d)", 2, len(wallet.paths))
	}

	// should load all accounts that have previous txns
	chainstateReader := newMockChainStateReader(4)
	wallet.SelfDerive([]accounts.DerivationPath{accounts.DefaultBaseDerivationPath}, chainstateReader)

	if len(wallet.accounts) != 3 {
		log.Fatalf("error deriving previous accounts. expected (%d) actual (%d)", 3, len(wallet.accounts))
	}

	if len(wallet.paths) != 3 {
		log.Fatalf("error deriving previous accounts. expected (%d) actual (%d)", 3, len(wallet.paths))
	}

	// should be able to access previous accounts
	prevAccount := wallet.GetAccount(0)
	if *prevAccount != wallet.accounts[0] {
		log.Fatalf("error fetching previous accounts. expected (%v) actual (%v)", wallet.accounts[0], prevAccount)
	}

	// should encrypt and store private key

	// load wallet using plain text and attempt to read

	// should load private key if found on first load

}
