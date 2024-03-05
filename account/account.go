package account

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type KLAccount struct {
	PrivateKey    *ecdsa.PrivateKey
	AddressCommon common.Address
	AddressStr    string
}

func New(hexkey string) (*KLAccount, error) {
	privateKey_, err := crypto.HexToECDSA(hexkey[4:])
	if err != nil {
		return nil, err
	}

	publicKey := privateKey_.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	common := crypto.PubkeyToAddress(*publicKeyECDSA)
	str := common.Hex()

	return &KLAccount{
		PrivateKey:    privateKey_,
		AddressCommon: common,
		AddressStr:    str,
	}, nil
}
