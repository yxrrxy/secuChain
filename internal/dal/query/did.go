package query

import (
	"blockSBOM/internal/dal/model"
	"context"

	"gorm.io/gorm"
)

type DIDRepository struct {
	db *gorm.DB
}

func NewDIDRepository(db *gorm.DB) *DIDRepository {
	return &DIDRepository{db: db}
}

func (r *DIDRepository) CreateDID(ctx context.Context, doc *model.DIDDocument) error {
	return r.db.WithContext(ctx).Create(doc).Error
}

func (r *DIDRepository) GetDID(ctx context.Context, id string) (*model.DIDDocument, error) {
	var doc model.DIDDocument
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&doc).Error; err != nil {
		return nil, err
	}
	return &doc, nil
}

func (r *DIDRepository) ListDIDs(ctx context.Context, offset, limit int) ([]*model.DIDDocument, int64, error) {
	var docs []*model.DIDDocument
	var total int64

	if err := r.db.WithContext(ctx).Model(&model.DIDDocument{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&docs).Error; err != nil {
		return nil, 0, err
	}

	return docs, total, nil
}

func (r *DIDRepository) UpdateDID(ctx context.Context, doc *model.DIDDocument) error {
	return r.db.WithContext(ctx).Save(doc).Error
}
