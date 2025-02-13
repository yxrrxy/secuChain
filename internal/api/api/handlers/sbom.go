package handlers

import (
	"blockSBOM/backend/internal/service/sbom"
	"context"

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

func (h *SBOMHandler) CreateSBOM(c context.Context, ctx *app.RequestContext) {
	var req sbom.CreateSBOMRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(consts.StatusBadRequest, ErrorResponse("无效的请求参数", err))
		return
	}

	doc, err := h.sbomService.CreateSBOM(c, &req)
	if err != nil {
		ctx.JSON(consts.StatusInternalServerError, ErrorResponse("创建SBOM失败", err))
		return
	}

	ctx.JSON(consts.StatusCreated, SuccessResponse("创建SBOM成功", doc))
}

func (h *SBOMHandler) GetSBOM(c context.Context, ctx *app.RequestContext) {
	id := ctx.Param("id")
	doc, err := h.sbomService.GetSBOM(c, id)
	if err != nil {
		ctx.JSON(consts.StatusNotFound, ErrorResponse("获取SBOM失败", err))
		return
	}

	ctx.JSON(consts.StatusOK, SuccessResponse("获取SBOM成功", doc))
}

func (h *SBOMHandler) ListSBOMsByDID(c context.Context, ctx *app.RequestContext) {
	did := ctx.Param("did")
	offset := ctx.DefaultQuery("offset", "0")
	limit := ctx.DefaultQuery("limit", "10")

	offsetInt, limitInt := ParsePagination(offset, limit)

	docs, total, err := h.sbomService.ListSBOMsByDID(c, did, offsetInt, limitInt)
	if err != nil {
		ctx.JSON(consts.StatusInternalServerError, ErrorResponse("获取SBOM列表失败", err))
		return
	}

	ctx.JSON(consts.StatusOK, PageResponse("获取SBOM列表成功", docs, total))
}

func (h *SBOMHandler) SearchSBOMs(c context.Context, ctx *app.RequestContext) {
	keyword := ctx.Query("keyword")
	offset := ctx.DefaultQuery("offset", "0")
	limit := ctx.DefaultQuery("limit", "10")

	offsetInt, limitInt := ParsePagination(offset, limit)

	docs, total, err := h.sbomService.SearchSBOMs(c, keyword, offsetInt, limitInt)
	if err != nil {
		ctx.JSON(consts.StatusInternalServerError, ErrorResponse("搜索SBOM失败", err))
		return
	}

	ctx.JSON(consts.StatusOK, PageResponse("搜索SBOM成功", docs, total))
}
