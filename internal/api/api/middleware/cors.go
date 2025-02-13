package middleware

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

// Cors 跨域中间件
func Cors() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin,Content-Type,Accept,Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")

		if string(c.Request.Method()) == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next(ctx)
	}
}
