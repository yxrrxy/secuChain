package sbom

import (
	"blockSBOM/internal/dal/model"
	"context"
	"testing"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// 创建 mock 结构体
type MockContract interface {
	StoreSBOM(ctx contractapi.TransactionContextInterface, id string, doc string) error
	GetSBOM(ctx contractapi.TransactionContextInterface, id string) (string, error)
}

type mockContractImpl struct {
	mock.Mock
}

var _ MockContract = (*mockContractImpl)(nil) // 确保实现了接口

func (m *mockContractImpl) StoreSBOM(ctx contractapi.TransactionContextInterface, id string, doc string) error {
	return nil
}

func (m *mockContractImpl) GetSBOM(ctx contractapi.TransactionContextInterface, id string) (string, error) {
	return "", nil
}

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) CreateSBOM(ctx context.Context, sbom *model.SBOM) error {
	return nil
}

func (m *MockRepo) GetSBOM(ctx context.Context, id string) (*model.SBOM, error) {
	return nil, nil
}

var mockContract = &mockContractImpl{}
var mockRepo = &MockRepo{}

func TestCreateSBOM(t *testing.T) {
	service := NewSBOMService(mockContract, mockRepo)
	req := &CreateSBOMRequest{
		DID:      "did:example:123456789",
		SPDXSBOM: &model.SPDXSBOM{Name: "example", Components: []string{"componentA", "componentB"}},
		CDXSBOM:  &model.CDXSBOM{Name: "example", Components: []string{"componentA", "componentB"}},
	}

	sbom, err := service.CreateSBOM(context.Background(), req)
	assert.NoError(t, err)
	assert.NotNil(t, sbom)
	assert.Equal(t, "example", sbom.Name)
}

func TestGetSBOM(t *testing.T) {
	service := NewSBOMService(mockContract, mockRepo)
	id := "example-id"

	sbom, err := service.GetSBOM(context.Background(), id)
	assert.NoError(t, err)
	assert.NotNil(t, sbom)
	assert.Equal(t, "example-id", sbom.ID)
}
