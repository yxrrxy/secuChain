#!/bin/bash

# 设置环境变量
export PATH=${PWD}/fabric/bin
export FABRIC_CFG_PATH=${PWD}/fabric/config

# 测试 peer 命令
echo "测试 peer 命令..."

# 检查 peer 命令是否存在
if ! command -v peer &> /dev/null; then
    echo "错误: peer 命令未找到"
    echo "PATH: $PATH"
    echo "当前目录: $(pwd)"
    exit 1
fi

# 显示 peer 版本
echo "Peer 版本信息:"
peer version

# 测试链码列表
echo "测试链码列表..."
cd fabric/test-network

# 设置必要的环境变量
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_ADDRESS=localhost:7051
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp

# 查询已安装的链码
echo "已安装的链码:"
peer lifecycle chaincode queryinstalled

# 查询已提交的链码
echo "已提交的链码:"
peer lifecycle chaincode querycommitted -C mychannel

echo "测试完成" 