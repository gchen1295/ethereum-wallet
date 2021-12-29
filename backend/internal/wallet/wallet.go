package wallet

import (
	"context"
	"errors"
	"sync"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tyler-smith/go-bip39"
)

const DEFAULT_BASE_DERIVATION_PATH = "m/44'/60'/0'/0/0"

// Wallet implements a manager for an HD wallet as defined by
// https://pkg.go.dev/github.com/ethereum/go-ethereum/accounts#Wallet
// An HD wallet must trace all accounts on first load.
// An account is marked as existing if there is a transaction on it.
type Wallet struct {
	masterKey *hdkeychain.ExtendedKey
	accounts  []accounts.Account
	paths     map[common.Address]accounts.DerivationPath
	path      accounts.DerivationPath // current wallet derivation path

	chainReader ethereum.ChainStateReader

	accountMax  uint32 // number of accounts on current path
	_selfDerive bool   // is self derive running?
	url         accounts.URL
	mux         sync.RWMutex
	ctx         context.Context
	Cleanup     func()
}

// New creates a new HD Wallet from a pass phrase.
// REMEMBER TO SAVE THE SEED PHRASE!!!
func New(password string) (*Wallet, string, error) {
	mnemonic := NewMnemonic()
	seed := bip39.NewSeed(mnemonic, password)
	wallet, err := newWallet(seed)
	if err != nil {
		return nil, "", err
	}

	return wallet, mnemonic, nil
}

// Load loads a wallet from a password and mnemonic.
// Must also provide a chain state reader to reads derivation path and load accounts.
func Load(mnemonic, password string, chain ethereum.ChainStateReader) (*Wallet, error) {
	seed := bip39.NewSeed(mnemonic, password)
	wallet, err := newWallet(seed)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

// SetAccounts sets the current path to
func (wallet *Wallet) SetAccounts(num uint32) {
	wallet.accountMax = num
}

// NewFromMnemonic returns a new wallet from a BIP-39 mnemonic.
// Use NewFromPassword instead.
func NewFromMnemonic(mnemonic string) (*Wallet, error) {
	if mnemonic == "" {
		return nil, errors.New("mnemonic is required")
	}

	if !bip39.IsMnemonicValid(mnemonic) {
		return nil, errors.New("mnemonic is invalid")
	}

	seed := bip39.NewSeed(mnemonic, "")
	wallet, err := newWallet(seed)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

// NewFromSeed returns a new wallet from a BIP-39 seed.
func NewFromSeed(seed []byte) (*Wallet, error) {
	if len(seed) == 0 {
		return nil, errors.New("seed is required")
	}

	return newWallet(seed)
}

func LoadFromString(key string, chain ethereum.ChainStateReader) (*Wallet, error) {
	wallet, err := newWalletFromString(key)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

func (wallet *Wallet) ToString() string {
	return wallet.masterKey.String()
}

// NewMnemonic generates a bip39 compliant mnemonic.
func NewMnemonic() string {
	entropy, _ := bip39.NewEntropy(256)
	mnemonic, _ := bip39.NewMnemonic(entropy)
	return mnemonic
}

// newWallet creates a new HD Wallet.
func newWallet(seed []byte) (*Wallet, error) {
	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.TODO())

	// Get metamask default derivation path
	path, _ := accounts.ParseDerivationPath(DEFAULT_BASE_DERIVATION_PATH)

	wallet := &Wallet{
		masterKey: masterKey,
		accounts:  []accounts.Account{},
		paths:     map[common.Address]accounts.DerivationPath{},
		path:      path,
		ctx:       ctx,
		Cleanup:   cancel,
	}

	_, err = wallet.Derive(path, true)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

// newWallet creates a new HD Wallet.
func newWalletFromString(key string) (*Wallet, error) {
	masterKey, err := hdkeychain.NewKeyFromString(key)

	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.TODO())

	// Get metamask default derivation path
	path, _ := accounts.ParseDerivationPath(DEFAULT_BASE_DERIVATION_PATH)

	wallet := &Wallet{
		masterKey: masterKey,
		accounts:  []accounts.Account{},
		paths:     map[common.Address]accounts.DerivationPath{},
		path:      path,
		ctx:       ctx,
		Cleanup:   cancel,
	}

	_, err = wallet.Derive(path, true)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

// SetChainReader sets the chain state reader for the wallet. Set to nil to turn off self-derive.
func (wallet *Wallet) SetChainReader(reader ethereum.ChainStateReader) {
	wallet.chainReader = reader
}

// traverse iterates the path by 1
func (wallet *Wallet) traverse() accounts.DerivationPath {
	wallet.path[len(wallet.path)-1]++
	return wallet.path
}

// URL retrieves the canonical path under which this wallet is reachable. It is
// user by upper layers to define a sorting order over all wallets from multiple
// backends.
//
// This function is used to satisfy the accounts.Wallet interface and currently
// has no use.
func (wallet *Wallet) URL() accounts.URL {
	return wallet.url
}

// Status returns a textual status to aid the user in the current state of the
// wallet. It also returns an error indicating any failure the wallet might have
// encountered.
//
// This function is used to satisfy the accounts.Wallet interface and currently
// has no use.
func (wallet *Wallet) Status() (string, error) {
	return "", nil
}

// Open initializes access to a wallet instance. It is not meant to unlock or
// decrypt account keys, rather simply to establish a connection to hardware
// wallets and/or to access derivation seeds.
//
// The passphrase parameter may or may not be used by the implementation of a
// particular wallet instance. The reason there is no passwordless open method
// is to strive towards a uniform wallet handling, oblivious to the different
// backend providers.
//
// Please note, if you open a wallet, you must close it to release any allocated
// resources (especially important when working with hardware wallets).
//
// This function is used to satisfy the accounts.Wallet interface and currently
// has no use.
func (wallet *Wallet) Open(passphrase string) error {
	return nil
}

// Close releases any resources held by an open wallet instance.
//
// This function is used to satisfy the accounts.Wallet interface and currently
// has no use.
func (wallet *Wallet) Close() error {
	return nil
}

// Accounts retrieves the list of signing accounts the wallet is currently aware
// of. For hierarchical deterministic wallets, the list will not be exhaustive,
// rather only contain the accounts explicitly pinned during account derivation.
func (wallet *Wallet) Accounts() []accounts.Account {
	wallet.mux.RLock()
	defer wallet.mux.RUnlock()
	return wallet.accounts
}

// Contains returns whether an account is part of this particular wallet or not.
func (wallet *Wallet) Contains(account accounts.Account) bool {
	wallet.mux.RLock()
	defer wallet.mux.RUnlock()

	if _, ok := wallet.paths[account.Address]; ok {
		return true
	}

	return false
}

// NewAccount creates a new account and adds it to the wallet.
func (wallet *Wallet) NewAccount() (*accounts.Account, error) {
	wallet.mux.RLock()
	path, _ := accounts.ParseDerivationPath(wallet.path.String())
	wallet.mux.RUnlock()
	path[len(path)-1]++

	account, err := wallet.Derive(path, true)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

// GetAccount returns an account at the provided index or nil
func (wallet *Wallet) GetAccount(idx int) *accounts.Account {
	if len(wallet.accounts) <= idx {
		return nil
	}

	wallet.mux.RLock()
	defer wallet.mux.RUnlock()

	return &wallet.accounts[idx]
}
