package handlers

import (
	"blockSBOM/internal/service/vuln"
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type VulnHandler struct {
	vulnService *vuln.VulnService
}

func NewVulnHandler(vulnService *vuln.VulnService) *VulnHandler {
	return &VulnHandler{
		vulnService: vulnService,
	}
}

// LoadVulnerabilityDatabase 加载漏洞库
func (h *VulnHandler) LoadVulnerabilityDatabase(c context.Context, ctx *app.RequestContext) {
	var reply []vuln.Vulnerability
	err := h.vulnService.LoadVulnerabilityDatabase(nil, &reply)
	if err != nil {
		ctx.JSON(consts.StatusInternalServerError, map[string]string{"error": "加载漏洞库失败", "message": err.Error()})
		return
	}
	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"message": "漏洞库加载成功",
		"data":    reply,
	})
}

// GetVulnerability 根据 ID 获取漏洞信息
func (h *VulnHandler) GetVulnerability(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	vulnerability, err := h.vulnService.GetVulnerability(c, id)
	if err != nil {
		ctx.JSON(consts.StatusNotFound, ErrorResponse("获取漏洞信息失败", err))
		return
	}

	ctx.JSON(consts.StatusOK, SuccessResponse("获取漏洞信息成功", vulnerability))
}

// ListVulnerabilities 列出漏洞信息
func (h *VulnHandler) ListVulnerabilities(c context.Context, ctx *app.RequestContext) {
	severity := ctx.DefaultQuery("severity", "")
	offset := ctx.DefaultQuery("offset", "0")
	limit := ctx.DefaultQuery("limit", "10")

	offsetInt, limitInt := ParsePagination(offset, limit)

	vulnerabilities, total, err := h.vulnService.ListVulnerabilities(c, severity, offsetInt, limitInt)
	if err != nil {
		ctx.JSON(consts.StatusInternalServerError, ErrorResponse("获取漏洞列表失败", err))
		return
	}

	ctx.JSON(consts.StatusOK, PageResponse("获取漏洞列表成功", vulnerabilities, total))
}

// SearchVulnerabilities 搜索漏洞信息
func (h *VulnHandler) SearchVulnerabilities(c context.Context, ctx *app.RequestContext) {
	keyword := ctx.Query("keyword")
	offset := ctx.DefaultQuery("offset", "0")
	limit := ctx.DefaultQuery("limit", "10")

	offsetInt, limitInt := ParsePagination(offset, limit)

	vulnerabilities, total, err := h.vulnService.SearchVulnerabilities(c, keyword, offsetInt, limitInt)
	if err != nil {
		ctx.JSON(consts.StatusInternalServerError, ErrorResponse("搜索漏洞失败", err))
		return
	}

	ctx.JSON(consts.StatusOK, PageResponse("搜索漏洞成功", vulnerabilities, total))
}
