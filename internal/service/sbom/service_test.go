package sbom

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateSBOM(t *testing.T) {
	service := NewSBOMService(mockContract, mockRepo)
	req := &CreateSBOMRequest{
		Name:       "example",
		Version:    "1.0.0",
		Components: []string{"componentA", "componentB"},
		Format:     "spdx",
		DID:        "did:example:123456789",
		Content:    "<SPDX SBOM content here>",
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
