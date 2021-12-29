package deth

import (
	"context"
	"crypto/ecdsa"
)

// Worker handles automation of contract minting
type Worker struct {
	key    *ecdsa.PrivateKey
	amount string
	mintFn string

	Cancel context.CancelFunc
}
