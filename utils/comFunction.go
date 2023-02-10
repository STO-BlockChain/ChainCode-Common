package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"strconv"
	"strings"

	model "github.com/STO-BlockChain/ChainCode-Common/model"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

// DoTransfer is 토큰 Transfer
func DoTransfer(stub shim.ChaincodeStubInterface, tokenName string, transfMeta model.TransferMetaN) peer.Response {

	chainCodeFunc := "transfer"

	invokeArgs := ToChaincodeArgs(chainCodeFunc, transfMeta.FromAddress, transfMeta.ToAddress, strconv.Itoa(int(transfMeta.Amount)))
	channel := stub.GetChannelID()
	response := stub.InvokeChaincode(tokenName, invokeArgs, channel)

	return response
}

// DoBalanceOf is 토큰 balanceOf
func DoBalanceOf(stub shim.ChaincodeStubInterface, toaddress string, tokenName string) peer.Response {

	// 지갑형 트랜잭션 VAILD WALLET CHECK 및 지갑주소/파라미터 파싱
	chainCodeFunc := "balanceOf"
	invokeArgs := ToChaincodeArgs(chainCodeFunc, toaddress)
	channel := stub.GetChannelID()
	response := stub.InvokeChaincode(tokenName, invokeArgs, channel)

	if response.Status != shim.OK {
		errStr := fmt.Sprintf("Failed to balanceOf chaincode. Got error: %s", string(response.Payload))
		fmt.Printf(errStr)
		return peer.Response{Status: 501, Message: "balanceOf Fail!", Payload: nil}
	}

	return response
}

// DoTokenFunc is 토큰 함수 실행 (burn, mint)
func DoTokenFunc(stub shim.ChaincodeStubInterface, funcName string, transParam string, tokenName string) peer.Response {
	chainCodeFunc := funcName
	invokeArgs := ToChaincodeArgs(chainCodeFunc, transParam)
	channel := stub.GetChannelID()
	response := stub.InvokeChaincode(tokenName, invokeArgs, channel)
	return response
}

// DoTransferMulti is 토큰 TransferMulti
func DoTransferMulti(stub shim.ChaincodeStubInterface, callerAddress string, stTransferMetaArr []model.TransferMeta, tokenName string) peer.Response {

	chainCodeFunc := "transferMulti"
	stTransferStr, _ := json.Marshal(stTransferMetaArr)
	invokeArgs := ToChaincodeArgs(chainCodeFunc, callerAddress, string(stTransferStr))
	channel := stub.GetChannelID()
	response := stub.InvokeChaincode(tokenName, invokeArgs, channel)

	if response.Status != shim.OK {
		errStr := fmt.Sprintf("Failed to transfer chaincode. Got error: %s", string(response.Payload))
		fmt.Printf(errStr)
		return peer.Response{Status: 501, Message: "transfer Fail!", Payload: nil}
	}

	return response
}

func DoTransferMultiNoneSafetyN(stub shim.ChaincodeStubInterface, stTransferMetaArr []model.TransferMetaN, tokenName string) peer.Response {

	chainCodeFunc := "transferMultiNoneSafetyN"
	stTransferStr, _ := json.Marshal(stTransferMetaArr)
	invokeArgs := ToChaincodeArgs(chainCodeFunc, string(stTransferStr))
	channel := stub.GetChannelID()
	response := stub.InvokeChaincode(tokenName, invokeArgs, channel)

	if response.Status != shim.OK {
		errStr := fmt.Sprintf("Failed to transfer chaincode. Got error: %s", string(response.Payload))
		fmt.Printf(errStr)
		return peer.Response{Status: 501, Message: "transfer Fail!", Payload: nil}
	}

	return response
}

// ToChaincodeArgs is 외부 체인코드 호출시 파라미터 만드는 함수
func ToChaincodeArgs(args ...string) [][]byte {
	bargs := make([][]byte, len(args))
	for i, arg := range args {
		bargs[i] = []byte(arg)
	}
	return bargs
}

// GetNowDt is 현재 시간 반환
func GetNowDt(stub shim.ChaincodeStubInterface) int64 {
	nowTimestamp, _ := stub.GetTxTimestamp()
	nowdt := nowTimestamp.GetSeconds()

	return nowdt
}

