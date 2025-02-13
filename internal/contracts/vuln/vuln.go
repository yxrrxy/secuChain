package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
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

// SmartContract 提供了漏洞相关的函数
type SmartContract struct {
	contractapi.Contract
}

// ReportVulnerability 报告新的漏洞
func (s *SmartContract) ReportVulnerability(ctx contractapi.TransactionContextInterface, id string, doc string) error {
	exists, err := ctx.GetStub().GetState(id)
	if err != nil {
		return fmt.Errorf("获取漏洞状态失败: %v", err)
	}
	if exists != nil {
		return fmt.Errorf("漏洞 %s 已存在", id)
	}

	// 验证漏洞信息格式
	var vuln Vulnerability
	if err := json.Unmarshal([]byte(doc), &vuln); err != nil {
		return fmt.Errorf("漏洞信息格式无效: %v", err)
	}

	return ctx.GetStub().PutState(id, []byte(doc))
}

// GetVulnerability 获取漏洞信息
func (s *SmartContract) GetVulnerability(ctx contractapi.TransactionContextInterface, id string) (string, error) {
	docBytes, err := ctx.GetStub().GetState(id)
	if err != nil {
		return "", fmt.Errorf("获取漏洞状态失败: %v", err)
	}
	if docBytes == nil {
		return "", fmt.Errorf("漏洞 %s 不存在", id)
	}

	return string(docBytes), nil
}

// GetVulnerabilitiesByPackage 获取指定包的所有漏洞
func (s *SmartContract) GetVulnerabilitiesByPackage(ctx contractapi.TransactionContextInterface, pkg string) ([]string, error) {
	iterator, err := ctx.GetStub().GetStateByPartialCompositeKey("pkg~id", []string{pkg})
	if err != nil {
		return nil, fmt.Errorf("查询漏洞失败: %v", err)
	}
	defer iterator.Close()

	var vulns []string
	for iterator.HasNext() {
		queryResponse, err := iterator.Next()
		if err != nil {
			return nil, fmt.Errorf("迭代漏洞失败: %v", err)
		}
		vulns = append(vulns, string(queryResponse.Value))
	}

	return vulns, nil
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
