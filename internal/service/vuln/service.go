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
	"path/filepath"
)

type VulnService struct {
	contract vuln.VulnContract
	repo     *query.VulnRepository
}
type Vulnerability struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

func NewVulnService(contract vuln.VulnContract, repo *query.VulnRepository) *VulnService {
	return &VulnService{
		contract: contract,
		repo:     repo,
	}
}

// ReportVulnRequest 定义报告新漏洞的请求结构体
type ReportVulnRequest struct {
	CVE                string   `json:"cve" binding:"required"`                                     // CVE 编号
	Description        string   `json:"description" binding:"required"`                             // 漏洞描述
	Severity           string   `json:"severity" binding:"required,oneof=low medium high critical"` // 漏洞严重性
	AffectedComponents []string `json:"affectedComponents" binding:"required"`                      // 受影响的组件
}

// UpdateVulnRequest 定义更新漏洞信息的请求结构体
type UpdateVulnRequest struct {
	ID                 string   `json:"id" binding:"required"`                                      // 漏洞唯一标识
	Description        string   `json:"description" binding:"required"`                             // 更新后的漏洞描述
	Severity           string   `json:"severity" binding:"required,oneof=low medium high critical"` // 更新后的漏洞严重性
	AffectedComponents []string `json:"affectedComponents" binding:"required"`                      // 更新后的受影响组件
}

// LoadVulnerabilityDatabase 从指定文件中加载本地软件漏洞库
func (s *VulnService) LoadVulnerabilityDatabase(_ *struct{}, reply *[]Vulnerability) error {
	// 加载本地软件漏洞库
	InitDatabase()

	relativePath := "../vuln/database.json"
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("获取当前工作目录失败: %w", err)
	}
	filePath := filepath.Join(currentDir, relativePath)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("读取文件 %s 失败: %w", filePath, err)
	}
	var vulnerabilities []Vulnerability
	if err := json.Unmarshal(data, &vulnerabilities); err != nil {
		return fmt.Errorf("从文件 %s 反序列化JSON数据失败: %w", filePath, err)
	}
	*reply = vulnerabilities
	return nil
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

// SearchVulnerabilities 搜索漏洞信息
func (s *VulnService) SearchVulnerabilities(ctx context.Context, keyword string, offset, limit int) ([]*model.Vulnerability, int64, error) {
	return s.repo.SearchVulnerabilities(ctx, keyword, offset, limit)
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
