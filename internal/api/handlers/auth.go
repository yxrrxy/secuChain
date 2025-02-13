package handlers

import (
	"blockSBOM/backend/internal/service/auth"
	"context"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type AuthHandler struct {
	authService *auth.AuthService
}

func NewAuthHandler(authService *auth.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register 处理注册请求
func (h *AuthHandler) Register(c context.Context, ctx *app.RequestContext) {
	var req auth.RegisterRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(consts.StatusBadRequest, ErrorResponse("无效的请求参数", err))
		return
	}

	user, err := h.authService.Register(c, &req)
	if err != nil {
		ctx.JSON(consts.StatusBadRequest, ErrorResponse("注册失败", err))
		return
	}

	ctx.JSON(consts.StatusCreated, SuccessResponse("注册成功", user))
}

// Login 处理登录请求
func (h *AuthHandler) Login(c context.Context, ctx *app.RequestContext) {
	var req auth.LoginRequest
	if err := ctx.BindAndValidate(&req); err != nil {
		ctx.JSON(consts.StatusBadRequest, ErrorResponse("无效的请求参数", err))
		return
	}

	resp, err := h.authService.Login(c, &req)
	if err != nil {
		ctx.JSON(consts.StatusUnauthorized, ErrorResponse("登录失败", err))
		return
	}

	ctx.JSON(consts.StatusOK, SuccessResponse("登录成功", resp))
}

// GetUserInfo 获取当前用户信息
func (h *AuthHandler) GetUserInfo(c context.Context, ctx *app.RequestContext) {
	// 从上下文中获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(consts.StatusUnauthorized, ErrorResponse("未登录", nil))
		return
	}

	// 使用 service 层获取用户信息
	user, err := h.authService.GetUserByID(c, userID.(uint))
	if err != nil {
		ctx.JSON(consts.StatusNotFound, ErrorResponse("获取用户信息失败", err))
		return
	}

	ctx.JSON(consts.StatusOK, SuccessResponse("获取用户信息成功", user))
}

// RefreshToken 处理令牌刷新请求
func (h *AuthHandler) RefreshToken(c context.Context, ctx *app.RequestContext) {
	token := string(ctx.GetHeader("Authorization"))
	if len(token) == 0 {
		ctx.JSON(consts.StatusUnauthorized, ErrorResponse("未提供token", nil))
		return
	}

	// 从 Bearer token 中提取 token
	parts := strings.Split(token, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		ctx.JSON(consts.StatusUnauthorized, ErrorResponse("无效的token格式", nil))
		return
	}

	// 使用 service 层的 RefreshToken 方法
	tokenPair, err := h.authService.RefreshToken(c, parts[1])
	if err != nil {
		ctx.JSON(consts.StatusUnauthorized, ErrorResponse("刷新token失败", err))
		return
	}

	ctx.JSON(consts.StatusOK, SuccessResponse("刷新token成功", tokenPair))
}
