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

func main() {
	chaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		fmt.Printf("创建链码失败: %v", err)
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("启动链码失败: %v", err)
	}
}
