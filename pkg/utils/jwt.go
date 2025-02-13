package utils

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var (
	jwtHandlerInstance *JWTHandler
	jwtHandlerOnce     sync.Once
)

// TokenType 定义令牌类型
type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

// Claims 定义 JWT 的声明结构
type Claims struct {
	UserID   uint      `json:"user_id"`
	Username string    `json:"username"`
	Type     TokenType `json:"type"`
	jwt.RegisteredClaims
}

// JWTHandler JWT 处理器
type JWTHandler struct {
	secret        []byte
	accessExpiry  time.Duration // 访问令牌过期时间
	refreshExpiry time.Duration // 刷新令牌过期时间
}

// TokenPair 包含访问令牌和刷新令牌
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"` // 访问令牌过期时间（秒）
}

// GetJWTHandler 获取 JWT 处理器的全局实例
func GetJWTHandler() *JWTHandler {
	if jwtHandlerInstance == nil {
		panic("JWT handler not initialized")
	}
	return jwtHandlerInstance
}

// InitJWTHandler 初始化 JWT 处理器
func InitJWTHandler(secret string) {
	jwtHandlerOnce.Do(func() {
		jwtHandlerInstance = NewJWTHandler(secret)
	})
}

// NewJWTHandler 创建新的 JWT 处理器
func NewJWTHandler(secret string) *JWTHandler {
	if secret == "" {
		panic("JWT secret cannot be empty")
	}
	return &JWTHandler{
		secret:        []byte(secret),
		accessExpiry:  2 * time.Hour,      // 访问令牌2小时过期
		refreshExpiry: 7 * 24 * time.Hour, // 刷新令牌7天过期
	}
}

// GenerateTokenPair 生成访问令牌和刷新令牌对
func (h *JWTHandler) GenerateTokenPair(userID uint, username string) (*TokenPair, error) {
	// 生成访问令牌
	accessToken, err := h.GenerateToken(userID, username, AccessToken, h.accessExpiry)
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	// 生成刷新令牌
	refreshToken, err := h.GenerateToken(userID, username, RefreshToken, h.refreshExpiry)
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(h.accessExpiry.Seconds()),
	}, nil
}

// GenerateToken 生成指定类型的令牌
func (h *JWTHandler) GenerateToken(userID uint, username string, tokenType TokenType, expiry time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:   userID,
		Username: username,
		Type:     tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(h.secret)
}

// ParseToken 解析并验证令牌
func (h *JWTHandler) ParseToken(tokenString string, expectedType TokenType) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return h.secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	// 验证令牌类型
	if claims.Type != expectedType {
		return nil, fmt.Errorf("invalid token type: expected %s, got %s", expectedType, claims.Type)
	}

	// 检查是否过期
	if time.Now().After(claims.ExpiresAt.Time) {
		return nil, errors.New("token has expired")
	}

	return claims, nil
}

// RefreshTokenPair 使用刷新令牌生成新的令牌对
func (h *JWTHandler) RefreshTokenPair(refreshToken string) (*TokenPair, error) {
	// 解析并验证刷新令牌
	claims, err := h.ParseToken(refreshToken, RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// 生成新的令牌对
	return h.GenerateTokenPair(claims.UserID, claims.Username)
}

// ValidateAccessToken 验证访问令牌
func (h *JWTHandler) ValidateAccessToken(tokenString string) (*Claims, error) {
	return h.ParseToken(tokenString, AccessToken)
}
