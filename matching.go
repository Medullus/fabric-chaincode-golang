package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
		"bytes"
	"encoding/json"
)

func init() {
	//logger.SetLevel(shim.LogDebug)
}

var unmatchList = make([]Unmatched,0)

type Unmatched struct {
	Amount          float32 `json:"amount,string"`
	Buyer           string `json:"buyer"`
	Currency        string `json:"currency"`
	DisputeReason   string `json:"disputeReason"`
	DisputeResDate  string `json:"disputeResDate"`
	DisputeResSteps string `json:"disputeResSteps"`
	InvStts         string `json:"invStts"`
	PoNum           string `json:"poNum"`
	Quantity        float32 `json:"quantity,string"`
	RefID           string `json:"refId"`
	Seller          string `json:"seller"`
	Sku             string `json:"sku"`
	Unit            float32 `json:"unit,string"`
}
func newUnmatched(i *Invoice, p *PurchaseOrder) *Unmatched {
	var un = &Unmatched{}
	if i != nil {
		un.Amount = i.Amount
		un.RefID = i.RefID
		un.Currency = i.Currency
		un.PoNum = i.PoNumber
		un.Quantity = i.Quantity
		un.Seller = i.Seller
		un.Sku = i.Sku
		un.Unit = i.UnitCost
		return un
	}
	if p != nil {
		un.Amount = p.Amount
		un.RefID = p.RefID
		un.Currency = p.Currency
		un.PoNum = p.RefID
		un.Quantity = p.Quantity
		un.Seller = p.Seller
		un.Sku = p.Sku
		un.Unit = p.UnitCost

		return un

	}
	return nil
}

func createUnMatch(i *Invoice, p *PurchaseOrder)(Unmatched){
	var u Unmatched
	u.Seller = i.Seller
	if p != nil{
		u.Buyer = p.Buyer
	}
	u.RefID = i.RefID
	u.PoNum = i.PoNumber
	u.Sku = i.Sku
	u.Quantity = i.Quantity
	u.Currency = i.Currency
	u.Unit = i.UnitCost
	u.Amount = i.Amount
	u.InvStts = "error"
	logger.Debugf("Create object unmatch:%v",u)
	return u
}
//TODO dont forget to check PO also and not only invoice to change value, use FixedPO's and FixedInvoices in the spreadsheet
//TODO to see which field from which doc requires changing, its highlighted. before not to get confuse better invoice.buyer | po.buyer
func (t *SimpleChaincode) match_invoice(stub shim.ChaincodeStubInterface, pk string, v *Invoice){

	logger.Debugf("Entering matched_invoice with refID:%s",v.RefID)
	defer logger.Debug("Exiting from invoice refID:"+v.RefID)

	if v.RefID == "80203" { // TC 1
		logger.Debug("inside matching for Invoice refid 80203")
		cn, _ := getCN(stub)
		var attr = []string{cn,v.PoNumber}
		pk1,_ := buildPK(stub, "PurchaseOrder", attr)

		poB, err := stub.GetState(pk1)
		if err != nil{
			logger.Error(err)
		}
		po := &PurchaseOrder{}
		err = json.Unmarshal(poB, po)

		if err != nil{
			logger.Error("Failed to unmarshal PO %v",pk1)
		}
		po.Quantity = 100
		po.Amount = 10000
		v.Quantity = 100
		v.Amount = 10000
		unmatch := createUnMatch(v, po)
		unmatch.InvStts = "error"
		unmatch.DisputeReason = "PO Exceed NTE quantity"
		unmatch.DisputeResDate = "25-Jan"
		unmatch.DisputeResSteps = "Old PO amended for the quantity"
		unmatch.Quantity = 180
		unmatch.Amount = 18000

		unmatchList = append(unmatchList, unmatch)
		logger.Debug("finished matching for invoice 80203")

		vBytes, _ := json.Marshal(po)
		//fmt.Printf("PurchaseOrder: %-v\n", po)
		logger.Debug("Matching for Invoice 80203 complete putting po back to state")
		err = stub.PutState(pk1, vBytes)
		if err != nil {
			logger.Errorf("Failed to save %s", vBytes)
		}
	} else if v.PoNumber == "A6908" && v.RefID == "1354651" {
		// TC 2
		logger.Debug("2nd unmatch transaction")
		cn, _ := getCN(stub)
		var attr = []string{cn,v.PoNumber}
		pk1,_ := buildPK(stub, "PurchaseOrder", attr)

		poB, err := stub.GetState(pk1)
		if err != nil{
			logger.Error(err)
		}
		po := &PurchaseOrder{}
		err = json.Unmarshal(poB, po)

		if err != nil{
			logger.Error("Failed to unmarshal PO %v",pk1)
		}

		po.Buyer = "A6"
		unmatch := createUnMatch(v, po)
		unmatch.Buyer = "A4"
		unmatch.DisputeResSteps = "CPTY corrected"
		unmatch.DisputeResDate ="5-Jan"
		unmatch.DisputeReason="Invalid PO # by CPTY"
		unmatchList = append(unmatchList, unmatch)
		logger.Debug("Second transaction unmatch complete")


		vBytes, _ := json.Marshal(po)
		//fmt.Printf("PurchaseOrder: %-v\n", po)
		logger.Debug("Matching for Invoice 80203 complete putting po back to state")
		err = stub.PutState(pk1, vBytes)
		if err != nil {
			logger.Errorf("Failed to save %s", vBytes)
		}

	} else if v.PoNumber == "A6910" && v.RefID == "546568" {
		// TC 3
		logger.Debug("3rd error transaction enter")
		cn, _ := getCN(stub)
		var attr = []string{cn,v.PoNumber}
		pk1,_ := buildPK(stub, "PurchaseOrder", attr)

		poB, err := stub.GetState(pk1)
		if err != nil{
			logger.Error(err)
		}
		po := &PurchaseOrder{}
		err = json.Unmarshal(poB, po)

		if err != nil{
			logger.Error("Failed to unmarshal PO %v",pk1)
		}
		po.UnitCost = 400
		u := createUnMatch(v, po)
		u.DisputeReason = "Price exceed PO price"
		u.DisputeResDate="22-Jan"
		u.DisputeResSteps="PO Corrected with new price"

		unmatchList = append(unmatchList, u)

		vBytes, _ := json.Marshal(po)
		//fmt.Printf("PurchaseOrder: %-v\n", po)
		logger.Debug("Matching for Invoice 80203 complete putting po back to state")
		err = stub.PutState(pk1, vBytes)
		if err != nil {
			logger.Errorf("Failed to save %s", vBytes)
		}

	} else if v.PoNumber == "A691000" && v.RefID == "56546" {
		// TC 4
		logger.Debug("last transaction mismatch")
		v.PoNumber = "A6909"
		cn, _ := getCN(stub)
		var attr = []string{cn,v.PoNumber}
		pk1,_ := buildPK(stub, "PurchaseOrder", attr)

		poB, err := stub.GetState(pk1)
		if err != nil{
			logger.Error(err)
		}
		po := &PurchaseOrder{}
		err = json.Unmarshal(poB, po)

		if err != nil{
			logger.Error("Failed to unmarshal PO %v",pk1)
		}
		u := createUnMatch(v, po)
		u.PoNum = "A691000"
		u.DisputeReason ="Invalid PO #"
		u.DisputeResDate="18-Jan"
		u.DisputeResSteps="PO# Corrected"
		unmatchList = append(unmatchList, u)
		logger.Debug("Finish 4th fixed transaction")



	}else if v.RefID == "1354651" && v.PoNumber == "A5686"{
		u := createUnMatch(v, nil)
		u.Buyer = "A6"
		u.DisputeReason ="Invalid CPT"
		u.DisputeResDate=""
		u.DisputeResSteps="Invoice remains in err as external invoice issued as I/C invoice"
		unmatchList = append(unmatchList, u)

		logger.Debug("5th unmatch complete")
	}else if v.RefID == "4684" && v.PoNumber =="A69879"{
		u := createUnMatch(v, nil)
		u.Buyer = "A5"
		u.DisputeReason ="Invalid PO #"
		u.DisputeResDate=""
		u.DisputeResSteps="Invoice remain in err, reason under investigation"
		unmatchList = append(unmatchList, u)

		logger.Debug("6th unmatch complete")
	}
	logger.Debug("exit matching")
}

