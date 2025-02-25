package handlers

import (
	"blockSBOM/internal/service/sbom"
	"context"
	"encoding/json"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type SBOMHandler struct {
	sbomService *sbom.SBOMService
}

func NewSBOMHandler(sbomService *sbom.SBOMService) *SBOMHandler {
	return &SBOMHandler{
		sbomService: sbomService,
	}
}

// CreateSBOM 创建 SBOM
func (h *SBOMHandler) CreateSBOM(c context.Context, ctx *app.RequestContext) {
	var req sbom.CreateSBOMRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]string{"error": "无效的请求参数", "message": err.Error()})
		return
	}

	// 调用 SBOMService 的 GenerateSBOM 方法
	args := sbom.Args{
		Language:    req.Language,
		Format:      req.Format,
		ProjectPath: req.ProjectPath,
		ConfigPath:  req.ConfigPath,
		Token:       req.Token,
	}
	var reply string
	err := h.sbomService.GenerateSBOM(&args, &reply)
	if err != nil {
		ctx.JSON(consts.StatusInternalServerError, map[string]string{"error": "创建 SBOM 失败", "message": err.Error()})
		return
	}

	ctx.JSON(consts.StatusCreated, map[string]string{"message": reply})
}

// LoadVulnerabilityDatabase 加载漏洞库
func (h *SBOMHandler) LoadVulnerabilityDatabase(c context.Context, ctx *app.RequestContext) {
	var reply []sbom.Vulnerability
	err := h.sbomService.LoadVulnerabilityDatabase(nil, &reply)
	if err != nil {
		ctx.JSON(consts.StatusInternalServerError, map[string]string{"error": "加载漏洞库失败", "message": err.Error()})
		return
	}

	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"message": "漏洞库加载成功",
		"data":    reply,
	})
}

// ScanForVulnerabilities 扫描漏洞
func (h *SBOMHandler) ScanForVulnerabilities(c context.Context, ctx *app.RequestContext) {
	var req sbom.ScanVulnerabilitiesRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]string{"error": "无效的请求参数", "message": err.Error()})
		return
	}

	args := sbom.Args{
		Language:    req.Language,
		Format:      req.Format,
		PackagePath: req.PackagePath,
		ConfigPath:  req.ConfigPath,
		Token:       req.Token,
	}
	var reply string
	err := h.sbomService.ScanForVulnerabilities(&args, &reply)
	if err != nil {
		ctx.JSON(consts.StatusInternalServerError, map[string]string{"error": "扫描漏洞失败", "message": err.Error()})
		return
	}

	ctx.JSON(consts.StatusOK, map[string]string{"message": reply})
}

// SaveSBOMToBlockchain 将 SBOM 保存到区块链
func (h *SBOMHandler) SaveSBOMToBlockchain(c context.Context, ctx *app.RequestContext) {
	var req sbom.SBOM
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(consts.StatusBadRequest, map[string]string{"error": "无效的请求参数", "message": err.Error()})
		return
	}

	// 将 SBOM 数据序列化为 JSON
	sbomData, err := json.Marshal(req)
	if err != nil {
		ctx.JSON(consts.StatusInternalServerError, map[string]string{"error": "序列化 SBOM 数据失败", "message": err.Error()})
		return
	}

	// 调用 SBOMService 的 SaveSBOMToBlockchain 方法
	sbomID, err := h.sbomService.SaveSBOMToBlockchain(string(sbomData))
	if err != nil {
		ctx.JSON(consts.StatusInternalServerError, map[string]string{"error": "保存 SBOM 到区块链失败", "message": err.Error()})
		return
	}

	ctx.JSON(consts.StatusCreated, map[string]string{"message": "SBOM 保存成功", "id": sbomID})
}

// GetSBOMFromBlockchain 根据 ID 从区块链获取 SBOM
func (h *SBOMHandler) GetSBOMFromBlockchain(c context.Context, ctx *app.RequestContext) {
	sbomID := ctx.Param("id")
	if sbomID == "" {
		ctx.JSON(consts.StatusBadRequest, map[string]string{"error": "无效的 SBOM ID"})
		return
	}

	// 调用 SBOMService 的 GetSBOMFromBlockchain 方法
	sbomData, err := h.sbomService.GetSBOMFromBlockchain(sbomID)
	if err != nil {
		ctx.JSON(consts.StatusNotFound, map[string]string{"error": "获取 SBOM 失败", "message": err.Error()})
		return
	}

	ctx.JSON(consts.StatusOK, map[string]string{"message": "获取 SBOM 成功", "data": sbomData})
}

// GetSBOMsByDIDFromBlockchain 根据 DID 从区块链获取所有 SBOM
func (h *SBOMHandler) GetSBOMsByDIDFromBlockchain(c context.Context, ctx *app.RequestContext) {
	did := ctx.Param("did")
	if did == "" {
		ctx.JSON(consts.StatusBadRequest, map[string]string{"error": "无效的 DID"})
		return
	}

	// 调用 SBOMService 的 GetSBOMsByDIDFromBlockchain 方法
	sboms, err := h.sbomService.GetSBOMsByDIDFromBlockchain(did)
	if err != nil {
		ctx.JSON(consts.StatusInternalServerError, map[string]string{"error": "获取 SBOM 列表失败", "message": err.Error()})
		return
	}

	ctx.JSON(consts.StatusOK, map[string]interface{}{
		"message": "获取 SBOM 列表成功",
		"data":    sboms,
	})
}
