package handlers

import (
	"blockSBOM/internal/service/did"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type DIDHandler struct {
	didService *did.DIDService
}

func NewDIDHandler(didService *did.DIDService) *DIDHandler {
	return &DIDHandler{
		didService: didService,
	}
}

// RegisterDID 处理创建 DID 的 API 请求
func (h *DIDHandler) RegisterDID(c context.Context, ctx *app.RequestContext) {
	// 调用服务层的 RegisterDID 方法
	didData, err := h.didService.RegisterDID(c)
	if err != nil {
		// 调用服务失败，返回 500 错误
		ctx.JSON(consts.StatusInternalServerError, ErrorResponse("DID 創建失敗", err))
		return
	}
	// 成功响应
	ctx.JSON(consts.StatusOK, SuccessResponse("DID 創建成功", didData))
}

// UpdateDID 用于更新 DID 的 API 方法
func (h *DIDHandler) UpdateDID(c context.Context, ctx *app.RequestContext) {
	var req UpdateDIDRequest // 假设你有一个更新DID请求的结构体
	// 绑定并验证请求体中的数据
	if err := ctx.BindAndValidate(&req); err != nil {
		// 请求参数无效，返回 400 错误
		ctx.JSON(consts.StatusBadRequest, ErrorResponse("无效的请求参数", err))
		return
	}

	// 调用服务层的 UpdateDID 方法
	result, err := h.didService.UpdateDID(c, req.DID, req.RecoveryKey, req.RecoveryPrivateKey)
	if err != nil {
		// 调用服务失败，返回 500 错误
		ctx.JSON(consts.StatusInternalServerError, ErrorResponse("更新 DID 失败", err))
		return
	}

	// 成功响应
	ctx.JSON(consts.StatusOK, SuccessResponse("DID 更新成功", result))
}

// DeleteDID 用于删除 DID 的 API 方法
func (h *DIDHandler) DeleteDID(c context.Context, ctx *app.RequestContext) {
	// 定义请求体结构体
	var req DeleteDIDRequest
	// 绑定并验证请求体中的数据
	if err := ctx.BindAndValidate(&req); err != nil {
		// 请求参数无效，返回 400 错误
		ctx.JSON(consts.StatusBadRequest, ErrorResponse("无效的请求参数", err))
		return
	}

	// 调用服务层的 DeleteDID 方法
	result, err := h.didService.DeleteDID(c, req.DID, req.RecoveryKey, req.RecoveryPrivateKey)
	if err != nil {
		// 调用服务失败，返回 500 错误
		ctx.JSON(consts.StatusInternalServerError, ErrorResponse("刪除 DID 失败", err))
		return
	}

	// 成功响应
	ctx.JSON(consts.StatusOK, SuccessResponse("DID 刪除成功", result))
}

func (h *DIDHandler) ResolveDIDAPI(c context.Context, ctx *app.RequestContext) {
	// 获取请求参数中的 DID
	did := ctx.Param("did")
	if did == "" {
		ctx.JSON(consts.StatusBadRequest, ErrorResponse("无效的请求参数", nil))
		return
	}

	// 调用服务层的 ResolveDID 方法
	doc, err := h.didService.ResolveDID(c, did)
	if err != nil {
		// 调用服务失败，返回 500 错误
		ctx.JSON(consts.StatusInternalServerError, ErrorResponse("解析 DID 失败", err))
		return
	}

	// 成功响应，返回 DID 文档
	ctx.JSON(consts.StatusOK, SuccessResponse("DID 更新成功", doc))
}
