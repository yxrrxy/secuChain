#!/bin/bash

set -e

GO_CMD=$(which go)
if [ -z "$GO_CMD" ]; then
    echo "错误: 无法找到 go 命令"
    echo "请安装 Go: https://golang.org/doc/install"
    exit 1
fi

export PATH=${PWD}/fabric/bin:/usr/local/bin:/usr/bin:$PATH
export FABRIC_CFG_PATH=${PWD}/fabric/config
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_ADDRESS=localhost:7051
export GOPROXY=https://goproxy.cn,direct

# 设置组织信息
ORG1_PEER_TLS_ROOTCERT=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
ORG1_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_TLS_ROOTCERT_FILE=$ORG1_PEER_TLS_ROOTCERT
export CORE_PEER_MSPCONFIGPATH=$ORG1_MSPCONFIGPATH

cd scripts/fabric/test-network


find . -type f -name "*.sh" -exec sed -i 's/\r$//' {} +
find . -type f -name "*.config" -exec sed -i 's/\r$//' {} +

sed -i 's/\r$//' ./organizations/fabric-ca/*.sh 2>/dev/null || true
sed -i 's/\r$//' ./organizations/**/*.sh 2>/dev/null || true

find . -type f -name "*.sh" -exec chmod +x {} +

chmod +x ./network.sh
chmod +x ./scripts/*.sh 2>/dev/null || true

bash ./network.sh down
bash ./network.sh up -ca
bash ./network.sh createChannel -c mychannel

check_tool() {
    if ! command -v $1 &> /dev/null; then
        echo "错误: $1 命令未找到"
        echo "请安装 $1"
        exit 1
    fi
}
check_tool "jq"
check_tool "tar"
check_tool "zip"

deploy_chaincode "did" "../internal/blockchain/contracts/did"
deploy_chaincode "sbom" "../internal/blockchain/contracts/sbom"
deploy_chaincode "vuln" "../internal/blockchain/contracts/vuln"

check_health
docker ps 

# 注册身份
#echo "注册管理员身份..."
#cd ../..
#go run cmd/fabric/enroll.go

echo "已安装的链码列表："
peer chaincode list --installed || echo "无法获取已安装链码列表"

check_health() {
    if ! docker ps --format '{{.Names}}' | grep -q "peer0.org1.example.com"; then
        echo "错误: Peer 节点未正常启动"
        exit 1
    fi
    docker ps
}

deploy_chaincode() {
    local NAME=$1
    local CHAINCODE_PATH=$2
    echo "部署 $NAME 链码..."
    
    local REAL_CHAINCODE_PATH="../../internal/blockchain/contracts/$NAME"
    
    echo "准备 $NAME 链码..."
    cd $REAL_CHAINCODE_PATH
    
    if [ ! -f "go.mod" ]; then
        echo "初始化 go.mod..."
        $GO_CMD mod init blockSBOM/internal/blockchain/contracts/$NAME
    fi
    
    echo "添加依赖..."
    $GO_CMD mod tidy
    $GO_CMD mod vendor
    
    echo "打包链码..."
    if /usr/bin/tar czf ${NAME}.tar.gz ./* 2>/dev/null; then
        echo "tar 打包成功"
    else
        echo "tar 打包失败，尝试使用 zip..."
        if /usr/bin/zip -r ${NAME}.zip ./* 2>/dev/null; then
            /bin/mv ${NAME}.zip ${NAME}.tar.gz
            echo "zip 打包成功"
        else
            echo "打包失败"
            exit 1
        fi
    fi
    
    echo "移动链码包..."
    /bin/mv ${NAME}.tar.gz ../../../../scripts/fabric/test-network/
    
    # 返回测试网络目录并部署
    cd ../../../../scripts/fabric/test-network/
    
    # 部署链码
    echo "部署链码..."
    bash ./network.sh deployCC \
        -ccn $NAME \
        -ccp $REAL_CHAINCODE_PATH \
        -ccl go \
        -ccv 1.0 \
        -ccs 1 \
        -verbose
    
    if command -v peer &> /dev/null; then
        peer chaincode list --installed
        peer chaincode list --instantiated -C mychannel
    else
        echo "警告: peer 命令未找到，跳过验证步骤"
    fi
}
