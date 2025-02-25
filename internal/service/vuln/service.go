package vuln

import (
	"blockSBOM/internal/blockchain/contracts/vuln"
	"blockSBOM/internal/dal/model"
	"blockSBOM/internal/dal/query"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/google/uuid"
)

type VulnService struct {
	contract vuln.VulnContract
	repo     *query.VulnRepository
}

func NewVulnService(contract vuln.VulnContract, repo *query.VulnRepository) *VulnService {
	return &VulnService{
		contract: contract,
		repo:     repo,
	}
}

// ReportVulnerability 报告新的漏洞
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

	// 序列化漏洞信息并写入区块链
	vulnStr, err := json.Marshal(vuln)
	if err != nil {
		return nil, fmt.Errorf("序列化漏洞信息失败: %v", err)
	}
	if err := s.contract.ReportVulnerability(ctx, vuln.ID, string(vulnStr)); err != nil {
		return nil, fmt.Errorf("报告区块链漏洞失败: %v", err)
	}

	// 写入数据库
	if err := s.repo.CreateVulnerability(ctx, vuln); err != nil {
		return nil, fmt.Errorf("存储数据库漏洞失败: %v", err)
	}

	return vuln, nil
}

// GetVulnerability 根据 ID 获取漏洞信息
func (s *VulnService) GetVulnerability(ctx context.Context, id string) (*model.Vulnerability, error) {
	// 优先从数据库查询
	vulnDoc, err := s.repo.GetVulnerability(ctx, id)
	if err == nil {
		return vulnDoc, nil
	}

	// 数据库查询失败，从区块链获取
	vulnStr, err := s.contract.GetVulnerability(ctx, id)
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

// ListVulnerabilities 列出漏洞信息
func (s *VulnService) ListVulnerabilities(ctx context.Context, severity string, offset, limit int) ([]*model.Vulnerability, int64, error) {
	return s.repo.ListVulnerabilities(ctx, severity, offset, limit)
}

// GetVulnerabilitiesByComponent 根据组件获取漏洞信息
func (s *VulnService) GetVulnerabilitiesByComponent(ctx context.Context, component string) ([]*model.Vulnerability, error) {
	return s.repo.GetVulnerabilitiesByComponent(ctx, component)
}

// SearchVulnerabilities 搜索漏洞信息
func (s *VulnService) SearchVulnerabilities(ctx context.Context, keyword string, offset, limit int) ([]*model.Vulnerability, int64, error) {
	return s.repo.SearchVulnerabilities(ctx, keyword, offset, limit)
}

// GenerateVulnerabilityDatabase 生成漏洞库
func (s *VulnService) GenerateVulnerabilityDatabase() error {
	// 调用外部脚本生成漏洞库
	cmd := exec.Command("go", "run", "scripts/deploy/vuln/vuln.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("生成漏洞库失败: %v", err)
	}
	return nil
}

// GenerateVulnerabilityCharts 生成漏洞统计图表
func (s *VulnService) GenerateVulnerabilityCharts() error {
	// 调用外部脚本生成漏洞统计图表
	cmd := exec.Command("python3", "scripts/chart/vuln_chart.py")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("生成漏洞统计图表失败: %v", err)
	}
	return nil
}