// JsonFromQueryResponse 은 iterator 를 json 으로 변환
// query result iterator 와 응답 메타데이터를 넘기면, 리턴할 json 으로 변환해 줌
func JsonFromQueryResponse(resultsIterator shim.StateQueryIteratorInterface, responseMetadata *peer.QueryResponseMetadata) (*bytes.Buffer, error) {
	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	//buffer.WriteString("[")
	buffer.WriteString("{\"bookmark\":")
	buffer.WriteString("\"")
	buffer.WriteString(responseMetadata.Bookmark)
	buffer.WriteString("\"")
	buffer.WriteString(",")
	buffer.WriteString("\"recordcnt\":")
	buffer.WriteString("\"")
	buffer.WriteString(fmt.Sprintf("%v", responseMetadata.FetchedRecordsCount))
	buffer.WriteString("\"")
	buffer.WriteString(", \"reqlist\":")
	buffer.WriteString("[")
	var i int64 = 0
	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// 두번째 array 부터는 , 붙이기
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		// json object 를 string 으로 변환해서
		buffer.WriteString(string(queryResponse.Value))
		bArrayMemberAlreadyWritten = true
		i++
	}
	buffer.WriteString("]")
	buffer.WriteString("}")
	return &buffer, nil
}

// SaveMetaData is 데이터 저장
func SaveMetaData(stub shim.ChaincodeStubInterface, dataKey string, metaDataBytes []byte) peer.Response {

	// 저장
	err := stub.PutState(dataKey, metaDataBytes)
	if err != nil {
		return shim.Error("failed to PutState, error: " + err.Error())
	}
	return shim.Success(nil)
}

// CreateCompositKeyAndPut is
// params : stub, keytype(composite objType), keycode(attribute), data(save data(doc))
// return : compositeKey(string), err
func CreateCompositKeyAndPut(stub shim.ChaincodeStubInterface, keytype string, keycode []string, data []byte) (*string, error) {
	// create composite key
	compositeKey, err := stub.CreateCompositeKey(keytype, keycode)
	if err != nil {
		return nil, model.NewCustomError(model.CreateCompositeKeyErrorType, "compositeKey", err.Error())
	}
	// save data
	err = stub.PutState(compositeKey, data)
	if err != nil {
		return nil, model.NewCustomError(model.PutStateErrorType, compositeKey, err.Error())
	}
	return &compositeKey, nil
}

//ConvertStringToUint64 is ..
func ConvertStringToUint64(typeName, value string) (*uint64, error) {

	// check amount is integer & positive
	intValue, err := strconv.Atoi(value)

	if err != nil {
		return nil, model.NewCustomError(model.ConvertErrorType, typeName, " must be integer")
	}

	if intValue < 0 {
		return nil, model.NewCustomError(model.ConvertErrorType, typeName, " must be positive")
	}
	uint64Value := uint64(intValue)
	return &uint64Value, nil
}

func GetFundAdmin(stub shim.ChaincodeStubInterface, fundid string) string {

	chainCodeFunc := "fundAddressSearch"
	invokeArgs := ToChaincodeArgs(chainCodeFunc, fundid)
	channel := stub.GetChannelID()
	response := stub.InvokeChaincode("fund", invokeArgs, channel)

	address := string(response.Payload)

	fmt.Println("FUND ADMIN:", fundid, address)

	return address
}

func GetCallerInfo() (funcName string, line int) {
	_, fn, line, ok := runtime.Caller(2)
	if ok {
		if index := strings.LastIndex(fn, "/"); index >= 0 {
			fn = fn[index+1:]
		}
	} else {
		fn = "???"
		line = 1
	}

	return fn, line
}

func InfoLogger(infoLog string) {
	funcName, line := GetCallerInfo()
	resultLog := fmt.Sprintf("%s:%d: %s", funcName, line, infoLog)
	log.SetPrefix("[INFO] ")
	log.Println(resultLog)
}

func ErrorLogger(errorLog string) {
	funcName, line := GetCallerInfo()
	resultLog := fmt.Sprintf("%s:%d: %s", funcName, line, errorLog)
	log.SetPrefix("[ERROR] ")
	log.Println(resultLog)
}
