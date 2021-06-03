package bcconnector

// // Test
// contract := bcconnector.NewConnector()
// log.Println(contract)
// bcconnector.ReadCodeInfo(contract,"TestCID")
// bcconnector.ReadTokenInfo(contract,"TestTID")

// // CreateCodeInfo Test
// mHashV := [32]byte{0}
// mCodeInfo := bcconnector.CodeInfo{
// 	InfoType:        "CodeInfo",
// 	ID_code:         "TestCID8",
// 	DID_RO:          "TestDIDRO2",
// 	DID_client:      "TESTDIDClient2",
// 	Scope:           "all",
// 	Hash_code:       mHashV,
// 	Time_issueed:    "TestTime2",
// 	URI_Redirection: "TESTURI2",
// 	Condition:       "Available",
// 	ID_token:        "TestTID2",
// }
// bcconnector.CreateCodeInfo(contract, mCodeInfo)

// // CreateTokenInfo Test
// mTokenInfo := bcconnector.TokenInfo{
// 	InfoType:        "TokenInfo",
// 	ID_token:        "TestTID8",
// 	DID_RO:          "TestDIDRO2",
// 	DID_client:      "TESTDIDClient2",
// 	Scope:           "all",
// 	Hash_token:      mHashV,
// 	Time_issueed:    "TestTime2",
// 	Time_expiration: "TestDuration",
// 	URI_Redirection: "TESTURI2",
// 	Condition:       "Available",
// }
// bcconnector.CreateTokenInfo(contract, mTokenInfo)

// bcconnector.GetCodeInfoByDID(contract, mCodeInfo.DID_RO)
// bcconnector.GetTokenInfoByDID(contract, mTokenInfo.DID_RO)

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

// CodeInfo describes basic details of what makes up a code informaion
type CodeInfo struct {
	InfoType        string `json:"InfoType"`
	ID_code         string `json:"ID_code"`
	DID_RO          string `json:"DID_RO"`
	DID_client      string `json:"DID_client"`
	Scope           string `json:"Scope"`
	Hash_code       string `json:"Hash_code"`
	Time_issueed    string `json:"Time_issueed"`
	URI_Redirection string `json:"URI_Redirection"`
	Condition       string `json:"Condition"`
	ID_token        string `json:"ID_token"`
}

// TokenInfo describes basic details of what makes up a token informaion
type TokenInfo struct {
	InfoType        string `json:"InfoType"`
	ID_token        string `json:"ID_token"`
	DID_RO          string `json:"DID_RO"`
	DID_client      string `json:"DID_client"`
	Scope           string `json:"Scope"`
	Hash_token      string `json:"Hash_code"`
	Time_issueed    string `json:"Time_issueed"`
	Time_expiration string `json:"Time_expiration"`
	URI_Redirection string `json:"URI_Redirection"`
	Condition       string `json:"Condition"`
}

func NewConnector() *gateway.Contract {
	log.Println("============ application-golang starts ============")

	err := os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	if err != nil {
		log.Fatalf("Error setting DISCOVERY_AS_LOCALHOST environemnt variable: %v", err)
	}

	err = os.RemoveAll("./wallet")
	if err != nil {
		log.Fatalf("Failed to REMOVE wallet directory: %v", err)
	}
	os.RemoveAll("./keystore")
	if err != nil {
		log.Fatalf("Failed to REMOVE keystore directory: %v", err)
	}

	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		log.Fatalf("Failed to create wallet: %v", err)
	}

	if !wallet.Exists("appUser") {
		err = populateWallet(wallet)
		if err != nil {
			log.Fatalf("Failed to populate wallet contents: %v", err)
		}
	}

	ccpPath := filepath.Join(
		"..",
		"test-network",
		"organizations",
		"peerOrganizations",
		"org1.example.com",
		"connection-org1.yaml",
	)

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, "appUser"),
	)
	if err != nil {
		log.Fatalf("Failed to connect to gateway: %v", err)
	}
	defer gw.Close()

	network, err := gw.GetNetwork("mychannel")
	if err != nil {
		log.Fatalf("Failed to get network: %v", err)
	}

	contract := network.GetContract("basic")

	log.Println("[Connector Test]--> Submit Transaction: InitLedger, function creates the initial set of assets on the ledger")
	result, err := contract.SubmitTransaction("InitLedger")
	if err != nil {
		log.Fatalf("Failed to Submit transaction: %v", err)
	}
	log.Println("[Connector Test] : ", string(result))

	log.Println("[Connector Test]--> Evaluate Transaction: ReadCodeInfo, function returns all the current assets on the ledger")
	result, err = contract.EvaluateTransaction("ReadCodeInfo", "TestCID")
	if err != nil {
		log.Fatalf("Failed to evaluate transaction: %v", err)
	}
	log.Println("[Connector Test] : ", string(result))

	return contract
}

