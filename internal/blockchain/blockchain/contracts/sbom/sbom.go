package contracts

import (
	"blockSBOM/backend/internal/blockchain/fabric"
	"blockSBOM/backend/internal/dal/model"
	"encoding/json"
	"fmt"
)

type SBOMContract struct {
	client *fabric.FabricClient
}

func NewSBOMContract(client *fabric.FabricClient) *SBOMContract {
	return &SBOMContract{client: client}
}

func (c *SBOMContract) StoreSBOM(sbom *model.SBOM) error {
	sbomBytes, err := json.Marshal(sbom)
	if err != nil {
		return fmt.Errorf("序列化SBOM失败: %v", err)
	}

	_, err = c.client.Contract.SubmitTransaction("StoreSBOM", sbom.ID, string(sbomBytes))
	if err != nil {
		return fmt.Errorf("存储SBOM失败: %v", err)
	}

	return nil
}

func (c *SBOMContract) GetSBOM(id string) (*model.SBOM, error) {
	result, err := c.client.Contract.EvaluateTransaction("GetSBOM", id)
	if err != nil {
		return nil, fmt.Errorf("获取SBOM失败: %v", err)
	}

	var sbom model.SBOM
	if err := json.Unmarshal(result, &sbom); err != nil {
		return nil, fmt.Errorf("反序列化SBOM失败: %v", err)
	}

	return &sbom, nil
}
