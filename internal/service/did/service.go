package did

import (
	"blockSBOM/internal/dal/dal/model"
	"blockSBOM/internal/dal/dal/query"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type DIDService struct {
	contract *contracts.DIDContract
	repo     *query.DIDRepository
}

func NewDIDService(contract *contracts.DIDContract, repo *query.DIDRepository) *DIDService {
	return &DIDService{
		contract: contract,
		repo:     repo,
	}
}

func (s *DIDService) CreateDID(ctx context.Context, req *CreateDIDRequest) (*model.DIDDocument, error) {
	id := fmt.Sprintf("did:blocksbom:%s", uuid.New().String())

	doc := &model.DIDDocument{
		ID:             id,
		PublicKey:      req.PublicKeys,
		Authentication: req.Authentication,
		Created:        time.Now().UTC(),
		Updated:        time.Now().UTC(),
	}

	// 先写入区块链
	if err := s.contract.CreateDID(id, doc); err != nil {
		return nil, fmt.Errorf("创建区块链DID失败: %v", err)
	}

	// 再写入数据库
	if err := s.repo.CreateDID(ctx, doc); err != nil {
		return nil, fmt.Errorf("创建数据库DID失败: %v", err)
	}

	return doc, nil
}

func (s *DIDService) ResolveDID(ctx context.Context, id string) (*model.DIDDocument, error) {
	// 优先从数据库查询
	doc, err := s.repo.GetDID(ctx, id)
	if err == nil {
		return doc, nil
	}

	// 数据库查询失败，从区块链获取
	doc, err = s.contract.ResolveDID(id)
	if err != nil {
		return nil, fmt.Errorf("解析DID失败: %v", err)
	}

	// 同步到数据库
	if err := s.repo.CreateDID(ctx, doc); err != nil {
		// 仅记录错误，不影响返回结果
		fmt.Printf("同步DID到数据库失败: %v\n", err)
	}

	return doc, nil
}

func (s *DIDService) ListDIDs(ctx context.Context, offset, limit int) ([]*model.DIDDocument, int64, error) {
	return s.repo.ListDIDs(ctx, offset, limit)
}

func (s *DIDService) UpdateDID(ctx context.Context, id string, req *UpdateDIDRequest) (*model.DIDDocument, error) {
	// 先检查DID是否存在
	doc, err := s.ResolveDID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("DID不存在: %v", err)
	}

	// 更新文档
	doc.PublicKey = req.PublicKeys
	doc.Authentication = req.Authentication
	doc.Updated = time.Now().UTC()

	// 先更新区块链
	if err := s.contract.CreateDID(id, doc); err != nil {
		return nil, fmt.Errorf("更新区块链DID失败: %v", err)
	}

	// 再更新数据库
	if err := s.repo.UpdateDID(ctx, doc); err != nil {
		return nil, fmt.Errorf("更新数据库DID失败: %v", err)
	}

	return doc, nil
}

func (s *DIDService) ValidateDID(ctx context.Context, id string) (bool, error) {
	doc, err := s.ResolveDID(ctx, id)
	if err != nil {
		return false, nil
	}
	return doc != nil, nil
}

// Request/Response types
type CreateDIDRequest struct {
	PublicKeys     []string `json:"publicKeys" binding:"required"`
	Authentication []string `json:"authentication" binding:"required"`
}

type UpdateDIDRequest struct {
	PublicKeys     []string `json:"publicKeys" binding:"required"`
	Authentication []string `json:"authentication" binding:"required"`
}
