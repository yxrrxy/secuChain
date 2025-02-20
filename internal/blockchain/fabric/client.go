package fabric

import (
	"crypto/x509"
	"fmt"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type FabricClient struct {
	gateway *client.Gateway
	Network *client.Network
	conn    *grpc.ClientConn
}

func NewFabricClient(certPath, keyPath, tlsCertPath, endpoint, channelName, mspID string) (*FabricClient, error) {
	// 加载身份信息
	clientCert, err := loadCertificate(certPath)
	if err != nil {
		return nil, fmt.Errorf("加载客户端证书失败: %v", err)
	}

	privateKey, err := loadPrivateKey(keyPath)
	if err != nil {
		return nil, fmt.Errorf("加载私钥失败: %v", err)
	}

	// 加载 TLS 证书
	tlsCredentials, err := loadTLSCredentials(tlsCertPath)
	if err != nil {
		return nil, fmt.Errorf("加载 TLS 证书失败: %v", err)
	}

	// 创建 gRPC 连接
	conn, err := grpc.Dial(endpoint,
		grpc.WithTransportCredentials(tlsCredentials),
	)
	if err != nil {
		return nil, fmt.Errorf("创建 gRPC 连接失败: %v", err)
	}

	cert, err := x509.ParseCertificate(clientCert)
	if err != nil {
		return nil, fmt.Errorf("解析证书失败: %v", err)
	}

	signer, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		return nil, fmt.Errorf("创建签名者失败: %v", err)
	}

	id, err := identity.NewX509Identity(mspID, cert)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("创建身份失败: %v", err)
	}

	// 连接到 Gateway
	gw, err := client.Connect(
		id,
		client.WithSign(signer),
		client.WithClientConnection(conn),
	)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("连接 Gateway 失败: %v", err)
	}

	// 获取网络
	network := gw.GetNetwork(channelName)

	return &FabricClient{
		gateway: gw,
		Network: network,
		conn:    conn,
	}, nil
}

func (fc *FabricClient) GetNetwork() *client.Network {
	return fc.Network
}

func (fc *FabricClient) Close() {
	if fc.gateway != nil {
		fc.gateway.Close()
	}
	if fc.conn != nil {
		fc.conn.Close()
	}
}

// 辅助函数
func loadCertificate(path string) ([]byte, error) {
	// 实现证书加载逻辑
	return nil, nil
}

func loadPrivateKey(path string) ([]byte, error) {
	// 实现私钥加载逻辑
	return nil, nil
}

func loadTLSCredentials(path string) (credentials.TransportCredentials, error) {
	// 实现 TLS 证书加载逻辑
	return nil, nil
}
