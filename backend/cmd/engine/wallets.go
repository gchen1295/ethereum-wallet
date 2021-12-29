package main

import (
	"context"
	"log"
	"nft-engine/internal/engine"
	"nft-engine/pkg/proto"
	"time"
)

// WalletStatusChan is used to relay messages to the frontend wallet handler
var WalletStatusChan = make(chan *proto.GenericResponse)

// VaultHandler implements the VaultHandlerServer
type VaultHandler struct {
	proto.UnimplementedVaultHandlerServer
}

// Init initiates a keystore and returns the addresses of all avaiable wallets
func (VaultHandler) Init(context.Context, *proto.KeystoreOptions) (*proto.KeystoreResponse, error) {
	// Initiate vault
	accounts := engine.InitVault()

	// Retrieve public addresses for frontend
	accPubkeys := []string{}
	for _, account := range accounts {
		accPubkeys = append(accPubkeys, account.Address.Hex())
	}

	return &proto.KeystoreResponse{Accounts: accPubkeys}, nil
}

func (VaultHandler) CreateWallet(ctx context.Context, opts *proto.KeystoreOptions) (*proto.KeystoreResponse, error) {
	accounts, err := engine.KeyVault.CreateWallet([]byte(opts.GetPassphrase()))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println(accounts)
	return &proto.KeystoreResponse{Accounts: accounts}, nil
}

func (VaultHandler) ImportWallet(ctx context.Context, opts *proto.KeystoreOptions) (*proto.KeystoreResponse, error) {
	accounts, err := engine.KeyVault.ImportWallet([]byte(opts.GetAddress()), []byte(opts.GetPassphrase()))
	if err != nil {
		return nil, err
	}

	return &proto.KeystoreResponse{Accounts: accounts}, nil
}

func (VaultHandler) CreateHDWallet(ctx context.Context, opts *proto.KeystoreOptions) (*proto.KeystoreResponse, error) {
	err := engine.KeyVault.ImportHDWallet([]byte(opts.GetMnemonic()), []byte(opts.GetPassphrase()))
	if err != nil {
		return nil, err
	}

	return &proto.KeystoreResponse{Accounts: engine.KeyVault.Addresses()}, nil
}

func (VaultHandler) DeleteWallet(ctx context.Context, opts *proto.KeystoreOptions) (*proto.KeystoreResponse, error) {
	accounts, err := engine.KeyVault.RemoveWallet([]byte(opts.GetAddress()), []byte(opts.GetPassphrase()))
	if err != nil {
		return nil, err
	}

	return &proto.KeystoreResponse{Accounts: accounts}, nil
}

func (VaultHandler) GenerateMnemonic(ctx context.Context, opts *proto.KeystoreOptions) (*proto.MnemonicResponse, error) {
	m, err := engine.GenerateMnemonic()
	if err != nil {
		return nil, err
	}

	return &proto.MnemonicResponse{Mnemonic: m}, nil
}

func (VaultHandler) ListenWallets(e *proto.Empty, stream proto.VaultHandler_ListenWalletsServer) error {
	for range time.Tick(time.Millisecond * 500) {
		stream.SendMsg(&proto.KeystoreResponse{
			Accounts: engine.KeyVault.Addresses(),
		})
	}

	return nil
}

func (VaultHandler) GetWallets(ctx context.Context, opts *proto.KeystoreOptions) (*proto.KeystoreResponse, error) {
	return &proto.KeystoreResponse{Accounts: engine.KeyVault.Addresses()}, nil
}