func (t *SimpleChaincode) getUnmatchedItems(stub shim.ChaincodeStubInterface, args []string) []string {
	logger.Debug("enter getUnmatchedItems")
	defer logger.Debug("exited getUnmatchedItems")

	getUnmatchedInv, err := stub.GetStateByPartialCompositeKey("unmatched~cn~ref~po", []string{"org"})

	if err != nil {
		return nil
	}
	logger.Debug("found items ")
	defer getUnmatchedInv.Close()

	var keys = make([]string, 0)

	// Iterate through result set and for each Marble found, transfer to newOwner
	var i int

	for i = 0; getUnmatchedInv.HasNext(); i++ {
		logger.Debugf("Unmatched items are:\n\n")
		logger.Debug()
		// Note that we don't get the value (2nd return variable), we'll just get the Marble name from the composite key
		responseRange, err := getUnmatchedInv.Next()
		logger.Debugf("%v:%s",responseRange.Value, responseRange.Key)
		if err != nil {
			return nil
		}

		keys= append(keys, string(responseRange.Value))
		// get the color and name from color~name composite key
		//_, compositeKeyParts, err := stub.SplitCompositeKey(responseRange.Key)
		//if err != nil {
		//	return nil
		//}

		///		returnedUuid := compositeKeyParts[1]
		//var tp = "Invoice"
		//var attr []string
		//if len(compositeKeyParts) == 2 {
		//	tp = "PurchaseOrder"
		//	attr = []string{"org", compositeKeyParts[1]}
		//} else {
		//	attr = []string{"org", compositeKeyParts[1], compositeKeyParts[2]}
		//}
		//
		//pk1, _ := buildPK(stub, tp, attr)
		//
		////logger.Errorf(" %v", pk1)
		//
		//keys = append(keys, pk1)
	}

	//logger.Debugf("- found an unmatched indexes: %s\n", keys)
	return keys
}

// Finds unmatched queries
func (t *SimpleChaincode) getUnmatched(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	logger.Debug("enter get unmatched")
	defer logger.Debug("exited get unmatched")

	if len(unmatchList) < 1{
		var buffer bytes.Buffer
		buffer.WriteString("[")
		buffer.WriteString("]")
		return shim.Success([]byte(buffer.Bytes()))
	}
	var buffer bytes.Buffer
	buffer.WriteString("[")
	for _,v :=range unmatchList{
		b, err := json.Marshal(v)
		if err != nil{
			logger.Error(err)
			shim.Error(err.Error())
		}
		buffer.WriteString(string(b))
		buffer.WriteString(",")
	}
	buffer.Truncate(buffer.Len()-1)
	buffer.WriteString("]")


	logger.Debugf("Sending back:\n%s",buffer)
	return shim.Success([]byte(buffer.Bytes()))

}
