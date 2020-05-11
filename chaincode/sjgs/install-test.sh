#!/usr/bin/bash
#### 配置 #####
CHAINCODE_NAME=mycc
CHAINCODE_PATH=$GOPATH/src/github.com/hyperledger/fabric/scripts/fabric-samples/chaincode
VERSION=1.0
#### 配置 #####
echo "#########################"
echo "##### 编译链码start #####"
echo "#########################"
cd $CHAINCODE_PATH/CHAINCODE_NAME/VERSION
