package sbom

import (
	"context"
	"encoding/json"
	"fmt"

	"blockSBOM/internal/blockchain/contracts/sbom"
	"blockSBOM/internal/dal/model"
	"blockSBOM/internal/dal/query"

	"github.com/google/uuid"
)

type SBOMService struct {
	contract *sbom.SmartContract
	repo     *query.SBOMRepository
}

func NewSBOMService(contract *sbom.SmartContract, repo *query.SBOMRepository) *SBOMService {
	return &SBOMService{
		contract: contract,
		repo:     repo,
	}
}

// 定义 CreateSBOMRequest 结构体
type CreateSBOMRequest struct {
	DID      string          `json:"did"`
	SPDXSBOM *model.SPDXSBOM `json:"spdxSBOM,omitempty"`
	CDXSBOM  *model.CDXSBOM  `json:"cdxSBOM,omitempty"`
}

func (s *SBOMService) CreateSBOM(ctx context.Context, req *CreateSBOMRequest) (*model.SBOM, error) {
	// 创建 SBOM 实例
	sbom := &model.SBOM{
		ID:       uuid.New().String(),
		DID:      req.DID,
		SPDXSBOM: req.SPDXSBOM,
		CDXSBOM:  req.CDXSBOM,
	}

	// 设置 SBOM 格式和内容
	if sbom.SPDXSBOM != nil {
		sbom.Format = "spdx"
		content, err := json.Marshal(sbom.SPDXSBOM)
		if err != nil {
			return nil, fmt.Errorf("序列化 SPDX SBOM 失败: %v", err)
		}
		sbom.Content = content
	} else if sbom.CDXSBOM != nil {
		sbom.Format = "cdx"
		content, err := json.Marshal(sbom.CDXSBOM)
		if err != nil {
			return nil, fmt.Errorf("序列化 CycloneDX SBOM 失败: %v", err)
		}
		sbom.Content = content
	} else {
		return nil, fmt.Errorf("必须提供 SPDX 或 CycloneDX SBOM 数据")
	}

	// 先写入区块链
	doc, err := json.Marshal(sbom)
	if err != nil {
		return nil, fmt.Errorf("序列化 SBOM 失败: %v", err)
	}

	if err := (*s.contract).StoreSBOM(ctx, sbom.ID, string(doc)); err != nil {
		return nil, fmt.Errorf("存储区块链 SBOM 失败: %v", err)
	}

	// 再写入数据库
	if err := s.repo.CreateSBOM(ctx, sbom); err != nil {
		return nil, fmt.Errorf("存储数据库 SBOM 失败: %v", err)
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
	sbomDoc, err := (*s.contract).GetSBOM(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("获取 SBOM 失败: %v", err)
	}

	// 反序列化 SBOM 文档
	var chainSBOM model.SBOM
	if err := json.Unmarshal([]byte(sbomDoc), &chainSBOM); err != nil {
		return nil, fmt.Errorf("反序列化 SBOM 失败: %v", err)
	}

	// 根据格式反序列化嵌套的 SBOM 数据
	if chainSBOM.SPDXSBOM != nil {
		spdxBytes, err := json.Marshal(chainSBOM.SPDXSBOM)
		if err != nil {
			return nil, fmt.Errorf("序列化 SPDX SBOM 失败: %v", err)
		}
		var spdxSBOM model.SPDXSBOM
		if err := json.Unmarshal(spdxBytes, &spdxSBOM); err != nil {
			return nil, fmt.Errorf("反序列化 SPDX SBOM 失败: %v", err)
		}
		chainSBOM.SPDXSBOM = &spdxSBOM
	} else if chainSBOM.CDXSBOM != nil {
		cdxBytes, err := json.Marshal(chainSBOM.CDXSBOM)
		if err != nil {
			return nil, fmt.Errorf("序列化 CycloneDX SBOM 失败: %v", err)
		}
		var cdxSBOM model.CDXSBOM
		if err := json.Unmarshal(cdxBytes, &cdxSBOM); err != nil {
			return nil, fmt.Errorf("反序列化 CycloneDX SBOM 失败: %v", err)
		}
		chainSBOM.CDXSBOM = &cdxSBOM
	}

	// 同步到数据库
	if err := s.repo.CreateSBOM(ctx, &model.SBOM{
		ID:       chainSBOM.ID,
		DID:      chainSBOM.DID,
		SPDXSBOM: chainSBOM.SPDXSBOM,
		CDXSBOM:  chainSBOM.CDXSBOM,
	}); err != nil {
		fmt.Printf("同步 SBOM 到数据库失败: %v\n", err)
	}

	return &model.SBOM{
		ID:       chainSBOM.ID,
		DID:      chainSBOM.DID,
		SPDXSBOM: chainSBOM.SPDXSBOM,
		CDXSBOM:  chainSBOM.CDXSBOM,
	}, nil
}

func (s *SBOMService) ListSBOMsByDID(ctx context.Context, did string, offset, limit int) ([]*model.SBOM, int64, error) {
	return s.repo.ListSBOMsByDID(ctx, did, offset, limit)
}

func (s *SBOMService) SearchSBOMs(ctx context.Context, keyword string, offset, limit int) ([]*model.SBOM, int64, error) {
	return s.repo.SearchSBOMs(ctx, keyword, offset, limit)
}
