package api

import (
	"blockSBOM/internal/api/handlers"
	"blockSBOM/internal/api/middleware"

	"github.com/cloudwego/hertz/pkg/app/server"
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(h *server.Hertz, authHandler *handlers.AuthHandler,
	sbomHandler *handlers.SBOMHandler,
	vulnHandler *handlers.VulnHandler) {
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
		//did := auth.Group("/did")
		//{
		//	did.POST("/registerDID", didHandler.RegisterDID)
		//	did.GET("/:did", didHandler.ResolveDIDAPI)
		//	did.PUT("/:did", didHandler.UpdateDID)
		//	did.DELETE("/:did", didHandler.DeleteDID)
		//}

		// SBOM 工具路由

		sbom := v1.Group("/sbom")
		{
			// 创建 SBOM
			sbom.POST("/create", sbomHandler.CreateSBOM)

			// 获取指定 ID 的 SBOM
			sbom.GET("/:id", sbomHandler.GetSBOMFromBlockchain)

			// 根据 DID 列出所有 SBOM
			sbom.GET("/did/:did", sbomHandler.GetSBOMsByDIDFromBlockchain)

			// 保存 SBOM 到区块链
			sbom.POST("/blockchain/save", sbomHandler.SaveSBOMToBlockchain)

			// 加载漏洞库
			sbom.POST("/vuln/load", sbomHandler.LoadVulnerabilityDatabase)

			// 扫描漏洞
			sbom.POST("/vuln/scan", sbomHandler.ScanForVulnerabilities)
		}
    
		// Vuln 工具路由
		vuln := v1.Group("/vuln")
		{
			// 报告新的漏洞
			vuln.POST("/report", vulnHandler.ReportVulnerability)

			// 根据 ID 获取漏洞信息
			vuln.GET("/:id", vulnHandler.GetVulnerability)

			// 列出漏洞信息
			vuln.GET("/list", vulnHandler.ListVulnerabilities)

			// 根据组件获取漏洞信息
			vuln.GET("/component/:component", vulnHandler.GetVulnerabilitiesByComponent)

			// 搜索漏洞信息
			vuln.GET("/search", vulnHandler.SearchVulnerabilities)
		}

	}
}
