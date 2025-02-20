package manage

import (
	"context"
)

// DIDContract 定义了与 DID 相关的接口方法
type DIDContract interface {
	// GetOverview 获取系统概览信息
	GetOverview(ctx context.Context) (overview map[string]interface{}, err error)

	// ListPackages 获取软件包列表
	ListPackages(ctx context.Context) (packages []string, err error)

	// GetPackageDetails 获取指定软件包的详细信息
	GetPackageDetails(ctx context.Context, id string) (details map[string]interface{}, err error)

	// ScanPackage 扫描指定软件包并返回扫描结果
	ScanPackage(ctx context.Context, id string) (result map[string]interface{}, err error)
}
