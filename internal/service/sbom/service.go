package sbom

import (
	"blockSBOM/internal/blockchain/contracts/sbom"
	"blockSBOM/internal/dal/query"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

var ConfigPath = "config.json"

// SBOMService 提供生成SBOM、加载漏洞库和扫描漏洞等功能
type SBOMService struct {
	contract sbom.SBOMContract // 智能合约接口
	repo     *query.SBOMRepository
}

// NewDIDService 创建一个 SBOMService 实例，并注入依赖
func NewSBOMService(contract sbom.SBOMContract, repo *query.SBOMRepository) *SBOMService {
	return &SBOMService{
		contract: contract,
		repo:     repo,
	}
}

// Args 表示生成SBOM或扫描漏洞的参数
type Args struct {
	Language    string
	Format      string
	ProjectPath string
	PackagePath string
	ConfigPath  string
	Token       string
}

type Reply struct {
	Message string `json:"message"`
	Address string `json:"address"`
}

// Vulnerability 表示漏洞信息
type Vulnerability struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

// CreateSBOMRequest 定义创建 SBOM 的请求参数
type CreateSBOMRequest struct {
	Language    string `json:"language" binding:"required,oneof=python java golang"` // 支持的语言
	Format      string `json:"format" binding:"required,oneof=spdx cdx"`             // 支持的格式
	ProjectPath string `json:"project_path" binding:"required"`                      // 项目路径
	ConfigPath  string `json:"config_path" binding:"required"`                       // 配置文件路径
	Token       string `json:"token" binding:"required"`                             // 授权令牌
}

// ScanVulnerabilitiesRequest 定义扫描漏洞的请求参数
type ScanVulnerabilitiesRequest struct {
	Language    string `json:"language" binding:"required,oneof=python java golang"` // 支持的语言
	Format      string `json:"format" binding:"required,oneof=spdx cdx"`             // 支持的格式
	PackagePath string `json:"package_path" binding:"required"`                      // 软件包路径
	ConfigPath  string `json:"config_path" binding:"required"`                       // 配置文件路径
	Token       string `json:"token" binding:"required"`                             // 授权令牌
}

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
	ID        string    `json:"id"`         // 区块链中SBOM的唯一标识
	DID       string    `json:"did"`        // 数字身份标识
	SPDXSBOM  string    `json:"spdxSBOM"`   // SPDX格式的SBOM内容
	CDXSBOM   string    `json:"cdxSBOM"`    // CycloneDX格式的SBOM内容
	CreatedAt time.Time `json:"created_at"` // 创建时间
}

// GenerateSBOM 生成软件SBOM，支持SPDX和CDX格式，支持Python Java和Golang语言
func (s *SBOMService) GenerateSBOM(args *Args, reply *Reply) error {
	cmd := exec.Command("opensca-cli", "-path", args.ProjectPath, "-config", ConfigPath, "-out", fmt.Sprintf("sbom.%s", args.Format))
	cmd.Stdout = log.Writer()
	cmd.Stderr = log.Writer()
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("生成SBOM失败: %w", err)
	}
	reply.Message = "SBOM生成成功"
	reply.Address = fmt.Sprintf("sbom.%s", args.Format)
	return nil
}

// LoadVulnerabilityDatabase 从指定文件中加载本地软件漏洞库
func (s *SBOMService) LoadVulnerabilityDatabase(_ *struct{}, reply *[]Vulnerability) error {
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

// ScanForVulnerabilities 扫描用户上传的软件包，生成漏洞清单信息
func (s *SBOMService) ScanForVulnerabilities(args *Args, reply *string) error {
	cmd := exec.Command("opensca-cli", "-path", args.PackagePath, "-config", args.ConfigPath, "-out", "vulnerabilities.json", "-token", args.Token)
	cmd.Stdout = log.Writer()
	cmd.Stderr = log.Writer()
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("扫描漏洞失败: %w", err)
	}
	*reply = "漏洞扫描完成"
	return nil
}

// SaveSBOMToBlockchain 将SBOM保存到区块链
func (s *SBOMService) SaveSBOMToBlockchain(sbomData string) (string, error) {
	// 生成唯一ID（示例：使用UUID）
	id := fmt.Sprintf("SBOM-%d", time.Now().UnixNano())

	// 调用智能合约保存SBOM
	err := s.contract.StoreSBOM(context.Background(), id, sbomData)
	if err != nil {
		return "", fmt.Errorf("保存SBOM到区块链失败: %w", err)
	}

	return id, nil
}

// GetSBOMFromBlockchain 根据ID从区块链获取SBOM
func (s *SBOMService) GetSBOMFromBlockchain(id string) (string, error) {
	sbomData, err := s.contract.GetSBOM(context.Background(), id)
	if err != nil {
		return "", fmt.Errorf("从区块链获取SBOM失败: %w", err)
	}
	return sbomData, nil
}

// GetSBOMsByDIDFromBlockchain 根据DID从区块链获取所有SBOM
func (s *SBOMService) GetSBOMsByDIDFromBlockchain(did string) ([]string, error) {
	sboms, err := s.contract.GetSBOMsByDID(context.Background(), did)
	if err != nil {
		return nil, fmt.Errorf("从区块链获取SBOM列表失败: %w", err)
	}
	return sboms, nil
}
