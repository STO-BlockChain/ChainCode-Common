/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"log"

	"github.com/hyperledger/fabric-chaincode-go/shim"
)

func main() {
	err := shim.Start(new(Chaincode))
	if err != nil {
		log.Printf("Error starting demo chaincode: %v", err)
	}
}
