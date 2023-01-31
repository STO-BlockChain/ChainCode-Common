package controller

import (
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

type Controller struct {
}

func (cc *Controller) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success([]byte("init succeed"))
}

func (cc *Controller) Controller(stub shim.ChaincodeStubInterface) peer.Response {
	fcn, _ := stub.GetFunctionAndParameters()

	switch fcn {
	case "PutDummyData":
		return shim.Error("Method Not Found")
	default:
		return shim.Error("Method Not Found")
	}

}
