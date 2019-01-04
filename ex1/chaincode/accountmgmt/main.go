package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type AccountMgmt struct{}

type Account struct {
	ID        int    `json:"-"`
	Name      string `json:"name"`
	PublicKey string `json:"public-key"`
}

const (
	argEmptyErrorf = "Argument[%s] is empty"
	argNumErrorf   = "Incorrect number of arguments. Expecting %d, actual %d"

	accountKeyPrefix = "account_"
)

func (am *AccountMgmt) queryIDByIDOrName(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error(fmt.Sprintf(argNumErrorf, 1, len(args)))
	}
	if _, err := strconv.Atoi(args[0]); err == nil {
		bytes, err := stub.GetState(accountKeyPrefix + args[0])
		if err != nil {
			return shim.Error(err.Error())
		}
		if bytes == nil {
			return shim.Success(nil)
		}
		return shim.Success([]byte(args[0]))
	} else {
		str := fmt.Sprintf("{\"selector\": {\"name\": \"%s\"}}", args[0])
		iterator, err := stub.GetQueryResult(str)
		if err != nil {
			return shim.Error(err.Error())
		}
		defer iterator.Close()
		if iterator.HasNext() == false {
			return shim.Success(nil)
		}
		kv, err := iterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success([]byte(strings.TrimPrefix(kv.Key, accountKeyPrefix)))
	}
}

func (am *AccountMgmt) createAccount(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if args
}

func (am *AccountMgmt) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (am *AccountMgmt) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	f, args := stub.GetFunctionAndParameters()
	switch f {
	case "QueryIDByIDOrName":
		return am.queryIDByIDOrName(stub, args)
	case "CreateAccount":
		return am.createAccount(stub, args)
	default:
		return shim.Error("Invalid function name")
	}
}

func main() {
	if err := shim.Start(new(AccountMgmt)); err != nil {
		fmt.Printf("Failed to start chaincode: %s", err.Error())
	}
}
