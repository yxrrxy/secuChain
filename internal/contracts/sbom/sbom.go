package contracts

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SPDXSBOM 表示 SPDX 格式的 SBOM
type SPDXSBOM struct {
	SPDXID       string       `json:"spdxID"`
	Name         string       `json:"name"`
	VersionInfo  string       `json:"versionInfo"`
	Supplier     string       `json:"supplier"`
	ExternalRefs ExternalRefs `json:"externalRefs"`
}

type ExternalRefs struct {
	ReferenceCategory string `json:"referenceCategory"`
	ReferenceLocator  string `json:"referenceLocator"`
	ReferenceType     string `json:"referenceType"`
}

// CDXSBOM 表示 CycloneDX 格式的 SBOM
type CDXSBOM struct {
	BomRef  string `json:"bom-ref"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	Version string `json:"version"`
	Purl    string `json:"purl"`
}

// SBOM 是存储在区块链上的通用 SBOM 结构
type SBOM struct {
	ID       string `json:"id"`       // 区块链中SBOM的唯一标识
	DID      string `json:"did"`      // 数字身份标识
	SPDXSBOM string `json:"spdxSBOM"` // SPDX格式的SBOM内容
	CDXSBOM  string `json:"cdxSBOM"`  // CycloneDX格式的SBOM内容
}

// SmartContract 提供了 SBOM 相关的函数
type SmartContract struct {
	contractapi.Contract
}

// StoreSBOM 存储新的 SBOM 文档
func (s *SmartContract) StoreSBOM(ctx contractapi.TransactionContextInterface, id string, doc string) error {
	// 检查是否已存在
	exists, err := ctx.GetStub().GetState(id)
	if err != nil {
		return fmt.Errorf("获取 SBOM 状态失败: %v", err)
	}
	if exists != nil {
		return fmt.Errorf("SBOM %s 已存在", id)
	}

	// 验证 SBOM 格式
	var sbom SBOM
	if err := json.Unmarshal([]byte(doc), &sbom); err != nil {
		return fmt.Errorf("SBOM 格式无效: %v", err)
	}

	// 存储 SBOM
	return ctx.GetStub().PutState(id, []byte(doc))
}

// GetSBOM 获取 SBOM 文档
func (s *SmartContract) GetSBOM(ctx contractapi.TransactionContextInterface, id string) (string, error) {
	docBytes, err := ctx.GetStub().GetState(id)
	if err != nil {
		return "", fmt.Errorf("获取 SBOM 状态失败: %v", err)
	}
	if docBytes == nil {
		return "", fmt.Errorf("SBOM %s 不存在", id)
	}

	return string(docBytes), nil
}

// GetSBOMsByDID 获取指定 DID 的所有 SBOM
func (s *SmartContract) GetSBOMsByDID(ctx contractapi.TransactionContextInterface, did string) ([]string, error) {
	// 创建复合键的迭代器
	iterator, err := ctx.GetStub().GetStateByPartialCompositeKey("did~id", []string{did})
	if err != nil {
		return nil, fmt.Errorf("查询 SBOM 失败: %v", err)
	}
	defer iterator.Close()

	var sboms []string
	for iterator.HasNext() {
		queryResponse, err := iterator.Next()
		if err != nil {
			return nil, fmt.Errorf("迭代 SBOM 失败: %v", err)
		}

		sboms = append(sboms, string(queryResponse.Value))
	}

	return sboms, nil
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
