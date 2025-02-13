package middleware

import (
	"blockSBOM/internal/api/api/handlers"
	"blockSBOM/pkg/utils"
	"context"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// Auth 认证中间件
func Auth() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		auth := string(ctx.GetHeader("Authorization"))
		if auth == "" {
			ctx.JSON(consts.StatusUnauthorized, handlers.ErrorResponse("未提供认证令牌", nil))
			ctx.Abort()
			return
		}

		// 检查并提取 Bearer 令牌
		parts := strings.Split(auth, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.JSON(consts.StatusUnauthorized, handlers.ErrorResponse("无效的认证格式", nil))
			ctx.Abort()
			return
		}

		// 使用 ParseToken 替代 ValidateAccessToken
		claims, err := utils.GetJWTHandler().ParseToken(parts[1], utils.AccessToken)
		if err != nil {
			ctx.JSON(consts.StatusUnauthorized, handlers.ErrorResponse("无效的认证令牌", err))
			ctx.Abort()
			return
		}

		// 将用户信息存储在上下文中
		ctx.Set("userID", claims.UserID)
		ctx.Set("username", claims.Username)

		ctx.Next(c)
	}
}
