package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type AssetMgmt struct{}

type Asset struct {
	Symbol    string `json:"-"`
	Supply    int    `json:"supply"`
	CreatorID string `json:"creator-id"`
}

type Balance struct {
	AccountID string
	Symbol    string
	Value     int
}

type TXRecord struct {
	TXTime    time.Time
	OpType    string
	AccountID string
	Symbol    string
	Amount    int
}

const (
	argEmptyErrorf = "Argument[%s] is empty"
	argNumErrorf   = "Incorrect number of arguments. Expecting %d, actual %d"

	assetKeyPrefix    = "asset_"
	balanceObjectType = "balance:accountid~symbol"
)

func toChaincodeArgs(args ...string) [][]byte {
	bytes := make([][]byte, len(args))
	for i, arg := range args {
		bytes[i] = []byte(arg)
	}
	return bytes
}

func (am *AssetMgmt) createAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error(fmt.Sprintf(argNumErrorf, 2, len(args)))
	}
	symbol := args[0]
	creator := args[1]
	supply, _ := strconv.Atoi(args[2])
	if len(symbol) == 0 {
		return shim.Error(fmt.Sprintf(argEmptyErrorf, "symbol"))
	}
	if len(creator) == 0 {
		return shim.Error(fmt.Sprintf(argEmptyErrorf, "creator"))
	}
	if supply == 0 {
		return shim.Error("Argument[supply] cannot be zero")
	}

	var asset Asset
	bytes, err := stub.GetState(assetKeyPrefix + symbol)
	if err != nil {
		return shim.Error(err.Error())
	}
	if bytes != nil {
		return shim.Error(fmt.Sprintf("Asset[%s] already exists", symbol))
	}

	response := stub.InvokeChaincode("accountmgmt",
		toChaincodeArgs("QueryIDByIDOrName", creator),
		stub.GetChannelID())
	if response.Status != shim.OK {
		return shim.Error(fmt.Sprintf("Failed to invoke chaincode[accountmgmt]. Got error: %s", string(response.Payload)))
	}

	creatorID, err := strconv.Atoi(string(response.Payload))
	if err != nil {
		return shim.Error(err.Error())
	}
	asset.Symbol = symbol
	asset.CreatorID = string(creatorID)
	asset.Supply = supply
	jasset, err := json.Marshal(asset)
	if err != nil {
		return shim.Error(err.Error())
	}
	if err = stub.PutState(assetKeyPrefix+asset.Symbol, jasset); err != nil {
		return shim.Error(err.Error())
	}

	key, err := stub.CreateCompositeKey(balanceObjectType, []string{asset.CreatorID, asset.Symbol})
	if err != nil {
		return shim.Error(err.Error())
	}
	bytes, err = stub.GetState(key)
	if err != nil {
		return shim.Error(err.Error())
	}
	if bytes == nil {
		return shim.Error(fmt.Sprintf("Balance[%s] already exists", key))
	}
	balance := Balance{
		AccountID: asset.CreatorID,
		Symbol:    asset.Symbol,
		Value:     asset.Supply,
	}
	jbalance, err := json.Marshal(balance)
	if err != nil {
		return shim.Error(err.Error())
	}
	if err = stub.PutState(key, jbalance); err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (am *AssetMgmt) transfer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success(nil)
}

// Init is called during Instantiate transaction after the chaincode container
// has been established for the first time, allowing the chaincode to
// initialize its internal data
func (am *AssetMgmt) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke is called to update or query the ledger in a proposal transaction.
// Updated state variables are not committed to the ledger until the
// transaction is committed.
func (am *AssetMgmt) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	f, args := stub.GetFunctionAndParameters()
	switch f {
	case "CreateAsset":
		return am.createAsset(stub, args)
	case "Transfer":
		return am.transfer(stub, args)
	default:
		return shim.Error("Invalid function name")
	}
}

func main() {
	if err := shim.Start(new(AssetMgmt)); err != nil {
		fmt.Printf("Error starting AssetMgmt chaincode: %s", err)
	}
}
