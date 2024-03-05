package client

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/ethereum/ethereumtool/account"
)

type logger interface {
	Println(a ...interface{})
	Printf(format string, a ...interface{})
}

type KLClient struct {
	Ethclient *ethclient.Client
	ChainID   *big.Int
}

func New(rpc string) (*KLClient, error) {
	client_, err := ethclient.Dial(rpc)
	if err != nil {
		return nil, err
	}

	chainID_, err := client_.NetworkID(timeoutCtx())
	if err != nil {
		return nil, err
	}

	return &KLClient{
		Ethclient: client_,
		ChainID:   chainID_,
	}, nil
}

func (client *KLClient) Nonce(account common.Address) (uint64, error) {
	return client.Ethclient.PendingNonceAt(timeoutCtx(), account)
}

func (client *KLClient) SuggestGasPrice() (*big.Int, error) {
	return client.Ethclient.SuggestGasPrice(timeoutCtx())
}

func (client *KLClient) EstimateGas(account common.Address, to string, callData []byte) (uint64, error) {
	contractAddress := common.HexToAddress(to)
	return client.Ethclient.EstimateGas(timeoutCtx(), ethereum.CallMsg{
		From: account,
		To:   &contractAddress,
		Data: callData,
	})
}

// callData, err := nftpool.Pack("dividingTime")
// resultData, err := client.Call(nftPoolAddress, callData)
// result, err := nftpool.Unpack("dividingTime", resultData)
func (client *KLClient) Call(from common.Address, to string, abi *abi.ABI, name string, args ...interface{}) ([]interface{}, error) {
	contractAddress := common.HexToAddress(to)

	callData, err := abi.Pack(name, args...)
	if err != nil {
		return nil, err
	}

	resultData, err := client.Ethclient.CallContract(timeoutCtx(), ethereum.CallMsg{
		From: from,
		To:   &contractAddress,
		Data: callData,
	}, nil)
	if err != nil {
		return nil, err
	}

	return abi.Unpack(name, resultData)
}

func (client *KLClient) SendData(klAccount *account.KLAccount, to string, callData []byte, log logger) (*types.Transaction, error) {
	contractAddress := common.HexToAddress(to)
	log.Println("contract address:", contractAddress.Hex())

	nonce, err := client.Nonce(klAccount.AddressCommon)
	if err != nil {
		log.Println("get nonce error:", err.Error())
		return nil, err
	}

	gasLimit, err := client.EstimateGas(klAccount.AddressCommon, to, callData)
	if err != nil {
		log.Println("get gaslimit error:", err.Error())
		return nil, err
	}
	gasLimit += 50000

	gasPrice, err := client.SuggestGasPrice()
	if err != nil {
		log.Println("get gasPrice error:", err.Error())
		return nil, err
	}
	log.Printf("nonce:%d gasPrice:%v gasLimit:%d\n", nonce, gasPrice, gasLimit)

	tx := types.NewTransaction(nonce, contractAddress, big.NewInt(0), gasLimit, gasPrice, callData)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(client.ChainID), klAccount.PrivateKey)
	if err != nil {
		log.Println("sign tx error:", err.Error())
		return nil, err
	}

	err = client.Ethclient.SendTransaction(timeoutCtx(), signedTx)
	if err != nil {
		log.Println("send transaction error:", err.Error())
		return signedTx, err
	}
	log.Println("tx broadcast:", signedTx.Hash().Hex())

	receipt, err := bind.WaitMined(timeoutCtx2(), client.Ethclient, signedTx)
	if err != nil {
		log.Println("wait mined error:", err)
		if err.Error() == "context deadline exceeded" {
			return signedTx, err
		} else {
			return signedTx, err
		}
	} else {
		log.Printf("receipted - status:%d, blockNumber:%s\n", receipt.Status, receipt.BlockNumber.String())
		if receipt.Status == 1 {
			return signedTx, nil
		} else {
			return signedTx, fmt.Errorf("receipted fail, status: %d", receipt.Status)
		}
	}
}

func (client *KLClient) SendTransaction(klAccount *account.KLAccount, tx *types.Transaction, log logger) error {
	log.Println("contract address:", tx.To().Hex())
	log.Printf("nonce:%d gasPrice:%v gasLimit:%d\n", tx.Nonce(), tx.GasPrice(), tx.Gas())

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(client.ChainID), klAccount.PrivateKey)
	if err != nil {
		log.Println("sign tx error:", err.Error())
		return err
	}

	err = client.Ethclient.SendTransaction(timeoutCtx(), signedTx)
	if err != nil {
		log.Println("send transaction error:", err.Error())
		return err
	}
	log.Println("tx broadcast:", signedTx.Hash().Hex())

	receipt, err := bind.WaitMined(timeoutCtx2(), client.Ethclient, signedTx)
	if err != nil {
		log.Println("wait mined error:", err)
		if err.Error() == "context deadline exceeded" {
			return nil
		} else {
			return err
		}
	} else {
		log.Printf("receipted - status:%d, blockNumber:%s\n", receipt.Status, receipt.BlockNumber.String())
		if receipt.Status == 1 {
			return nil
		} else {
			return fmt.Errorf("receipted fail, status: %d", receipt.Status)
		}
	}
}

func timeoutCtx() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), time.Minute*2)
	return ctx
}

func timeoutCtx2() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), time.Minute*10)
	return ctx
}
