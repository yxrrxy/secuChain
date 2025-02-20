package did

import (
	"context"
	"fmt"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

// DIDDocument 表示 DID 文档结构
type DIDDocument struct {
	ID        string            `json:"id"`
	PublicKey map[string]string `json:"publicKey"`
	Service   map[string]string `json:"service"`
	Created   string            `json:"created"`
	Updated   string            `json:"updated"`
}

// DIDContract 定义了与区块链交互的接口
type DIDContract interface {
	// 调整方法签名以匹配 SmartContract 的实现
	CreateDID(ctx context.Context, id string, doc string) error
	ResolveDID(ctx context.Context, id string) (string, error)
}

// SmartContract 提供了 DID 相关的函数
type SmartContract struct {
	contract *client.Contract
}

// NewDIDContract 创建新的 DID 合约实例
func NewDIDContract(network *client.Network) (*SmartContract, error) {
	contract := network.GetContract("did")
	return &SmartContract{contract: contract}, nil
}

// CreateDID 创建新的 DID 文档
func (s *SmartContract) CreateDID(ctx context.Context, id string, doc string) error {
	_, err := s.contract.SubmitTransaction("CreateDID", id, doc)
	if err != nil {
		return fmt.Errorf("创建 DID 失败: %v", err)
	}
	return nil
}

// ResolveDID 解析 DID 文档
func (s *SmartContract) ResolveDID(ctx context.Context, id string) (string, error) {
	result, err := s.contract.EvaluateTransaction("ResolveDID", id)
	if err != nil {
		return "", fmt.Errorf("解析 DID 失败: %v", err)
	}
	return string(result), nil
}

// 确保 SmartContract 实现了 DIDContract 接口
var _ DIDContract = (*SmartContract)(nil)
