package test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"blockSBOM/internal/blockchain/contracts/sbom"
	"blockSBOM/internal/dal/model"
	"blockSBOM/internal/dal/query"
	msbom "blockSBOM/internal/service/sbom"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSmartContract 实现 sbom.SmartContract 接口
type MockSmartContract struct {
	mock.Mock
	contract *sbom.SmartContract
}

func (m *MockSmartContract) StoreSBOM(ctx contractapi.TransactionContextInterface, id string, doc string) error {
	args := m.Called(ctx, id, doc)
	return args.Error(0)
}

func (m *MockSmartContract) GetSBOM(ctx contractapi.TransactionContextInterface, id string) (string, error) {
	args := m.Called(ctx, id)
	return args.String(0), args.Error(1)
}

// MockSBOMRepository 实现 query.SBOMRepository 接口
type MockSBOMRepository struct {
	mock.Mock
	repo *query.SBOMRepository
}

func (m *MockSBOMRepository) CreateSBOM(ctx context.Context, sbom *model.SBOM) error {
	args := m.Called(ctx, sbom)
	return args.Error(0)
}

func (m *MockSBOMRepository) GetSBOM(ctx context.Context, id string) (*model.SBOM, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.SBOM), args.Error(1)
}

// TestCreateSBOM 测试创建 SBOM
func TestCreateSBOM(t *testing.T) {
	// 创建 Mock 实例
	mockContract := &MockSmartContract{}
	mockRepo := &MockSBOMRepository{}

	// 设置 Mock 的预期行为
	mockContract.On("StoreSBOM", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockRepo.On("CreateSBOM", mock.Anything, mock.Anything).Return(nil)

	// 创建 SBOMService 实例
	service := msbom.NewSBOMService(mockContract.contract, mockRepo.repo)
	// 构造测试请求
	req := &msbom.CreateSBOMRequest{
		DID: "did:example:123456789",
		SPDXSBOM: &model.SPDXSBOM{
			SPDXID:      "SPDXRef-1",
			Name:        "example",
			VersionInfo: "1.0.0",
			Supplier:    "Example Supplier",
			ExternalRefs: model.ExternalRefs{
				ReferenceCategory: "package-manager",
				ReferenceLocator:  "https://example.com",
				ReferenceType:     "purl",
			},
		},
	}

	// 调用 CreateSBOM 方法
	sbom, err := service.CreateSBOM(context.Background(), req)

	// 断言
	assert.NoError(t, err)
	assert.NotNil(t, sbom)
	assert.Equal(t, req.DID, sbom.DID)
	assert.Equal(t, req.SPDXSBOM.Name, sbom.SPDXSBOM.Name)
	assert.Equal(t, req.SPDXSBOM.VersionInfo, sbom.SPDXSBOM.VersionInfo)
	assert.Equal(t, req.SPDXSBOM.Supplier, sbom.SPDXSBOM.Supplier)
	assert.Equal(t, req.SPDXSBOM.ExternalRefs.ReferenceCategory, sbom.SPDXSBOM.ExternalRefs.ReferenceCategory)
	assert.Equal(t, req.SPDXSBOM.ExternalRefs.ReferenceLocator, sbom.SPDXSBOM.ExternalRefs.ReferenceLocator)
	assert.Equal(t, req.SPDXSBOM.ExternalRefs.ReferenceType, sbom.SPDXSBOM.ExternalRefs.ReferenceType)

	// 验证 Mock 是否被正确调用
	mockContract.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

// TestGetSBOM 测试获取 SBOM
func TestGetSBOM(t *testing.T) {
	// 创建 Mock 实例
	mockContract := &MockSmartContract{}
	mockRepo := &MockSBOMRepository{}

	// 设置 Mock 的预期行为
	// 第一次调用 GetSBOM 从数据库返回 nil，触发从区块链获取
	mockRepo.On("GetSBOM", mock.Anything, "example-id").Return(nil, fmt.Errorf("not found"))
	// 设置区块链返回的 SBOM 文档
	spdxSBOM := model.SPDXSBOM{
		SPDXID:      "SPDXRef-1",
		Name:        "example",
		VersionInfo: "1.0.0",
		Supplier:    "Example Supplier",
		ExternalRefs: model.ExternalRefs{
			ReferenceCategory: "package-manager",
			ReferenceLocator:  "https://example.com",
			ReferenceType:     "purl",
		},
	}
	sbomData := model.SBOM{
		ID:       "example-id",
		DID:      "did:example:123456789",
		Format:   "spdx",
		SPDXSBOM: &spdxSBOM,
	}
	sbomJSON, _ := json.Marshal(sbomData)
	mockContract.On("GetSBOM", mock.Anything, "example-id").Return(string(sbomJSON), nil)
	// 设置同步到数据库的调用
	mockRepo.On("CreateSBOM", mock.Anything, mock.Anything).Return(nil)

	// 创建 SBOMService 实例
	service := msbom.NewSBOMService(mockContract.contract, mockRepo.repo)

	// 调用 GetSBOM 方法
	sbom, err := service.GetSBOM(context.Background(), "example-id")

	// 断言
	assert.NoError(t, err)
	assert.NotNil(t, sbom)
	assert.Equal(t, "example-id", sbom.ID)
	assert.Equal(t, "did:example:123456789", sbom.DID)
	assert.Equal(t, "example", sbom.SPDXSBOM.Name)
	assert.Equal(t, "1.0.0", sbom.SPDXSBOM.VersionInfo)
	assert.Equal(t, "Example Supplier", sbom.SPDXSBOM.Supplier)
	assert.Equal(t, "package-manager", sbom.SPDXSBOM.ExternalRefs.ReferenceCategory)
	assert.Equal(t, "https://example.com", sbom.SPDXSBOM.ExternalRefs.ReferenceLocator)
	assert.Equal(t, "purl", sbom.SPDXSBOM.ExternalRefs.ReferenceType)

	// 验证 Mock 是否被正确调用
	mockContract.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}
