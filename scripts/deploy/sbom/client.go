package main

import (
	"fmt"
	"net/rpc"
)

// Args 表示生成SBOM或扫描漏洞的参数
type Args struct {
	Language    string
	Format      string
	ProjectPath string
	PackagePath string
}

// Vulnerability represents the structure of a vulnerability entry
type Vulnerability struct {
	ID                    string   `json:"id"`
	Source                string   `json:"sourceIdentifier"`
	Published             string   `json:"published"`
	LastModified          string   `json:"lastModified"`
	Status                string   `json:"vulnStatus"`
	Description           string   `json:"description"`
	Severity              string   `json:"severity"`
	AttackVector          string   `json:"attackVector"`
	AttackComplexity      string   `json:"attackComplexity"`
	PrivilegesRequired    string   `json:"privilegesRequired"`
	UserInteraction       string   `json:"userInteraction"`
	Scope                 string   `json:"scope"`
	ConfidentialityImpact string   `json:"confidentialityImpact"`
	IntegrityImpact       string   `json:"integrityImpact"`
	AvailabilityImpact    string   `json:"availabilityImpact"`
	CVSSScore             float64  `json:"cvssScore"`
	CWE                   string   `json:"cwe"`
	References            []string `json:"references"`
}

func main() {
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		fmt.Println("Dialing error:", err)
		return
	}

	// 调用 GenerateSBOM 方法
	args := Args{Language: "java", Format: "spdx", ProjectPath: "./path/to/project"}
	var sbomReply string
	err = client.Call("SBOMService.GenerateSBOM", &args, &sbomReply)
	if err != nil {
		fmt.Println("SBOMService.GenerateSBOM error:", err)
		return
	}
	fmt.Println(sbomReply)

	// 调用 LoadVulnerabilityDatabase 方法
	var vulnerabilities []Vulnerability
	err = client.Call("SBOMService.LoadVulnerabilityDatabase", new(struct{}), &vulnerabilities)
	if err != nil {
		fmt.Println("SBOMService.LoadVulnerabilityDatabase error:", err)
		return
	}
	fmt.Printf("Loaded %d vulnerabilities\n", len(vulnerabilities))

	// 调用 ScanForVulnerabilities 方法
	args = Args{PackagePath: "./path/to/uploaded/package"}
	var scanReply string
	err = client.Call("SBOMService.ScanForVulnerabilities", &args, &scanReply)
	if err != nil {
		fmt.Println("SBOMService.ScanForVulnerabilities error:", err)
		return
	}
	fmt.Println(scanReply)
}
