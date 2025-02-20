#!/bin/bash

set -e

# 检查操作系统环境
check_environment() {
    echo "检查操作系统环境..."
    
    # 检查操作系统类型
    case "$(uname -s)" in
        Linux*)     
            export OS_TYPE=Linux
            if grep -q Microsoft /proc/version; then
                export OS_TYPE=WSL
                echo "检测到 WSL 环境"
            else
                echo "检测到 Linux 环境"
            fi
            ;;
        Darwin*)    
            export OS_TYPE=Mac
            echo "检测到 MacOS 环境"
            ;;
        CYGWIN*|MINGW*|MSYS*) 
            export OS_TYPE=Windows
            echo "检测到 Windows 环境"
            ;;
        *)          
            echo "未知操作系统"
            exit 1
            ;;
    esac

    # 检查必要命令
    if ! command -v go &> /dev/null; then
        echo "错误: go 命令未找到"
        echo "请安装 Go: https://golang.org/doc/install"
        exit 1
    fi

    if ! command -v docker &> /dev/null; then
        echo "错误: docker 命令未找到"
        echo "请安装 Docker: https://docs.docker.com/get-docker/"
        exit 1
    fi

    echo "环境检查完成"
}

# 执行环境检查
check_environment

# 获取 go 命令的完整路径
GO_CMD=$(which go)
if [ -z "$GO_CMD" ]; then
    echo "错误: 无法找到 go 命令"
    exit 1
fi

echo "使用 Go: $GO_CMD"
echo "Go 版本: $($GO_CMD version)"

echo "开始设置 Fabric 环境..."

# 环境变量设置
echo "设置环境变量..."
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

cd fabric/test-network

echo "清理现有网络..."
./network.sh down

echo "启动网络并创建通道..."
./network.sh up -ca
./network.sh createChannel -c mychannel

# 部署链码函数
deploy_chaincode() {
    local NAME=$1
    local CHAINCODE_PATH=$2
    echo "部署 $NAME 链码..."
    
    # 使用正确的链码路径
    local REAL_CHAINCODE_PATH="../../../internal/blockchain/contracts/$NAME"
    echo "使用链码路径: $REAL_CHAINCODE_PATH"
    
    # 准备链码
    echo "准备 $NAME 链码..."
    cd $REAL_CHAINCODE_PATH
    
    # 检查 go.mod 文件
    if [ ! -f "go.mod" ]; then
        echo "初始化 go.mod..."
        $GO_CMD mod init blockSBOM/internal/blockchain/contracts/$NAME
    fi
    
    # 添加必要的依赖
    echo "添加依赖..."
    $GO_CMD mod tidy
    $GO_CMD mod vendor
    
    # 打包链码
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
    
    # 移动链码包
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
    
    # 验证部署
    if command -v peer &> /dev/null; then
        peer chaincode list --installed
        peer chaincode list --instantiated -C mychannel
    else
        echo "警告: peer 命令未找到，跳过验证步骤"
    fi
}

# 检查必要的工具
echo "检查必要工具..."
check_tool() {
    if ! command -v $1 &> /dev/null; then
        echo "错误: $1 命令未找到"
        echo "请安装 $1"
        exit 1
    fi
}

# 检查必要的工具
check_tool "jq"
check_tool "tar"
check_tool "zip"

# 部署链码
deploy_chaincode "did" "../internal/blockchain/contracts/did"
deploy_chaincode "sbom" "../internal/blockchain/contracts/sbom"
deploy_chaincode "vuln" "../internal/blockchain/contracts/vuln"

# 注册身份
echo "注册管理员身份..."
cd ../..
go run cmd/fabric/enroll.go

# 健康检查
check_health() {
    echo "检查网络健康状态..."
    if ! docker ps --format '{{.Names}}' | grep -q "peer0.org1.example.com"; then
        echo "错误: Peer 节点未正常启动"
        exit 1
    fi
    docker ps
}

check_health

echo "设置完成！"
echo "DID 和 SBOM 链码已部署"
echo "可以开始运行应用程序了"

# 显示运行的容器
echo "当前运行的 Fabric 容器："
docker ps 