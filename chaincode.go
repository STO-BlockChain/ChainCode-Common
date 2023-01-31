/*
 * SejongTelecom 코어기술개발팀
 * @author JinSan
 */

package main

import (
	"github.com/STO-BlockChain/Common/controller"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

// Chaincode is the definition of the chaincode structure.
type Chaincode struct {
	controller *controller.Controller
}

// Init is called when the chaincode is instantiated by the blockchain network.
func (cc *Chaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {

	return shim.Success(nil)
}

// Invoke is called as a result of an application request to run the chaincode.
func (cc *Chaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	return cc.controller.Controller(stub)
}
