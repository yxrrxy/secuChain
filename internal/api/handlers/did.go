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

func (h *DIDHandler) CreateDID(c context.Context, ctx *app.RequestContext) {
	var req did.CreateDIDRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(consts.StatusBadRequest, ErrorResponse("无效的请求参数", err))
		return
	}

	doc, err := h.didService.CreateDID(c, &req)
	if err != nil {
		ctx.JSON(consts.StatusInternalServerError, ErrorResponse("创建DID失败", err))
		return
	}

	ctx.JSON(consts.StatusCreated, SuccessResponse("创建DID成功", doc))
}

func (h *DIDHandler) GetDID(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	doc, err := h.didService.ResolveDID(c, id)
	if err != nil {
		ctx.JSON(consts.StatusNotFound, ErrorResponse("获取DID失败", err))
		return
	}

	ctx.JSON(consts.StatusOK, SuccessResponse("获取DID成功", doc))
}

func (h *DIDHandler) UpdateDID(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	var req did.UpdateDIDRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(consts.StatusBadRequest, ErrorResponse("无效的请求参数", err))
		return
	}

	doc, err := h.didService.UpdateDID(c, id, &req)
	if err != nil {
		ctx.JSON(consts.StatusInternalServerError, ErrorResponse("更新DID失败", err))
		return
	}

	ctx.JSON(consts.StatusOK, SuccessResponse("更新DID成功", doc))
}

func (h *DIDHandler) ListDIDs(c context.Context, ctx *app.RequestContext) {
	offset := ctx.DefaultQuery("offset", "0")
	limit := ctx.DefaultQuery("limit", "10")

	// 转换为整数
	offsetInt, limitInt := ParsePagination(offset, limit)

	docs, total, err := h.didService.ListDIDs(c, offsetInt, limitInt)
	if err != nil {
		ctx.JSON(consts.StatusInternalServerError, ErrorResponse("获取DID列表失败", err))
		return
	}

	ctx.JSON(consts.StatusOK, PageResponse("获取DID列表成功", docs, total))
}
