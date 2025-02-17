package handlers

import (
	"context"
	"io"
	"mime/multipart"
	"os"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func (h *ManagementHandler) UploadPackage(c context.Context, ctx *app.RequestContext) {
	// 示例：处理文件上传
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(consts.StatusBadRequest, ErrorResponse("上传文件失败", err))
		return
	}
	defer file.Close()

	// 示例：保存文件到本地
	filePath := "uploads/" + header.Filename
	if err := saveFile(file, filePath); err != nil {
		ctx.JSON(consts.StatusInternalServerError, ErrorResponse("保存文件失败", err))
		return
	}

	// 示例：调用管理服务处理上传的文件
	err = h.managementService.UploadPackage(c, filePath)
	if err != nil {
		ctx.JSON(consts.StatusInternalServerError, ErrorResponse("处理文件失败", err))
		return
	}

	ctx.JSON(consts.StatusOK, SuccessResponse("文件上传成功", nil))
}

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
