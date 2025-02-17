package sbom

import (
	"blockSBOM/internal/dal/dal/model"
	"blockSBOM/internal/dal/dal/query"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type SBOMService struct {
	contract *contracts.SBOMContract
	repo     *query.SBOMRepository
}

func NewSBOMService(contract *contracts.SBOMContract, repo *query.SBOMRepository) *SBOMService {
	return &SBOMService{
		contract: contract,
		repo:     repo,
	}
}

func (s *SBOMService) CreateSBOM(ctx context.Context, req *CreateSBOMRequest) (*model.SBOM, error) {
	sbom := &model.SPDXSBOM{
		ID:         uuid.New().String(),
		Name:       req.Name,
		Version:    req.Version,
		Components: req.Components,
		Format:     req.Format,
		Created:    time.Now().UTC(),
		DID:        req.DID,
	}

	// 先写入区块链
	if err := s.contract.StoreSBOM(sbom); err != nil {
		return nil, fmt.Errorf("存储区块链SBOM失败: %v", err)
	}

	// 再写入数据库
	if err := s.repo.CreateSBOM(ctx, sbom); err != nil {
		return nil, fmt.Errorf("存储数据库SBOM失败: %v", err)
	}

	return sbom, nil
}

func (s *SBOMService) GetSBOM(ctx context.Context, id string) (*model.SBOM, error) {
	// 优先从数据库查询
	sbom, err := s.repo.GetSBOM(ctx, id)
	if err == nil {
		return sbom, nil
	}

	// 数据库查询失败，从区块链获取
	sbom, err = s.contract.GetSBOM(id)
	if err != nil {
		return nil, fmt.Errorf("获取SBOM失败: %v", err)
	}

	// 同步到数据库
	if err := s.repo.CreateSBOM(ctx, sbom); err != nil {
		fmt.Printf("同步SBOM到数据库失败: %v\n", err)
	}

	return sbom, nil
}

func (s *SBOMService) ListSBOMsByDID(ctx context.Context, did string, offset, limit int) ([]*model.SBOM, int64, error) {
	return s.repo.ListSBOMsByDID(ctx, did, offset, limit)
}

func (s *SBOMService) SearchSBOMs(ctx context.Context, keyword string, offset, limit int) ([]*model.SBOM, int64, error) {
	return s.repo.SearchSBOMs(ctx, keyword, offset, limit)
}

func (s *SBOMService) ValidateSBOM(ctx context.Context, id string) (bool, error) {
	sbom, err := s.GetSBOM(ctx, id)
	if err != nil {
		return false, nil
	}
	return sbom != nil, nil
}

// Request/Response types
type CreateSBOMRequest struct {
	Name       string   `json:"name" binding:"required"`
	Version    string   `json:"version" binding:"required"`
	Components []string `json:"components" binding:"required"`
	Format     string   `json:"format" binding:"required,oneof=spdx cyclonedx swid"`
	DID        string   `json:"did" binding:"required"`
}

type UpdateSBOMRequest struct {
	Name       string   `json:"name" binding:"required"`
	Version    string   `json:"version" binding:"required"`
	Components []string `json:"components" binding:"required"`
	Format     string   `json:"format" binding:"required,oneof=spdx cyclonedx swid"`
}
