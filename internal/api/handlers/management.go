package handlers

import (
	"context"
	"io"
	"mime/multipart"
	"os"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// ManagementHandler 处理管理相关的请求
type ManagementHandler struct {
	managementService ManagementService
}

// ManagementService 定义管理服务接口
type ManagementService interface {
	UploadPackage(ctx context.Context, file *multipart.FileHeader) error
	GetOverview(ctx context.Context) (interface{}, error)
}

// NewManagementHandler 创建一个新的 ManagementHandler 实例
func NewManagementHandler(service ManagementService) *ManagementHandler {
	return &ManagementHandler{
		managementService: service,
	}
}

// UploadPackage 处理文件上传请求
func (h *ManagementHandler) UploadPackage(ctx *app.RequestContext) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(consts.StatusBadRequest, ErrorResponse("获取上传文件失败", err))
		return
	}
	//先注释
	_ = file
	//err = h.managementService.UploadPackage(ctx.Request.Context(), file)
	if err != nil {
		ctx.JSON(consts.StatusInternalServerError, ErrorResponse("上传文件失败", err))
		return
	}

	ctx.JSON(consts.StatusOK, SuccessResponse("上传文件成功", nil))
}

// saveFile 将上传的文件保存到指定路径
func saveFile(file multipart.File, filePath string) error {
	// 示例：保存文件到指定路径
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	return err
}

func (h *ManagementHandler) GetOverview(c context.Context, ctx *app.RequestContext) {
	overview, err := h.managementService.GetOverview(c)
	if err != nil {
		ctx.JSON(consts.StatusInternalServerError, ErrorResponse("获取概览失败", err))
		return
	}
	ctx.JSON(consts.StatusOK, SuccessResponse("获取概览成功", overview))
}
