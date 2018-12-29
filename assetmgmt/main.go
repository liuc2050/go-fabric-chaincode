package main

import (
	"fmt"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type AssetMgmt struct{}

type Asset struct {
	Symbol string
	Supply int
}

type Account struct {
	ID        int
	Name      string
	PublicKey string
}

type Balance struct {
	AccountID int
	Symbol    string
	Value     int
}

type TXRecord struct {
	TXTime    time.Time
	OpType    string
	AccountID int
	Symbol    string
	Amount    int
}

type cmd struct {
	name string
	fn   func()
	args []string
}

func (am *AssetMgmt) createAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// if len(args) < {
	// 	return shim.Error(fmt.Sprintf("Incorrect number of arguments. Expecting "))
	// }
	//stub.GetState()
	return shim.Success(nil)
}

// Init is called during Instantiate transaction after the chaincode container
// has been established for the first time, allowing the chaincode to
// initialize its internal data
func (am *AssetMgmt) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Printf("GetStringArgs:%v\n", stub.GetStringArgs())
	fn, args := stub.GetFunctionAndParameters()
	fmt.Printf("GetFunctionAndParamters: fn=%s, args=%v\n", fn, args)

	return shim.Success(nil)
}

// Invoke is called to update or query the ledger in a proposal transaction.
// Updated state variables are not committed to the ledger until the
// transaction is committed.
func (am *AssetMgmt) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	f, args := stub.GetFunctionAndParameters()
	switch f {
	case "CreateAsset":
		return
	}
	return shim.Success(nil)
}

func main() {
	if err := shim.Start(new(AssetMgmt)); err != nil {
		fmt.Printf("Error starting AssetMgmt chaincode: %s", err)
	}
}
