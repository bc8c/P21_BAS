package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// CodeInfo describes basic details of what makes up a code informaion
type CodeInfo struct {
	InfoType        string   `json:"InfoType"`
	ID_code         string   `json:"ID_code"`
	DID_RO          string   `json:"DID_RO"`
	DID_client      string   `json:"DID_client"`
	Scope           string   `json:"Scope"`
	Hash_code       string	 `json:"Hash_code"`
	Time_issueed    string   `json:"Time_issueed"`
	URI_Redirection string   `json:"URI_Redirection"`
	Condition       string   `json:"Condition"`
	ID_token        string   `json:"ID_token"`
}

// TokenInfo describes basic details of what makes up a token informaion
type TokenInfo struct {
	InfoType        string   `json:"InfoType"`
	ID_token        string   `json:"ID_token"`
	DID_RO          string   `json:"DID_RO"`
	DID_client      string   `json:"DID_client"`
	Scope           string   `json:"Scope"`
	Hash_token      string 	 `json:"Hash_code"`
	Time_issueed    string   `json:"Time_issueed"`
	Time_expiration string   `json:"Time_expiration"`
	URI_Redirection string   `json:"URI_Redirection"`
	Condition       string   `json:"Condition"`
}

// InitLedger adds a base set of CodeInfo to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {

	mCodeInfo := CodeInfo{InfoType: `CodeInfo`, ID_code: "TestCID", DID_RO: "TestDIDRO", DID_client: "TESTDIDClient", Scope: "all", Hash_code: "TestCHV", Time_issueed: "TestTime", URI_Redirection: "TESTURI", Condition: "Available", ID_token: "TestTID"}

	mTokenInfo := TokenInfo{InfoType: `TokenInfo`, ID_token: "TestTID", DID_RO: "TestDIDRO", DID_client: "TESTDIDClient", Scope: "all", Hash_token: "TestTHV", Time_issueed: "TestTime", Time_expiration: "TestDuration", URI_Redirection: "TESTURI", Condition: "Available"}

	iCodeInfoJSON, err := json.Marshal(mCodeInfo)
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(mCodeInfo.ID_code, iCodeInfoJSON)
	if err != nil {
		return fmt.Errorf("failed to put to world state. %v", err)
	}

	iTokenInfoJSON, err := json.Marshal(mTokenInfo)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(mTokenInfo.ID_token, iTokenInfoJSON)
	if err != nil {
		return fmt.Errorf("failed to put to world state. %v", err)
	}

	return nil
}

func (s *SmartContract) CreateCodeInfo(ctx contractapi.TransactionContextInterface, pID_code string, pDID_RO string, pDID_client string, pScope string, pHash_code string, pTime_issueed string, pURI_Redirection string, pCondition string, pID_token string) error {

	exists, err := s.InfoExists(ctx, pID_code)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the CodeInfo %s already exists", pID_code)
	}

	// mHashV := []byte(pHash_code)
	// mHashV32 := [32]byte{0}
	// copy(mHashV32[:], mHashV)

	mCodeInfo := CodeInfo{
		InfoType:        "CodeInfo",
		ID_code:         pID_code,
		DID_RO:          pDID_RO,
		DID_client:      pDID_client,
		Scope:           pScope,
		Hash_code:       pHash_code,
		Time_issueed:    pTime_issueed,
		URI_Redirection: pURI_Redirection,
		Condition:       pCondition,
		ID_token:        pID_token,
	}
	iCodeInfoJSON, err := json.Marshal(mCodeInfo)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(mCodeInfo.ID_code, iCodeInfoJSON)
}

func (s *SmartContract) CreateTokenInfo(ctx contractapi.TransactionContextInterface, pID_token string, pDID_RO string, pDID_client string, pScope string, pHash_token string, pTime_issueed string, pTime_expiration string, pURI_Redirection string, pCondition string) error {

	exists, err := s.InfoExists(ctx, pID_token)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the TokenInfo %s already exists", pID_token)
	}

	// mHashV := []byte(pHash_token)
	// mHashV32 := [32]byte{0}
	// copy(mHashV32[:], mHashV)

	mTokenInfo := TokenInfo{
		InfoType:        "TokenInfo",
		ID_token:        pID_token,
		DID_RO:          pDID_RO,
		DID_client:      pDID_client,
		Scope:           pScope,
		Hash_token:      pHash_token,
		Time_issueed:    pTime_issueed,
		Time_expiration: pTime_expiration,
		URI_Redirection: pURI_Redirection,
		Condition:       pCondition,
	}
	iTokenInfoJSON, err := json.Marshal(mTokenInfo)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(mTokenInfo.ID_token, iTokenInfoJSON)
}

