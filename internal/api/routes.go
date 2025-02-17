package api

import (
	"blockSBOM/internal/api/handlers"
	"blockSBOM/internal/api/middleware"

	"github.com/cloudwego/hertz/pkg/app/server"
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(h *server.Hertz, authHandler *handlers.AuthHandler,
	didHandler *handlers.DIDHandler,
	sbomHandler *handlers.SBOMHandler,
	vulnHandler *handlers.VulnHandler,
	managementHandler *handlers.ManagementHandler) {
	// API 版本分组
	v1 := h.Group("/api/v1")

	// 公开路由
	v1.POST("/auth/register", authHandler.Register)
	v1.POST("/auth/login", authHandler.Login)

	// 需要认证的路由
	auth := v1.Group("/", middleware.Auth())
	{
		// Auth 路由
		auth.POST("/auth/refresh", authHandler.RefreshToken)

		// 软件标识工具路由
		did := auth.Group("/did")
		{
			did.POST("/", didHandler.CreateDID)
			did.GET("/:id", didHandler.GetDID)
			did.PUT("/:id", didHandler.UpdateDID)
			did.GET("/", didHandler.ListDIDs)
		}

		// SBOM 工具路由
		sbom := auth.Group("/sbom")
		{
			sbom.POST("/", sbomHandler.CreateSBOM)
			sbom.GET("/:id", sbomHandler.GetSBOM)
			sbom.GET("/did/:did", sbomHandler.ListSBOMsByDID)
			sbom.GET("/search", sbomHandler.SearchSBOMs)
		}

		// 漏洞扫描工具路由
		vuln := auth.Group("/vuln")
		{
			vuln.POST("/", vulnHandler.ReportVulnerability)
			vuln.GET("/:id", vulnHandler.GetVulnerability)
			vuln.GET("/", vulnHandler.ListVulnerabilities)
			vuln.GET("/component/:component", vulnHandler.GetVulnerabilitiesByComponent)
			vuln.GET("/search", vulnHandler.SearchVulnerabilities)
		}

		// 管理系统路由
		management := auth.Group("/management")
		{
			management.GET("/overview", managementHandler.GetOverview)
			management.POST("/upload", managementHandler.UploadPackage)
			management.GET("/packages", managementHandler.ListPackages)
			management.GET("/packages/:id", managementHandler.GetPackageDetails)
			management.GET("/scan/:id", managementHandler.ScanPackage)
		}
	}
}
