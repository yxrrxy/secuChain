package vuln

import (
	contracts "blockSBOM/internal/contracts/vuln"
	"blockSBOM/internal/dal/model"
	"blockSBOM/internal/dal/query"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type VulnService struct {
	contract contracts.VulnContract
	repo     *query.VulnRepository
}

func NewVulnService(contract contracts.VulnContract, repo *query.VulnRepository) *VulnService {
	return &VulnService{
		contract: contract,
		repo:     repo,
	}
}

func (s *VulnService) ReportVulnerability(ctx context.Context, req *ReportVulnRequest) (*model.Vulnerability, error) {
	vuln := &model.Vulnerability{
		ID:                 uuid.New().String(),
		CVE:                req.CVE,
		Description:        req.Description,
		Severity:           req.Severity,
		AffectedComponents: req.AffectedComponents,
		Published:          time.Now().UTC(),
		Updated:            time.Now().UTC(),
	}

	// 先写入区块链
	vulnStr, err := json.Marshal(vuln)
	if err != nil {
		return nil, fmt.Errorf("序列化漏洞信息失败: %v", err)
	}
	if err := s.contract.ReportVulnerability(ctx.(contractapi.TransactionContextInterface), vuln.ID, string(vulnStr)); err != nil {
		return nil, fmt.Errorf("报告区块链漏洞失败: %v", err)
	}

	// 再写入数据库
	if err := s.repo.CreateVulnerability(ctx, vuln); err != nil {
		return nil, fmt.Errorf("存储数据库漏洞失败: %v", err)
	}

	return vuln, nil
}

func (s *VulnService) GetVulnerability(ctx context.Context, id string) (*model.Vulnerability, error) {
	// 优先从数据库查询
	vulnDoc, err := s.repo.GetVulnerability(ctx, id)
	if err == nil {
		return vulnDoc, nil
	}

	// 数据库查询失败，从区块链获取
	vulnStr, err := s.contract.GetVulnerability(ctx.(contractapi.TransactionContextInterface), id)
	if err != nil {
		return nil, fmt.Errorf("获取漏洞信息失败: %v", err)
	}

	// 反序列化漏洞信息
	var vuln model.Vulnerability
	if err := json.Unmarshal([]byte(vulnStr), &vuln); err != nil {
		return nil, fmt.Errorf("反序列化漏洞信息失败: %v", err)
	}

	// 同步到数据库
	if err := s.repo.CreateVulnerability(ctx, &vuln); err != nil {
		fmt.Printf("同步漏洞信息到数据库失败: %v\n", err)
	}

	return &vuln, nil
}

func (s *VulnService) ListVulnerabilities(ctx context.Context, severity string, offset, limit int) ([]*model.Vulnerability, int64, error) {
	return s.repo.ListVulnerabilities(ctx, severity, offset, limit)
}

func (s *VulnService) GetVulnerabilitiesByComponent(ctx context.Context, component string) ([]*model.Vulnerability, error) {
	return s.repo.GetVulnerabilitiesByComponent(ctx, component)
}

func (s *VulnService) SearchVulnerabilities(ctx context.Context, keyword string, offset, limit int) ([]*model.Vulnerability, int64, error) {
	return s.repo.SearchVulnerabilities(ctx, keyword, offset, limit)
}

// Request/Response types
type ReportVulnRequest struct {
	CVE                string   `json:"cve" binding:"required"`
	Description        string   `json:"description" binding:"required"`
	Severity           string   `json:"severity" binding:"required,oneof=low medium high critical"`
	AffectedComponents []string `json:"affectedComponents" binding:"required"`
}

type UpdateVulnRequest struct {
	Description        string   `json:"description" binding:"required"`
	Severity           string   `json:"severity" binding:"required,oneof=low medium high critical"`
	AffectedComponents []string `json:"affectedComponents" binding:"required"`
}
