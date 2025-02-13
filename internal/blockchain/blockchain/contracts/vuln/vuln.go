package contracts

import (
	"blockSBOM/backend/internal/blockchain/fabric"
	"blockSBOM/backend/internal/dal/model"
	"encoding/json"
	"fmt"
)

type VulnContract struct {
	client *fabric.FabricClient
}

func NewVulnContract(client *fabric.FabricClient) *VulnContract {
	return &VulnContract{client: client}
}

func (c *VulnContract) ReportVulnerability(vuln *model.Vulnerability) error {
	vulnBytes, err := json.Marshal(vuln)
	if err != nil {
		return fmt.Errorf("序列化漏洞信息失败: %v", err)
	}

	_, err = c.client.Contract.SubmitTransaction("ReportVulnerability", vuln.ID, string(vulnBytes))
	if err != nil {
		return fmt.Errorf("报告漏洞失败: %v", err)
	}

	return nil
}

func (c *VulnContract) GetVulnerability(id string) (*model.Vulnerability, error) {
	result, err := c.client.Contract.EvaluateTransaction("GetVulnerability", id)
	if err != nil {
		return nil, fmt.Errorf("获取漏洞信息失败: %v", err)
	}

	var vuln model.Vulnerability
	if err := json.Unmarshal(result, &vuln); err != nil {
		return nil, fmt.Errorf("反序列化漏洞信息失败: %v", err)
	}

	return &vuln, nil
}
