package test

import (
	"context"
	"testing"

	"blockSBOM/internal/dal/model"
	"blockSBOM/internal/service/sbom"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSmartContract struct {
	mock.Mock
}

func (m *MockSmartContract) StoreSBOM(ctx context.Context, id, doc string) error {
	args := m.Called(ctx, id, doc)
	return args.Error(0)
}

func (m *MockSmartContract) GetSBOM(ctx context.Context, id string) (string, error) {
	args := m.Called(ctx, id)
	return args.String(0), args.Error(1)
}

type MockSBOMRepository struct {
	mock.Mock
}

func (m *MockSBOMRepository) CreateSBOM(ctx context.Context, sbom *model.SBOM) error {
	args := m.Called(ctx, sbom)
	return args.Error(0)
}

func (m *MockSBOMRepository) GetSBOM(ctx context.Context, id string) (*model.SBOM, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.SBOM), args.Error(1)
}

func (m *MockSBOMRepository) ListSBOMsByDID(ctx context.Context, did string, offset, limit int) ([]*model.SBOM, int64, error) {
	args := m.Called(ctx, did, offset, limit)
	return args.Get(0).([]*model.SBOM), args.Get(1).(int64), args.Error(2)
}

func (m *MockSBOMRepository) SearchSBOMs(ctx context.Context, keyword string, offset, limit int) ([]*model.SBOM, int64, error) {
	args := m.Called(ctx, keyword, offset, limit)
	return args.Get(0).([]*model.SBOM), args.Get(1).(int64), args.Error(2)
}

func TestCreateSBOM(t *testing.T) {
	mockContract := new(MockSmartContract)
	mockRepo := new(MockSBOMRepository)
	service := sbom.NewSBOMService(mockContract, mockRepo)

	request := &sbom.CreateSBOMRequest{
		DID:      "test-did",
		SPDXSBOM: &model.SPDXSBOM{
			// Fill with necessary fields
		},
	}

	mockRepo.On("CreateSBOM", mock.Anything, mock.AnythingOfType("*model.SBOM")).Return(nil)
	mockContract.On("StoreSBOM", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)

	sbom, err := service.CreateSBOM(context.Background(), request)
	assert.NoError(t, err)
	assert.NotNil(t, sbom)
	assert.Equal(t, "test-did", sbom.DID)
}

func TestGetSBOM(t *testing.T) {
	mockContract := new(MockSmartContract)
	mockRepo := new(MockSBOMRepository)
	service := sbom.NewSBOMService(mockContract, mockRepo)

	sbomID := uuid.New().String()
	expectedSBOM := &model.SBOM{
		ID:       sbomID,
		DID:      "test-did",
		SPDXSBOM: &model.SPDXSBOM{
			// Fill with necessary fields
		},
	}

	mockRepo.On("GetSBOM", mock.Anything, sbomID).Return(expectedSBOM, nil)

	sbom, err := service.GetSBOM(context.Background(), sbomID)
	assert.NoError(t, err)
	assert.NotNil(t, sbom)
	assert.Equal(t, expectedSBOM.ID, sbom.ID)
	assert.Equal(t, expectedSBOM.DID, sbom.DID)
}
