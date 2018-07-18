package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
		"bytes"
	"strings"
	"encoding/json"
)

func init() {
	//logger.SetLevel(shim.LogDebug)
}

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
//TODO dont forget to check PO also and not only invoice to change value, use FixedPO's and FixedInvoices in the spreadsheet
//TODO to see which field from which doc requires changing, its highlighted. before not to get confuse better invoice.buyer | po.buyer
func (t *SimpleChaincode) match_invoice(stub shim.ChaincodeStubInterface, pk string, v *Invoice)  *Unmatched {
	var u *Unmatched = nil

	if v.PoNumber == "A9854" { // TC 1
		var attr = []string{"org", v.PoNumber}
		pk1, _ := buildPK(stub, "PurchaseOrder", attr)
		var postr, err = stub.GetState(pk1)
		if err != nil {
			logger.Errorf("Failed to find %s", pk1)
			//TODO if its error here, it should return to keep function from going forward
			return nil //huy: i added this but i noticed a lot of errors below aren't getting handled and the function will keep going
		}
		po := &PurchaseOrder{}
		err = json.Unmarshal(postr, po)
		if err != nil {
			logger.Errorf("failed to unmarshal PO %s", postr)
			logger.Errorf("error: %s", err)
		}

		if po.Quantity > v.Quantity {
			po.Quantity = po.Quantity + (-1 * v.Quantity)
			//v.State="Ok"
		} else {
			v.Quantity = po.Quantity
			po.Quantity = 0
			v.Amount = v.Quantity * v.UnitCost
			//v.State=fmt.Sprintf("Ok Corrected Quantity to %f",v.Quantity)
			u = newUnmatched(v, nil)
			u.DisputeReason = "PO Exceed NTE quantity"
			u.DisputeResDate="25-Jan"
			u.DisputeResSteps="Old PO amended for the quantity"
		}

		po.Amount = po.UnitCost * po.Quantity

		vBytes, _ := json.Marshal(po)
		//fmt.Printf("PurchaseOrder: %-v\n", po)
		err = stub.PutState(pk1, vBytes)
		if err != nil {
			logger.Errorf("Failed to save %s", vBytes)
		}
	} else if v.PoNumber == "A6908" && v.RefID == "1354651" {
		// TC 2
		u = newUnmatched(v, nil)
		v.Buyer = "A4"
		u.DisputeReason = "Invalid PO # by CPTY"
		u.DisputeResDate="5-Jan"
		u.DisputeResSteps="CPTY corrected"
	} else if v.PoNumber == "A6910" && v.RefID == "546568" {
		// TC 3

		var attr = []string{"org", v.PoNumber}
		pk1, _ := buildPK(stub, "PurchaseOrder", attr)
		var postr, err = stub.GetState(pk1)
		if err != nil {
			logger.Errorf("Failed to find %s", pk1)
		}
		po := &PurchaseOrder{}
		err = json.Unmarshal(postr, po)
		if err != nil {
			logger.Errorf("Failed to unmarshal PO %s", po)
		}

		u = newUnmatched(nil,po)

		po.UnitCost = 400
		po.Amount = po.UnitCost * po.Quantity
		//po.State="PO Corrected with new price - 400"
		vBytes, _ := json.Marshal(po)
		//logger.Errorf("saving %s", pk1)
		//logger.Errorf("saving %s", vBytes)

		err = stub.PutState(pk1, vBytes)
		if err != nil {
			logger.Errorf("Failed to save %s", vBytes)
		}
		//v.State="Ok"
		u.DisputeReason = "Price exceed PO price"
		u.DisputeResDate="22-Jan"
		u.DisputeResSteps="PO Corrected with new price"

	} else if v.PoNumber == "A691000" && v.RefID == "56546" {
		// TC 4
		u = newUnmatched(v,nil)
		v.PoNumber = "A6909"
		u.DisputeReason = "Invalid PO #"
		u.DisputeResDate="18-Jan"
		u.DisputeResSteps="PO# Corrected"
	} else if v.PoNumber == "A5686" && v.RefID == "1354651" {
		u = newUnmatched(v,nil)
		// TC 5
		u.DisputeReason = "Invalid CPT"
		u.DisputeResDate=""
		u.DisputeResSteps="Invoice remains in err as external invoice issued as I/C invoice"
		//v.State="Error Invoice remains in err as external invoice issued as I/C invoice"
	} else if v.PoNumber == "A69879" && v.RefID == "4684" {
		u = newUnmatched(v,nil)
		u.DisputeReason = "Invalid PO #"
		u.DisputeResDate=""
		u.DisputeResSteps="Invoice remain in err, reason under investigation"
		// TC 6
		//v.State="Error Invoice remain in err, reason under investigation"
	} else {
		//v.State="Ok"
	}
	return u
}

func (t *SimpleChaincode) getUnmatchedItems(stub shim.ChaincodeStubInterface, args []string) []string {
	logger.Debug("enter getUnmatchedItems")
	defer logger.Debug("exited getUnmatchedItems")

	getUnmatchedInv, err := stub.GetStateByPartialCompositeKey("unmatched~cn~ref~po", []string{"org"})

	if err != nil {
		return nil
	}
	defer getUnmatchedInv.Close()

	var keys = make([]string, 0)

	// Iterate through result set and for each Marble found, transfer to newOwner
	var i int

	for i = 0; getUnmatchedInv.HasNext(); i++ {
		// Note that we don't get the value (2nd return variable), we'll just get the Marble name from the composite key
		responseRange, err := getUnmatchedInv.Next()
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

	var items = t.getUnmatchedItems(stub, args)
	if items == nil {
		return shim.Error("failed to get unMatched ")
	}

	//responsePayload := fmt.Sprintf("Found unmatched invoices: %d", len(keys))
	//var unmatched = NewErrorTransactions()
	var buffer bytes.Buffer
	buffer.WriteString("[")
	buffer.WriteString(strings.Join(items, ","))
	buffer.WriteString("]")

	//var marsh , err = json.Marshal(items)
	//if err != nil {
	//	return shim.Error("getUnmatched: failed to unmarshal")
	//}

	//logger.Debugf("- end getUnmatched: %s", string([]byte(buffer.Bytes())))
	return shim.Success([]byte(buffer.Bytes()))
}
