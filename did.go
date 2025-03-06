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

func (r *DIDRepository) CreateDID(ctx context.Context, doc *model.DID) error {
	return r.db.WithContext(ctx).Create(doc).Error
}

func (r *DIDRepository) GetDID(ctx context.Context, id string) (*model.DID, error) {
	var doc model.DID
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&doc).Error; err != nil {
		return nil, err
	}
	return &doc, nil
}
func (r *DIDRepository) DeleteDID(ctx context.Context, id string) error {
	// 查找 DID
	var doc model.DID
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&doc).Error; err != nil {
		// 如果 DID 不存在，返回错误
		return err
	}
	// 删除 DID
	if err := r.db.WithContext(ctx).Delete(&doc).Error; err != nil {
		// 如果删除失败，返回错误
		return err
	}
	return nil
}

func (r *DIDRepository) ListDIDs(ctx context.Context, offset, limit int) ([]*model.DID, int64, error) {
	var docs []*model.DID
	var total int64

	if err := r.db.WithContext(ctx).Model(&model.DIDDocument{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&docs).Error; err != nil {
		return nil, 0, err
	}

	return docs, total, nil
}

func (r *DIDRepository) UpdateDID(ctx context.Context, doc *model.DID) error {
	return r.db.WithContext(ctx).Save(doc).Error
}
