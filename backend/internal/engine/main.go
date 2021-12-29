package engine

import (
	"nft-engine/pkg/proto"

	"github.com/ethereum/go-ethereum/accounts"
)

var (
	KeyVault *Vault

	StatusChan = make(chan *proto.Notification)
)

func Notify(notif *proto.Notification) {
	go func() {
		StatusChan <- notif
	}()
}

// InitVault setups and loads in the backend keystore for the Vault
func InitVault() []accounts.Account {
	KeyVault = NewVault(APPLICATION_CONFIG_PATH)
	return KeyVault.Accounts()
}
