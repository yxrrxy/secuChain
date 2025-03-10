package did

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/pkg/errors"
	"time"
)

type KeyPairs struct {
	PrivateKey string
	PublicKey  string
	Type       string
}

type Proof struct {
	Creator   string
	Signature string
	Type      string
}

type Authentication struct {
	PublicKey string
	Type      string
}

type Recovery struct {
	PublicKey string
	Type      string
}

type DIDDocument struct {
	did            string
	authentication Authentication
	recovery       Recovery
	proof          Proof
	created        string
	updated        string
}

type DID struct {
	did      string
	authKey  KeyPairs
	recyKey  KeyPairs
	document DIDDocument
}

func generateDID() string {
	specificId := uuid.New().String()
	return "did:" + "fabric:" + specificId
}

// 生成 EC (secp256k1) 密钥对

func generateKeyPair2() (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	curve := elliptic.P256() // Go 标准库没有 secp256k1，可以用 P256，其他库支持 secp256k1
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	return privateKey, &privateKey.PublicKey, nil
}

// 对公钥进行 Base64 编码

func encodePublicKey(pubKey *ecdsa.PublicKey) string {
	pubBytes := elliptic.Marshal(pubKey.Curve, pubKey.X, pubKey.Y)
	return base64.StdEncoding.EncodeToString(pubBytes)
}

// 对私钥进行 Base64 编码

func encodePrivateKey(privKey *ecdsa.PrivateKey) string {
	privBytes := privKey.D.Bytes()
	return base64.StdEncoding.EncodeToString(privBytes)
}

// signDidDocument 使用私钥对 DID Document 进行签名

func signDidDocument(privateKey *ecdsa.PrivateKey, didDocument DIDDocument) (string, error) {
	// 将 DID Document 序列化为 JSON
	documentJson, err := json.Marshal(didDocument)
	if err != nil {
		return "", err
	}

	// 计算 SHA256 哈希值
	hash := sha256.Sum256(documentJson)

	// 使用 ECDSA 对哈希值进行签名
	signature, err := ecdsa.SignASN1(rand.Reader, privateKey, hash[:])
	if err != nil {
		return "", err
	}

	// 返回签名结果的 Base64 编码
	return base64.StdEncoding.EncodeToString(signature), nil
}

// generatorDIDDocument 生成 DID 文档

func generatorDIDDocument(did, publicKey, publicKey2 string) DIDDocument {
	didDocument := DIDDocument{
		did: did,
		authentication: Authentication{
			Type:      "EC",      // 使用 EC 类型
			PublicKey: publicKey, // 填充第一个公钥
		},
		recovery: Recovery{
			Type:      "EC",       // 使用 EC 类型
			PublicKey: publicKey2, // 填充第二个公钥
		},
	}
	return didDocument
}

// createDID 创建 DID

func kkDID() (*DID, error) {
	did := &DID{
		did: generateDID(),
	}

	// 生成主密钥对
	privKey, pubKey, err := generateKeyPair2()
	if err != nil {
		return nil, err
	}

	keyPairs := KeyPairs{
		Type:       "EC",
		PublicKey:  encodePublicKey(pubKey),
		PrivateKey: encodePrivateKey(privKey),
	}

	did.authKey = keyPairs

	// 生成备份密钥对
	privKey2, pubKey2, err := generateKeyPair2()
	if err != nil {
		return nil, err
	}

	keyPairs2 := KeyPairs{
		Type:       "EC",
		PublicKey:  encodePublicKey(pubKey2),
		PrivateKey: encodePrivateKey(privKey2),
	}

	did.recyKey = keyPairs2

	// 生成 DID 文档
	didDocument := generatorDIDDocument(did.did, keyPairs.PublicKey, keyPairs2.PublicKey)
	did.document = didDocument

	// 创建 Proof
	signature, err := signDidDocument(privKey, didDocument)
	if err != nil {
		return nil, err
	}

	proof := Proof{
		Creator:   did.did,
		Type:      "EcdsaSecp256k1Signature",
		Signature: signature,
	}
	did.document.proof = proof
	did.document.created = time.Now().Format(time.RFC3339)
	did.document.updated = did.document.created

	return did, nil
}

func verifyRecoveryKey(publicKey string, privateKey string) bool {
	return publicKey != "" && privateKey != ""
}

// DIDContract 定义了与区块链交互的接口

type DIDContract interface {
	CreateDID(ctx context.Context) (string, error)
	UpdateDID(ctx context.Context, did string, recoveryKey string, recoveryPrivateKey string) (string, error)
	DeleteDID(ctx context.Context, did string, recoveryKey string, recoveryPrivateKey string) error
	GetDID(ctx context.Context, did string) (string, error)
}
type SmartContract struct {
	contract *client.Contract
}

