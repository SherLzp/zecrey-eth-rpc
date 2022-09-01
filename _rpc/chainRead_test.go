package _rpc

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/zecrey-labs/zecrey-eth-rpc/_const"
	"log"
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestGetBalance(t *testing.T) {
	toAddress := "0xE9b15a2D396B349ABF60e53ec66Bcf9af262D449"
	cli, err := NewClient(_const.InfuraRinkebyNetwork)
	balance, err := cli.GetBalance(toAddress)
	if err != nil {
		panic(err)
	}
	fmt.Println("balance:", balance)
}

func TestIsContract(t *testing.T) {
	type args struct {
		address string
	}
	cli, _ := NewClient(_const.InfuraRinkebyNetwork)
	tests := []struct {
		name           string
		args           args
		wantIsContract bool
		wantErr        bool
	}{
		{
			name: "valid",
			args: args{
				// smart zecrey address
				address: "0x8b2a865c5856571bc7f9951fee16215a6b2250e1",
			},
			wantIsContract: true,
			wantErr:        false,
		},
		{
			name: "invalid",
			args: args{
				// eth account
				address: "0xfef01c3494e9fbad65a5d12b3852ca87361bc9b2",
			},
			wantIsContract: false,
			wantErr:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIsContract, err := cli.IsContract(tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsContract() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotIsContract != tt.wantIsContract {
				t.Errorf("IsContract() = %v, want %v", gotIsContract, tt.wantIsContract)
			}
		})
	}
}

func TestGetBlockHeaderByNumber(t *testing.T) {
	cli, err := NewClient(_const.LocalNetwork)
	defer cli.Close()
	header, err := cli.GetBlockHeaderByNumber(big.NewInt(1))
	if err != nil {
		t.Error(err)
	}
	fmt.Println("header info:", header)
}

func TestGetBlockInfoByNumber(t *testing.T) {
	cli, err := NewClient(_const.InfuraRinkebyNetwork)
	defer cli.Close()
	blockInfo, err := cli.GetBlockInfoByNumber(big.NewInt(2))
	if err != nil {
		t.Error(err)
	}
	fmt.Println("block info:", blockInfo)
}

func TestGetHeight(t *testing.T) {
	cli, err := NewClient(_const.InfuraRinkebyNetwork)
	defer cli.Close()
	height, err := cli.GetHeight()
	if err != nil {
		t.Error(err)
	}
	fmt.Println("height:", height)
}

func TestGetTransactionByHash(t *testing.T) {
	cli, err := NewClient("http://tf-dex-preview-validator-nlb-6fd109ac8b9d390a.elb.ap-northeast-1.amazonaws.com:8545")
	defer cli.Close()
	hash := "0x3322dfec34ceed12d6d2ca2bbc2004e450eeb31d4eabb3660b324dd52ac382aa"
	tx, isPending, err := cli.GetTransactionByHash(hash)
	if err != nil {
		t.Error(err)
	}
	receipt, err := cli.GetTransactionReceipt(hash)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(receipt.Logs)
	receiptBytes, err := receipt.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	log.Println(string(receiptBytes))
	txBytes, err := tx.MarshalJSON()
	if err != nil {
		t.Error(err)
	}
	fmt.Println("tx:", string(txBytes))
	fmt.Println("isPending:", isPending)
}

func TestGetTransactionReceipt(t *testing.T) {
	cli, err := NewClient(_const.InfuraRinkebyNetwork)
	defer cli.Close()
	successHash := "0xf900253477a50a1cd808f61058f68eb2e73afcb0161c31e82ecafa034d7c8eec"
	receipt, err := cli.GetTransactionReceipt(successHash)
	if err != nil {
		t.Error(err)
	}
	receiptBytes, err := receipt.MarshalJSON()
	if err != nil {
		t.Error(err)
	}
	fmt.Println("success receipt:", string(receiptBytes))
	errorHash := "0x782d3bd436bea25d6d687daa8527abb2487b4d73eb5f9e1350ab70a11adf13b3"
	receipt, err = cli.GetTransactionReceipt(errorHash)
	if err != nil {
		t.Error(err)
	}
	receiptBytes, err = receipt.MarshalJSON()
	if err != nil {
		t.Error(err)
	}
	fmt.Println("error receipt:", string(receiptBytes))
}

func TestPrivateKeyToAddress(t *testing.T) {
	type args struct {
		privateKey *ecdsa.PrivateKey
	}
	sk, err := crypto.GenerateKey()
	if err != nil {
		t.Error(err)
	}
	pubKey := sk.PublicKey
	address := crypto.PubkeyToAddress(pubKey)
	tests := []struct {
		name        string
		args        args
		wantAddress common.Address
		wantErr     bool
	}{
		{
			name: "private key",
			args: args{
				sk,
			},
			wantAddress: address,
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAddress, err := PrivateKeyToAddress(tt.args.privateKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("PrivateKeyToAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotAddress, tt.wantAddress) {
				t.Errorf("PrivateKeyToAddress() = %v, want %v", gotAddress, tt.wantAddress)
			}
		})
	}
}

func TestProviderClient_GetTransactionByHash(t *testing.T) {
	cli, err := NewClient(_const.InfuraRinkebyNetwork)
	tx, _, err := cli.GetTransactionByHash("0xd5dc99aa9d25f510e6e7639327747b2e2cc82cbddfcdc3cbf771922fbf80640d")
	nativeChainId, err := cli.ChainID(context.Background())
	if err != nil {
	}
	fmt.Println(nativeChainId.String())
	msg, err := tx.AsMessage(types.NewLondonSigner(nativeChainId), nil)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(msg.From().Hex())
}

func TestProviderClient_GetTransactionByHash2(t *testing.T) {
	cli, _ := NewClient("http://tf-dex-preview-validator-nlb-6fd109ac8b9d390a.elb.ap-northeast-1.amazonaws.com:8545")
	oTx, _, _ := cli.GetTransactionByHash("0xb17e5e0f5ba1943df54a0f0acc861c4f42df1e628deccc0150afb36a8fcfb663")
	fmt.Println(oTx.To().Hex())
	txInfo, _ := cli.GetTransactionReceipt("0xb17e5e0f5ba1943df54a0f0acc861c4f42df1e628deccc0150afb36a8fcfb663")
	infoBytes, _ := json.Marshal(txInfo)
	fmt.Println(string(infoBytes))
	fmt.Println(string(common.FromHex("0x08c379a00000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000d494e56414c49445f494e50555400000000000000000000000000000000000000")))
}
