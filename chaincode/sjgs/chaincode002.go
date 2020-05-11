package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"strconv"
)

type SmartContract struct {
}

type Produce struct {
	Temp     string `json:"temp"`
	Humidity string `json:"humidity"`
}

/**
初始化
*/
func (s *SmartContract) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("chain code init")
	return shim.Success(nil)
}

/**
调用
*/
func (s *SmartContract) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	// api 提取调用的函数和参数，第一个是参数，后面参数数组
	fn, args := stub.GetFunctionAndParameters()
	// 参数校验
	if fn == "save" {
		return s.save(stub, args)
	} else if fn == "query" {
		return s.query(stub, args)
	} else if fn == "del" {
		return s.del(stub, args)
	} else if fn == "getHistoryForKey" {
		return s.getHistoryForKey(stub, args)
	}
	return shim.Error("No func")
}

/**
保存 key
*/
func (s *SmartContract) save(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		return shim.Error("args error")
	}
	var produce = Produce{Temp: args[1], Humidity: args[2]}
	byteData, _ := json.Marshal(produce)
	// 在账本中插入或更新键值对
	err := stub.PutState(args[0], byteData)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success([]byte(stub.GetTxID()))
}

/**
查询单个key当前数据
*/
func (s *SmartContract) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("args error")
	}
	byteData, _ := stub.GetState(args[0])
	return shim.Success(byteData)
}

/**
删除单个key
*/
func (s *SmartContract) del(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("args error")
	}
	err := stub.DelState(args[0])
	if err != nil {
		return shim.Error("del error")
	}
	return shim.Success([]byte(stub.GetTxID()))
}

/**
查询单个key的历史
*/
func (s *SmartContract) getHistoryForKey(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("args error")
	}
	resultsIterator, err := stub.GetHistoryForKey(args[0])
	if err != nil {
		return shim.Error("getHistoryForKey error")
	}
	defer resultsIterator.Close()

	// buffer 存储数据，组装json
	var buffer bytes.Buffer

	buffer.WriteString("[")

	isWrite := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		if isWrite == true {
			buffer.WriteString(",")
		}

		buffer.WriteString("{")
		buffer.WriteString("\"txid\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.GetTxId())
		buffer.WriteString("\"")

		buffer.WriteString(",\"timestamp\": ")
		buffer.WriteString("\"")
		//buffer.WriteString(time.Unix(queryResponse.Timestamp.Seconds, int64(queryResponse.Timestamp.Nanos)).String())
		timestamp := strconv.FormatInt(queryResponse.Timestamp.Seconds, 10)
		buffer.WriteString(timestamp)
		buffer.WriteString("\"")

		buffer.WriteString(",\"value\": ")
		// 空数据判断，删除key时候数据就会变空
		if queryResponse.IsDelete {
			buffer.WriteString("\"\"")
		} else {
			buffer.WriteString(string(queryResponse.Value))
		}

		buffer.WriteString(",\"isDelete\": ")
		buffer.WriteString(strconv.FormatBool(queryResponse.IsDelete))

		buffer.WriteString("}")
		isWrite = true
	}

	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())
}

func main() {
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Println("start error")
	}
}
