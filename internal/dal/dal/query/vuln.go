package query

import (
	"blockSBOM/backend/internal/dal/model"
	"context"

	"gorm.io/gorm"
)

type VulnRepository struct {
	db *gorm.DB
}

func NewVulnRepository(db *gorm.DB) *VulnRepository {
	return &VulnRepository{db: db}
}

func (r *VulnRepository) CreateVulnerability(ctx context.Context, vuln *model.Vulnerability) error {
	return r.db.WithContext(ctx).Create(vuln).Error
}

func (r *VulnRepository) GetVulnerability(ctx context.Context, id string) (*model.Vulnerability, error) {
	var vuln model.Vulnerability
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&vuln).Error; err != nil {
		return nil, err
	}
	return &vuln, nil
}

func (r *VulnRepository) ListVulnerabilities(ctx context.Context, severity string, offset, limit int) ([]*model.Vulnerability, int64, error) {
	var vulns []*model.Vulnerability
	var total int64

	query := r.db.WithContext(ctx)
	if severity != "" {
		query = query.Where("severity = ?", severity)
	}

	if err := query.Model(&model.Vulnerability{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(limit).Find(&vulns).Error; err != nil {
		return nil, 0, err
	}

	return vulns, total, nil
}

func (r *VulnRepository) GetVulnerabilitiesByComponent(ctx context.Context, component string) ([]*model.Vulnerability, error) {
	var vulns []*model.Vulnerability
	if err := r.db.WithContext(ctx).
		Where("affected_components LIKE ?", "%"+component+"%").
		Find(&vulns).Error; err != nil {
		return nil, err
	}
	return vulns, nil
}

func (r *VulnRepository) SearchVulnerabilities(ctx context.Context, keyword string, offset, limit int) ([]*model.Vulnerability, int64, error) {
	var vulns []*model.Vulnerability
	var total int64

	query := r.db.WithContext(ctx).Where(
		"cve LIKE ? OR description LIKE ? OR affected_components LIKE ?",
		"%"+keyword+"%",
		"%"+keyword+"%",
		"%"+keyword+"%",
	)

	if err := query.Model(&model.Vulnerability{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(limit).Find(&vulns).Error; err != nil {
		return nil, 0, err
	}

	return vulns, total, nil
}
