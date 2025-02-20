package vuln

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

// Vulnerability 表示漏洞信息结构
type Vulnerability struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Severity    string   `json:"severity"`
	CVSS        float64  `json:"cvss"`
	AffectedPkg []string `json:"affectedPkg"`
	Created     string   `json:"created"`
	References  []string `json:"references"`
	DID         string   `json:"did"`
}

// VulnContract 定义接口
type VulnContract interface {
	ReportVulnerability(ctx context.Context, id string, doc string) error
	GetVulnerability(ctx context.Context, id string) (string, error)
	GetVulnerabilitiesByPackage(ctx context.Context, pkg string) ([]string, error)
}

// SmartContract 提供了漏洞相关的函数
type SmartContract struct {
	contract *client.Contract
}

// NewVulnContract 创建新的漏洞合约实例
func NewVulnContract(network *client.Network) (*SmartContract, error) {
	contract := network.GetContract("vuln")
	return &SmartContract{contract: contract}, nil
}

// ReportVulnerability 报告新的漏洞
func (s *SmartContract) ReportVulnerability(ctx context.Context, id string, doc string) error {
	// 验证漏洞信息格式
	var vuln Vulnerability
	if err := json.Unmarshal([]byte(doc), &vuln); err != nil {
		return fmt.Errorf("漏洞信息格式无效: %v", err)
	}

	_, err := s.contract.SubmitTransaction("ReportVulnerability", id, doc)
	if err != nil {
		return fmt.Errorf("报告漏洞失败: %v", err)
	}
	return nil
}

// GetVulnerability 获取漏洞信息
func (s *SmartContract) GetVulnerability(ctx context.Context, id string) (string, error) {
	result, err := s.contract.EvaluateTransaction("GetVulnerability", id)
	if err != nil {
		return "", fmt.Errorf("获取漏洞信息失败: %v", err)
	}
	return string(result), nil
}

// GetVulnerabilitiesByPackage 获取指定包的所有漏洞
func (s *SmartContract) GetVulnerabilitiesByPackage(ctx context.Context, pkg string) ([]string, error) {
	result, err := s.contract.EvaluateTransaction("GetVulnerabilitiesByPackage", pkg)
	if err != nil {
		return nil, fmt.Errorf("获取包漏洞列表失败: %v", err)
	}

	var vulns []string
	if err := json.Unmarshal(result, &vulns); err != nil {
		return nil, fmt.Errorf("解析漏洞列表失败: %v", err)
	}
	return vulns, nil
}

// 确保 SmartContract 实现了 VulnContract 接口
var _ VulnContract = (*SmartContract)(nil)
