package sbom

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-gateway/pkg/client"
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
	contract *client.Contract
}

type SBOMContract interface {
	StoreSBOM(ctx context.Context, id string, doc string) error
	GetSBOM(ctx context.Context, id string) (string, error)
	GetSBOMsByDID(ctx context.Context, did string) ([]string, error)
}

// NewSBOMContract 创建新的 SBOM 合约实例
func NewSBOMContract(network *client.Network) (*SmartContract, error) {
	contract := network.GetContract("sbom")
	return &SmartContract{contract: contract}, nil
}

// StoreSBOM 存储新的 SBOM 文档
func (s *SmartContract) StoreSBOM(ctx context.Context, id string, doc string) error {
	// 验证 SBOM 格式
	var sbom SBOM
	if err := json.Unmarshal([]byte(doc), &sbom); err != nil {
		return fmt.Errorf("SBOM 格式无效: %v", err)
	}

	// 提交交易
	_, err := s.contract.SubmitTransaction("StoreSBOM", id, doc)
	if err != nil {
		return fmt.Errorf("存储 SBOM 失败: %v", err)
	}

	return nil
}

// GetSBOM 获取 SBOM 文档
func (s *SmartContract) GetSBOM(ctx context.Context, id string) (string, error) {
	// 评估交易
	result, err := s.contract.EvaluateTransaction("GetSBOM", id)
	if err != nil {
		return "", fmt.Errorf("获取 SBOM 失败: %v", err)
	}

	return string(result), nil
}

// GetSBOMsByDID 获取指定 DID 的所有 SBOM
func (s *SmartContract) GetSBOMsByDID(ctx context.Context, did string) ([]string, error) {
	// 评估交易
	result, err := s.contract.EvaluateTransaction("GetSBOMsByDID", did)
	if err != nil {
		return nil, fmt.Errorf("获取 SBOM 列表失败: %v", err)
	}

	var sboms []string
	if err := json.Unmarshal(result, &sboms); err != nil {
		return nil, fmt.Errorf("解析 SBOM 列表失败: %v", err)
	}

	return sboms, nil
}

// 确保 SmartContract 实现了 SBOMContract 接口
var _ SBOMContract = (*SmartContract)(nil)
