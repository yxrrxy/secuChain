package query

import (
	"blockSBOM/internal/dal/model"
	"context"

	"gorm.io/gorm"
)

type SBOMRepository struct {
	db *gorm.DB
}

func NewSBOMRepository(db *gorm.DB) *SBOMRepository {
	return &SBOMRepository{db: db}
}

func (r *SBOMRepository) CreateSBOM(ctx context.Context, sbom *model.SBOM) error {
	return r.db.WithContext(ctx).Create(sbom).Error
}

func (r *SBOMRepository) GetSBOM(ctx context.Context, id string) (*model.SBOM, error) {
	var sbom model.SBOM
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&sbom).Error; err != nil {
		return nil, err
	}
	return &sbom, nil
}

func (r *SBOMRepository) ListSBOMsByDID(ctx context.Context, did string, offset, limit int) ([]*model.SBOM, int64, error) {
	var sboms []*model.SBOM
	var total int64

	query := r.db.WithContext(ctx).Where("did = ?", did)

	if err := query.Model(&model.SBOM{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(limit).Find(&sboms).Error; err != nil {
		return nil, 0, err
	}

	return sboms, total, nil
}

func (r *SBOMRepository) SearchSBOMs(ctx context.Context, keyword string, offset, limit int) ([]*model.SBOM, int64, error) {
	var sboms []*model.SBOM
	var total int64

	query := r.db.WithContext(ctx).Where(
		"name LIKE ? OR version LIKE ? OR components LIKE ?",
		"%"+keyword+"%",
		"%"+keyword+"%",
		"%"+keyword+"%",
	)

	if err := query.Model(&model.SBOM{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(limit).Find(&sboms).Error; err != nil {
		return nil, 0, err
	}

	return sboms, total, nil
}
