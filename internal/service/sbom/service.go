package sbom

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"os/exec"
	"path/filepath"
)

type Vulnerability struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

// SBOMService 提供生成SBOM、加载漏洞库和扫描漏洞等功能
type SBOMService struct{}

// Args 表示生成SBOM或扫描漏洞的参数
type Args struct {
	Language    string
	Format      string
	ProjectPath string
	PackagePath string
	ConfigPath  string
	Token       string
}

// GenerateSBOM 生成软件SBOM，支持SPDX和CDX格式，支持Python Java和Golang语言
func (s *SBOMService) GenerateSBOM(args *Args, reply *string) error {
	cmd := exec.Command("opensca-cli", "-path", args.ProjectPath, "-config", args.ConfigPath, "-out", fmt.Sprintf("sbom.%s", args.Format), "-token", args.Token)
	cmd.Stdout = log.Writer()
	cmd.Stderr = log.Writer()
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("生成SBOM失败: %w", err)
	}
	*reply = "SBOM生成成功"
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

func main() {
	sbomService := new(SBOMService)
	rpc.Register(sbomService)
	l, err := net.Listen("tcp", ":12345")
	if err != nil {
		fmt.Println("监听错误:", err)
		return
	}
	defer l.Close()
	fmt.Println("在端口12345上提供RPC服务")
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("接受连接错误:", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}