func NewDIDContract(network *client.Network) (*SmartContract, error) {
	contract := network.GetContract("did")
	return &SmartContract{contract: contract}, nil
}

// DID 合约结构体

// registerDID 注册一个新的 DID
func (s *SmartContract) CreateDID(ctx context.Context) (string, error) {
	did, err := kkDID()
	if err != nil {
		return "", fmt.Errorf("創建did失敗: %v", err)
	}
	didJSON, err := json.Marshal(did)
	if err != nil {
		return "", errors.New("序列化 DID 文档失败: " + err.Error())
	}
	// 提交交易到区块链
	_, err = s.contract.SubmitTransaction("CreateDID", did.did, string(didJSON))
	if err != nil {
		return "", fmt.Errorf("提交 DID 注册交易失败: %v", err)
	}
	return string(didJSON), err
}

func (s *SmartContract) UpdateDID(ctx context.Context, did string, recoveryKey string, recoveryPrivateKey string) (string, error) {
	// 读取 DID 数据
	didData, err := s.contract.EvaluateTransaction("QueryDID", did)

	if err != nil {
		return "", errors.New("从世界状态读取 DID 失败: " + err.Error())
	}
	if didData == nil {
		return "", errors.New("DID 不存在")
	}
	// 解析 DID 文档
	var didDoc DID
	err = json.Unmarshal(didData, &didDoc)
	if err != nil {
		return "", errors.New("解码 DID 文档失败: " + err.Error())
	}

	// 验证恢复密钥
	if !verifyRecoveryKey(recoveryKey, recoveryPrivateKey) {
		return "", errors.New("无效的恢复密钥")
	}

	// 验证公钥和私钥匹配
	if didDoc.recyKey.PublicKey != recoveryKey {
		return "", errors.New("恢复密钥不匹配")
	}
	if didDoc.recyKey.PrivateKey != recoveryPrivateKey {
		return "", errors.New("恢复私钥不匹配")
	}

	// 生成新的密钥对
	newprivKey, newpubKey, err := generateKeyPair2()

	newKeyPairs := KeyPairs{
		Type:       "EC",
		PublicKey:  encodePublicKey(newpubKey),
		PrivateKey: encodePrivateKey(newprivKey),
	}

	// 更新 DID 文档
	didDoc.authKey = newKeyPairs
	didDoc.document.updated = time.Now().Format(time.RFC3339)
	didDoc.document.authentication.Type = "EC"
	didDoc.document.authentication.PublicKey = newKeyPairs.PublicKey

	// 将更新后的 DID 文档转换为 JSON
	newJson, err := json.Marshal(didDoc)
	if err != nil {
		return "", fmt.Errorf("failed to marshal updated DID document: %v", err)
	}

	// 更新 DID 文档
	_, err = s.contract.SubmitTransaction("UpdateDID", did, string(newJson))
	if err != nil {
		return "", fmt.Errorf("failed to update DID in world state: %v", err)
	}

	// 返回更新后的 DID 文档 JSON 字符串
	return string(newJson), nil
}

// revoke 注销 DID

func (s *SmartContract) DeleteDID(ctx context.Context, did string, recoveryKey string, recoveryPrivateKey string) error {
	// 读取 DID 数据
	didData, err := s.contract.EvaluateTransaction("GetDID", did)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}
	if didData == nil || len(didData) == 0 {
		return fmt.Errorf("DID does not exist")
	}

	var didDoc DIDDocument
	err = json.Unmarshal(didData, &didDoc)
	if err != nil {
		return fmt.Errorf("failed to unmarshal DID document: %v", err)
	}

	// 验证恢复密钥
	if !verifyRecoveryKey(recoveryKey, recoveryPrivateKey) {
		return fmt.Errorf("invalid recovery key")
	}
	if didDoc.recovery.PublicKey != recoveryKey {
		return fmt.Errorf("recovery key mismatch")
	}

	// 删除 DID
	_, err = s.contract.SubmitTransaction("DeleteDID", did)
	if err != nil {
		return fmt.Errorf("failed to delete state: %v", err)
	}

	return err
}

// queryDID 查询 DID
func (s *SmartContract) GetDID(ctx context.Context, did string) (string, error) {
	// 获取 DID 文档
	existing, err := s.contract.EvaluateTransaction("GetDID", did)
	if err != nil {
		return "", fmt.Errorf("獲取世界狀態失敗: %v", err)
	}
	if existing == nil || len(existing) == 0 {
		return "", fmt.Errorf("DID不存在")
	}
	var DIDs string
	if err := json.Unmarshal(existing, &DIDs); err != nil {
		return "", fmt.Errorf("解析 DID失败: %v", err)
	}

	return DIDs, nil
}

var _ DIDContract = (*SmartContract)(nil)
