package handlers

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type ManagementHandler struct {
	managementService *management.ManagementService
}

func NewManagementHandler(managementService *management.ManagementService) *ManagementHandler {
	return &ManagementHandler{
		managementService: managementService,
	}
}

func (h *ManagementHandler) GetOverview(c context.Context, ctx *app.RequestContext) {
	overview, err := h.managementService.GetOverview(c)
	if err != nil {
		ctx.JSON(consts.StatusInternalServerError, ErrorResponse("获取概览失败", err))
		return
	}

	ctx.JSON(consts.StatusOK, SuccessResponse("获取概览成功", overview))
}

func (h *ManagementHandler) UploadPackage(c context.Context, ctx *app.RequestContext) {
	// 处理文件上传和保存
}

func (h *ManagementHandler) ListPackages(c context.Context, ctx *app.RequestContext) {
	packages, err := h.managementService.ListPackages(c)
	if err != nil {
		ctx.JSON(consts.StatusInternalServerError, ErrorResponse("获取软件包列表失败", err))
		return
	}

	ctx.JSON(consts.StatusOK, SuccessResponse("获取软件包列表成功", packages))
}

func (h *ManagementHandler) GetPackageDetails(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	packageDetails, err := h.managementService.GetPackageDetails(c, id)
	if err != nil {
		ctx.JSON(consts.StatusNotFound, ErrorResponse("获取软件包详情失败", err))
		return
	}

	ctx.JSON(consts.StatusOK, SuccessResponse("获取软件包详情成功", packageDetails))
}

func (h *ManagementHandler) ScanPackage(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	scanResult, err := h.managementService.ScanPackage(c, id)
	if err != nil {
		ctx.JSON(consts.StatusInternalServerError, ErrorResponse("扫描软件包失败", err))
		return
	}

	ctx.JSON(consts.StatusOK, SuccessResponse("扫描软件包成功", scanResult))
}
