package handlers

import (
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// ManagementHandler 处理管理相关的请求
type ManagementHandler struct {
	managementService ManagementService
}

// ManagementService 定义管理服务接口
type ManagementService interface {
	UploadPackage(ctx context.Context, filePath string) error
}

// NewManagementHandler 创建一个新的 ManagementHandler 实例
func NewManagementHandler(service ManagementService) *ManagementHandler {
	return &ManagementHandler{
		managementService: service,
	}
}

// UploadPackage 处理文件上传请求
func (h *ManagementHandler) UploadPackage(c context.Context, ctx *gin.Context) {
	// 示例：处理文件上传
	header, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "上传文件失败", "message": err.Error()})
		return
	}

	// 打开文件并处理错误
	file, err := header.Open()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "打开文件失败", "message": err.Error()})
		return
	}
	defer file.Close()

	// 示例：保存文件到本地
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "打开文件失败", "message": err.Error()})
		return
	}
	defer file.Close()

	filePath := "uploads/" + header.Filename
	err = os.MkdirAll("uploads", os.ModePerm)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "创建目录失败", "message": err.Error()})
		return
	}

	out, err := os.Create(filePath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "创建文件失败", "message": err.Error()})
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "文件上传成功"})
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
