package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var logger = shim.NewLogger("main logger")

const (
	// INVOICE FC
	ADD_INVOICES    = "AddInvoices"
	GET_INVOICES    = "RetrieveInvoice"
	UPDATE_INVOICES = "UpdateInvoices"

	//POs
	ADD_PO    = "AddPO"
	GET_PO    = "RetrievePO" // ( poNumber )
	UPDATE_PO = "UpdatePO"   //( PurchaseOrder[] )

	//EntityMaster
	ADD_ENTITYMASTER    = "AddEntityMaster"      // ( ImportBlob[] )
	GET_ENTITYMASTER    = "RetrieveEntityMaster" // ( importNumber )
	UPDATE_ENTITYMASTER = "UpdateEntityMaster"   // ( ImportBlob[] )

	// Unmatched
	GET_UNMATCHED = "GetUnmatched" // ()
	DELETE_UNMATCHED = "DeleteUnmatched"

	//Documents

	ADD_DOCUMENTS    = "AddDocuments"      // (args, documentPK)
	GET_DOCUMENTS    = "RetrieveDocuments" // (documentPK)
	UPDATE_DOCUMENTS = "UpdateDocument"    // (args, documentPK)
)

// ===================================================================================
// Main
// ===================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

type SimpleChaincode struct{}

// Init initializes chaincode
// ===========================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("ex02 Init")

	return shim.Success(nil)
}

// Invoke - Our entry point for Invocations
// ========================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	//fmt.Println("invoke is running " + function)

	// Check for health
	//if function == "Ping" {
	//	return t.Ping(stub)
	//}

	// PO operations
	if function == ADD_PO {
		return t.initPurchaseOrders(stub, args)
	} else if function == GET_PO {
		return t.readPurchaseOrder(stub, args)
	} else if function == UPDATE_PO {
		return t.updatePurchaseOrders(stub, args)
	}

	// Entity operations
	if function == ADD_ENTITYMASTER {
		return t.addEntity(stub, args)
	}else if function == GET_ENTITYMASTER {
		return t.getEntity(stub, args)
	}else if function == UPDATE_ENTITYMASTER {
		return t.updateEntities(stub, args)
	}
	// Document operations
	if function == ADD_DOCUMENTS {
		return t.initDocuments(stub, args)
	} else if function == GET_DOCUMENTS {
		return t.readDocument(stub, args)
	} else if function == UPDATE_DOCUMENTS {
		return t.updateDocuments(stub, args)
	}

	// Invoice operations
	if function == ADD_INVOICES {
		return t.addInvoices(stub, args)
	} else if function == GET_INVOICES {
		return t.getInvoice(stub, args)
	} else if function == UPDATE_INVOICES {
		return t.updateInvoices(stub, args)
	}
	// Unmatched Ops
	if function == GET_UNMATCHED {
		return t.getUnmatched(stub, args)
	}else if function == DELETE_UNMATCHED{
		return t.deleteUnmatch()
	}

	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}

//func (this *SimpleChaincode) Ping(stub shim.ChaincodeStubInterface) pb.Response {
//	logger.Info("Ping: enter")
//	defer logger.Info("Ping: exit")
//
//	return shim.Success([]byte("Ok"))
//}
