package main

import (
	"net"
	"net/rpc"
	"testing"
)

// 设置测试环境
func setupServer(t *testing.T) net.Listener {
	sbomService := new(SBOMService)
	rpc.Register(sbomService)
	l, err := net.Listen("tcp", "localhost:12345")
	if err != nil {
		t.Fatalf("Listen error: %v", err)
	}
	go rpc.Accept(l)
	return l
}

// 关闭测试环境
func teardownServer(l net.Listener) {
	l.Close()
}

// 测试生成SBOM功能
func TestGenerateSBOM(t *testing.T) {
	l := setupServer(t)
	defer teardownServer(l)

	client, err := rpc.Dial("tcp", "localhost:12345")
	if err != nil {
		t.Fatalf("Dialing error: %v", err)
	}

	args := Args{Language: "java", Format: "spdx", ProjectPath: "./path/to/project"}
	var reply string

	err = client.Call("SBOMService.GenerateSBOM", &args, &reply)
	if err != nil {
		t.Fatalf("SBOMService.GenerateSBOM error: %v", err)
	}

	expected := "SBOM generated successfully"
	if reply != expected {
		t.Errorf("Expected %q, got %q", expected, reply)
	}
}

// 测试加载漏洞库功能
func TestLoadVulnerabilityDatabase(t *testing.T) {
	l := setupServer(t)
	defer teardownServer(l)

	client, err := rpc.Dial("tcp", "localhost:12345")
	if err != nil {
		t.Fatalf("Dialing error: %v", err)
	}

	var vulnerabilities []Vulnerability

	err = client.Call("SBOMService.LoadVulnerabilityDatabase", new(struct{}), &vulnerabilities)
	if err != nil {
		t.Fatalf("SBOMService.LoadVulnerabilityDatabase error: %v", err)
	}

	if len(vulnerabilities) == 0 {
		t.Error("Expected vulnerabilities to be loaded, but got none")
	}
}

// 测试扫描漏洞功能
func TestScanForVulnerabilities(t *testing.T) {
	l := setupServer(t)
	defer teardownServer(l)

	client, err := rpc.Dial("tcp", "localhost:12345")
	if err != nil {
		t.Fatalf("Dialing error: %v", err)
	}

	args := Args{PackagePath: "./path/to/uploaded/package"}
	var reply string

	err = client.Call("SBOMService.ScanForVulnerabilities", &args, &reply)
	if err != nil {
		t.Fatalf("SBOMService.ScanForVulnerabilities error: %v", err)
	}

	expected := "Vulnerability scan completed successfully"
	if reply != expected {
		t.Errorf("Expected %q, got %q", expected, reply)
	}
}
