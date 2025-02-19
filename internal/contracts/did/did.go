package contracts

import (
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
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
	CreateDID(ctx contractapi.TransactionContextInterface, id string, doc string) error
	ResolveDID(ctx contractapi.TransactionContextInterface, id string) (string, error)
}

// SmartContract 提供了 DID 相关的函数
type SmartContract struct {
	contractapi.Contract
}

// CreateDID 创建新的 DID 文档
func (s *SmartContract) CreateDID(ctx contractapi.TransactionContextInterface, id string, doc string) error {
	exists, err := ctx.GetStub().GetState(id)
	if err != nil {
		return fmt.Errorf("获取 DID 状态失败: %v", err)
	}
	if exists != nil {
		return fmt.Errorf("DID %s 已存在", id)
	}

	return ctx.GetStub().PutState(id, []byte(doc))
}

// ResolveDID 解析 DID 文档
func (s *SmartContract) ResolveDID(ctx contractapi.TransactionContextInterface, id string) (string, error) {
	docBytes, err := ctx.GetStub().GetState(id)
	if err != nil {
		return "", fmt.Errorf("获取 DID 状态失败: %v", err)
	}
	if docBytes == nil {
		return "", fmt.Errorf("DID %s 不存在", id)
	}

	return string(docBytes), nil
}

// 确保 SmartContract 实现了 DIDContract 接口
var _ DIDContract = (*SmartContract)(nil)
