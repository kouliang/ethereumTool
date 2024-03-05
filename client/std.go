package client

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/kouliang/ethereumtool/account"
)

var std *KLClient

func InitStd(rpc string) (err error) {
	std, err = New(rpc)
	return
}

func ChainID() *big.Int {
	return std.ChainID
}

func Nonce(account common.Address) (uint64, error) {
	return std.Nonce(account)
}

func SuggestGasPrice() (*big.Int, error) {
	return std.SuggestGasPrice()
}

func EstimateGas(account common.Address, to string, callData []byte) (uint64, error) {
	return std.EstimateGas(account, to, callData)
}

func Call(from common.Address, to string, abi *abi.ABI, name string, args ...interface{}) ([]interface{}, error) {
	return std.Call(from, to, abi, name, args...)
}

func SendData(klAccount *account.KLAccount, to string, callData []byte, log logger) (*types.Transaction, error) {
	return std.SendData(klAccount, to, callData, log)
}

func SendTransaction(klAccount *account.KLAccount, tx *types.Transaction, log logger) error {
	return std.SendTransaction(klAccount, tx, log)
}
