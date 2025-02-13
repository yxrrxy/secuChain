package auth

import (
	"blockSBOM/internal/dal/dal/model"
	"blockSBOM/internal/dal/dal/query"
	"blockSBOM/pkg/utils"
	"context"
	"errors"
	"fmt"
)

type AuthService struct {
	repo       *query.UserRepository
	jwtHandler *utils.JWTHandler
}

func NewAuthService(repo *query.UserRepository, jwtHandler *utils.JWTHandler) *AuthService {
	return &AuthService{
		repo:       repo,
		jwtHandler: jwtHandler,
	}
}

func (s *AuthService) Register(ctx context.Context, req *RegisterRequest) (*model.User, error) {
	// 检查用户是否已存在
	if err := s.checkUserExists(ctx, req.Username, req.Email); err != nil {
		return nil, err
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败: %v", err)
	}

	// 创建用户
	user := &model.User{
		Username: req.Username,
		Password: hashedPassword,
		Email:    req.Email,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("创建用户失败: %v", err)
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	// 获取用户
	user, err := s.repo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 验证密码
	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("用户名或密码错误")
	}

	// 生成token对
	tokenPair, err := s.jwtHandler.GenerateTokenPair(user.ID, user.Username)
	if err != nil {
		return nil, fmt.Errorf("生成token失败: %v", err)
	}

	// 更新最后登录时间
	if err := s.repo.UpdateLastLogin(ctx, user.ID); err != nil {
		// 仅记录错误，不影响登录
		fmt.Printf("更新最后登录时间失败: %v\n", err)
	}

	return &LoginResponse{
		Token:        tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
		Username:     user.Username,
		Email:        user.Email,
	}, nil
}

func (s *AuthService) checkUserExists(ctx context.Context, username, email string) error {
	// 检查用户名
	exists, err := s.repo.ExistsByUsername(ctx, username)
	if err != nil {
		return fmt.Errorf("检查用户名失败: %v", err)
	}
	if exists {
		return errors.New("用户名已存在")
	}

	// 检查邮箱
	exists, err = s.repo.ExistsByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("检查邮箱失败: %v", err)
	}
	if exists {
		return errors.New("邮箱已被注册")
	}

	return nil
}

func (s *AuthService) GetUserByID(ctx context.Context, id uint) (*model.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("获取用户失败: %v", err)
	}
	return user, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, token string) (*utils.TokenPair, error) {
	// 使用 ParseToken 验证刷新令牌
	claims, err := s.jwtHandler.ParseToken(token, utils.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("无效的token: %v", err)
	}

	// 生成新的令牌对
	return s.jwtHandler.GenerateTokenPair(claims.UserID, claims.Username)
}

// Request/Response types
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=32"`
	Email    string `json:"email" binding:"required,email"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	Username     string `json:"username"`
	Email        string `json:"email"`
}