func ReadCodeInfo(contract *gateway.Contract, ID_code string) CodeInfo {

	log.Println("--> Evaluate Transaction: ReadCodeInfo, function returns all the current assets on the ledger")
	result, err := contract.EvaluateTransaction("ReadCodeInfo", ID_code)
	if err != nil {
		// log.Fatalf("Failed to evaluate transaction: %v", err)
		log.Println("ReadCodeInfo : Failed to evaluate transaction")
		log.Println(err)

	}
	log.Println("[ReadCodeInfo result]", string(result))
	var iCodeInfo CodeInfo
	json.Unmarshal(result, &iCodeInfo)

	// log.Println("[ReadCodeInfo result]",iCodeInfo.ID_code)

	return iCodeInfo
}

func ReadTokenInfo(contract *gateway.Contract, ID_token string) TokenInfo {

	log.Println("--> Evaluate Transaction: ReadTokenInfo, function returns all the current assets on the ledger")
	result, err := contract.EvaluateTransaction("ReadTokenInfo", ID_token)
	if err != nil {
		log.Fatalf("Failed to evaluate transaction: %v", err)
	}
	log.Println("[ReadTokenInfo result]", string(result))
	log.Println("[--------------------------------------------------------------]", string(result))
	var iTokenInfo TokenInfo
	json.Unmarshal(result, &iTokenInfo)

	return iTokenInfo
}

func CreateCodeInfo(contract *gateway.Contract, mCodeInfo CodeInfo) string {

	log.Println("--> Submit Transaction: CreateCodeInfo, creates new CodeInfo")
	// mHashV := [32]byte{0}
	// result, err = contract.SubmitTransaction("CreateCodeInfo", "TestCID5", "TestDIDRO2", "TESTDIDClient2", "all", string(mHashV[:]), "TestTime2", "TESTURI2", "Available", "TestTID2")
	// result, err = contract.SubmitTransaction("CreateCodeInfo", "TestCID2", "TestDIDRO2", "TESTDIDClient2", "all", "HashTest", "TestTime2", "TESTURI2", "Available", "TestTID2")

	result, err := contract.SubmitTransaction("CreateCodeInfo", mCodeInfo.ID_code, mCodeInfo.DID_RO, mCodeInfo.DID_client, mCodeInfo.Scope, mCodeInfo.Hash_code, mCodeInfo.Time_issueed, mCodeInfo.URI_Redirection, mCodeInfo.Condition, mCodeInfo.ID_token)
	if err != nil {
		log.Fatalf("Failed to Submit transaction: %v", err)
	}
	log.Println("[CreateCodeInfo result]", string(result))

	return string(result)
}

func CreateTokenInfo(contract *gateway.Contract, mTokenInfo TokenInfo) string {

	log.Println("--> Submit Transaction: CreateTokenInfo, creates new CodeInfo")
	result, err := contract.SubmitTransaction("CreateTokenInfo", mTokenInfo.ID_token, mTokenInfo.DID_RO, mTokenInfo.DID_client, mTokenInfo.Scope, mTokenInfo.Hash_token, mTokenInfo.Time_issueed, mTokenInfo.Time_expiration, mTokenInfo.URI_Redirection, mTokenInfo.Condition)
	if err != nil {
		log.Fatalf("Failed to Submit transaction: %v", err)
	}
	log.Println("[CreateTokenInfo result]", string(result))

	return string(result)
}

func GetCodeInfoByDID(contract *gateway.Contract, DID_RO string) string {

	log.Println("--> Evaluate Transaction: GetCodeInfoByDID, function returns all the current assets on the ledger")
	result, err := contract.EvaluateTransaction("GetCodeInfoByDID", DID_RO)
	if err != nil {
		log.Fatalf("Failed to evaluate transaction: %v", err)
	}
	log.Println("[GetCodeInfoByDID result]", string(result))

	return string(result)
}

func GetTokenInfoByDID(contract *gateway.Contract, DID_RO string) string {

	log.Println("--> Evaluate Transaction: GetTokenInfoByDID, function returns all the current assets on the ledger")
	result, err := contract.EvaluateTransaction("GetTokenInfoByDID", DID_RO)
	if err != nil {
		// log.Fatalf("Failed to evaluate transaction: %v", err)
		log.Println(err)
	}
	log.Println("[GetTokenInfoByDID result]", string(result))

	return string(result)
}

func populateWallet(wallet *gateway.Wallet) error {
	log.Println("============ Populating wallet ============")
	credPath := filepath.Join(
		"..",
		"test-network",
		"organizations",
		"peerOrganizations",
		"org1.example.com",
		"users",
		"User1@org1.example.com",
		"msp",
	)

	certPath := filepath.Join(credPath, "signcerts", "cert.pem")
	// read the certificate pem
	cert, err := ioutil.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return err
	}

	keyDir := filepath.Join(credPath, "keystore")
	// there's a single file in this dir containing the private key
	files, err := ioutil.ReadDir(keyDir)
	if err != nil {
		return err
	}
	if len(files) != 1 {
		return fmt.Errorf("keystore folder should have contain one file")
	}
	keyPath := filepath.Join(keyDir, files[0].Name())
	key, err := ioutil.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return err
	}

	identity := gateway.NewX509Identity("Org1MSP", string(cert), string(key))

	return wallet.Put("appUser", identity)
}