// ReadCodeInfo returns the CodeInfo stored in the world state with given ID_code.
func (s *SmartContract) ReadCodeInfo(ctx contractapi.TransactionContextInterface, pID_code string) (*CodeInfo, error) {
	iCodeInfoJSON, err := ctx.GetStub().GetState(pID_code)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if iCodeInfoJSON == nil {
		return nil, fmt.Errorf("the CodeInfo %s does not exist", pID_code)
	}

	var mCodeInfo CodeInfo
	err = json.Unmarshal(iCodeInfoJSON, &mCodeInfo)
	if err != nil {
		return nil, err
	}

	return &mCodeInfo, nil
}

// ReadTokenInfo returns the TokenInfo stored in the world state with given ID_token.
func (s *SmartContract) ReadTokenInfo(ctx contractapi.TransactionContextInterface, pID_token string) (*TokenInfo, error) {
	iTokenInfoJSON, err := ctx.GetStub().GetState(pID_token)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if iTokenInfoJSON == nil {
		return nil, fmt.Errorf("the CodeInfo %s does not exist", pID_token)
	}

	var mTokenInfo TokenInfo
	err = json.Unmarshal(iTokenInfoJSON, &mTokenInfo)
	if err != nil {
		return nil, err
	}

	return &mTokenInfo, nil
}

// UpdateCodeInfo updates an existing asset in the world state with provided parameters.
func (s *SmartContract) UpdateCodeInfo(ctx contractapi.TransactionContextInterface, pID_code string, pDID_RO string, pDID_client string, pScope string, pHash_code string, pTime_issueed string, pURI_Redirection string, pCondition, string, pID_token string) error {

	exists, err := s.InfoExists(ctx, pID_code)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the CodeInfo %s does not exist", pID_code)
	}

	mCodeInfo := CodeInfo{
		InfoType:        "CodeInfo",
		ID_code:         pID_code,
		DID_RO:          pDID_RO,
		DID_client:      pDID_client,
		Scope:           pScope,
		Hash_code:       pHash_code,
		Time_issueed:    pTime_issueed,
		URI_Redirection: pURI_Redirection,
		Condition:       pCondition,
		ID_token:        pID_token,
	}
	iCodeInfoJSON, err := json.Marshal(mCodeInfo)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(mCodeInfo.ID_code, iCodeInfoJSON)
}

func (s *SmartContract) GetCodeInfoByDID(ctx contractapi.TransactionContextInterface, pDID_RO string) ([]*CodeInfo, error) {

	queryString := fmt.Sprintf("{\"selector\":{\"InfoType\":\"CodeInfo\",\"DID_RO\":\"%s\"}}", pDID_RO)
	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var mCodeInfo []*CodeInfo
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var iCodeInfo CodeInfo
		err = json.Unmarshal(queryResult.Value, &iCodeInfo)
		if err != nil {
			return nil, err
		}
		mCodeInfo = append(mCodeInfo, &iCodeInfo)
	}
	return mCodeInfo, nil
}

func (s *SmartContract) GetTokenInfoByDID(ctx contractapi.TransactionContextInterface, pDID_RO string) ([]*TokenInfo, error) {

	queryString := fmt.Sprintf("{\"selector\":{\"InfoType\":\"TokenInfo\",\"DID_RO\":\"%s\"}}", pDID_RO)
	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var mTokenInfo []*TokenInfo
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var iTokenInfo TokenInfo
		err = json.Unmarshal(queryResult.Value, &iTokenInfo)
		if err != nil {
			return nil, err
		}
		mTokenInfo = append(mTokenInfo, &iTokenInfo)
	}
	return mTokenInfo, nil
}

// InfoExists returns true when Info with given ID exists in world state
func (s *SmartContract) InfoExists(ctx contractapi.TransactionContextInterface, pid string) (bool, error) {
	iInfoJSON, err := ctx.GetStub().GetState(pid)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	return iInfoJSON != nil, nil
}
