package did

import (
	"blockSBOM/internals/blockchain/contracts/did"
	"blockSBOM/internals/dal/model"
	"blockSBOM/internals/dal/query"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
)

type DIDService struct {
	contract did.DIDContract
	repo     *query.DIDRepository
}

func NewDIDService(contract did.DIDContract, repo *query.DIDRepository) *DIDService {
	return &DIDService{
		contract: contract,
		repo:     repo,
	}
}

// RegisterDID 注册DID
func (s *DIDService) RegisterDID(c context.Context) (string, error) {
	var didData string
	// 创建DID
	didByte, err := s.contract.CreateDID(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed to register DID: %v", err)
	}

	// 定义一个空的 DID 结构体
	var did model.DID

	// 将返回的 JSON 数据解析为 DID 结构体
	if err := json.Unmarshal([]byte(didByte), &did); err != nil {
		return "", fmt.Errorf("failed to unmarshal DID data: %v", err)
	}

	// 再写入数据库
	if err := s.repo.CreateDID(c, &did); err != nil {
		return "", fmt.Errorf("创建数据库 DID 失败: %v", err)
	}
	return didData, nil
}

// UpdateDID 更新 DID
func (s *DIDService) UpdateDID(ctx context.Context, did string, recoveryKey string, recoveryPrivateKey string) (string, error) {
	// 调用智能合约更新 DID
	resultBytes, err := s.contract.UpdateDID(context.Background(), did, recoveryKey, recoveryPrivateKey)
	if err != nil {
		return "", fmt.Errorf("更新 DID 失败: %v", err)
	}
	var result string
	result = string(resultBytes)
	var did2 model.DID

	// 将返回的 JSON 数据解析为 DID 结构体
	if err := json.Unmarshal([]byte(resultBytes), &did2); err != nil {
		return "", fmt.Errorf("failed to unmarshal DID data: %v", err)
	}

	// 再更新数据库
	if err := s.repo.UpdateDID(ctx, &did2); err != nil {
		return "", fmt.Errorf("更新数据库 DID 失败: %v", err)
	}

	// 返回更新结果
	return result, nil
}

// DeleteDID 删除 DID
func (s *DIDService) DeleteDID(c context.Context, did string, recoveryKey string, recoveryPrivateKey string) (string, error) {
	// 调用智能合约撤销 DID
	err := s.contract.DeleteDID(context.Background(), did, recoveryKey, recoveryPrivateKey)
	if err != nil {
		return "", fmt.Errorf("撤销 DID 失败: %v", err)
	}
	// 同步到数据库
	if err := s.repo.DeleteDID(c, did); err != nil {
		// 仅记录错误，不影响返回结果
		fmt.Printf("同步 DID 到数据库失败: %v\n", err)
	}
	// 返回操作成功信息
	return "注销成功", nil
}

func (s *DIDService) ResolveDID(ctx context.Context, did string) (*model.DID, error) {
	// 优先从数据库查询
	doc, err := s.repo.GetDID(ctx, did)
	if err == nil {
		return doc, nil
	}
	// 数据库查询失败，从区块链获取
	docString, err := s.contract.GetDID(context.Background(), did)
	var did2 model.DID
	// 将返回的 JSON 数据解析为 DID 结构体
	if err := json.Unmarshal([]byte(docString), &did2); err != nil {
		return nil, fmt.Errorf("failed to unmarshal DID data: %v", err)
	}
	if err != nil {
		return &did2, errors.New("解析区块链 DID 失败")
	}
	// 同步到数据库
	if err := s.repo.CreateDID(ctx, &did2); err != nil {
		// 仅记录错误，不影响返回结果
		fmt.Printf("同步 DID 到数据库失败: %v\n", err)
	}

	return &did2, nil
}
