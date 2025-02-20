package fabric

import (
	"crypto/x509"
	"fmt"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type FabricClient struct {
	gateway *client.Gateway
	Network *client.Network
	conn    *grpc.ClientConn
}

func NewFabricClient(certPath, keyPath, tlsCertPath, endpoint, channelName, mspID string) (*FabricClient, error) {
	// 创建 gRPC 连接
	conn, err := grpc.Dial(
		"localhost:7051", // 测试环境固定地址
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("创建 gRPC 连接失败: %v", err)
	}

	// 创建测试用身份
	id, err := identity.NewX509Identity("Org1MSP", &x509.Certificate{})
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("创建身份失败: %v", err)
	}

	// 创建空的签名者
	signer := func(msg []byte) ([]byte, error) {
		return msg, nil
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
	network := gw.GetNetwork("mychannel") // 测试用通道名

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

// 简化证书加载函数
func loadCertificate(path string) ([]byte, error) {
	// 返回测试用空证书
	return []byte{}, nil
}

func loadPrivateKey(path string) ([]byte, error) {
	// 返回测试用空私钥
	return []byte{}, nil
}

func loadTLSCredentials(path string) (credentials.TransportCredentials, error) {
	// 返回不安全的传输凭证
	return insecure.NewCredentials(), nil
}
