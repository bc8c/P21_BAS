package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

//웹페이지
var welcome_page string = "static/index.html"
var ctrct *gateway.Contract

const UserDID = "did:BAS:u12345678abcdefghi"

func main() {

	ctrct = NewConnector()

	//파일시스템 포인팅 : css, js
	fs := http.StripPrefix("/static/", http.FileServer(http.Dir("static")))
	http.Handle("/static/", fs)

	//라우터
	http.HandleFunc("/", mainPage)
	http.HandleFunc("/code", codeFunc)
	http.HandleFunc("/token", tokenFunc)

	//서버가동
	log.Println("서버 시작")
	http.ListenAndServe(":5000", nil)
}

////////////////////////////////////////////////////////////////
func mainPage(res http.ResponseWriter, req *http.Request) {
	webpage, _ := ioutil.ReadFile(welcome_page)
	res.Header().Set("Content-Type", "text.html")
	res.Write(webpage)
}

func codeFunc(res http.ResponseWriter, req *http.Request) {
	doc := GetCodeInfoByDID(ctrct, UserDID)

	//JSON 인코딩
	body, _ := json.Marshal(doc)

	res.Header().Set("Content-Type", "application/json")
	res.Write(body)
}

func tokenFunc(res http.ResponseWriter, req *http.Request) {
	doc := GetTokenInfoByDID(ctrct, UserDID)

	//JSON 인코딩
	body, _ := json.Marshal(doc)

	res.Header().Set("Content-Type", "application/json")
	res.Write(body)
}

////////////////////////////////////////////////////////////////
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

	log.Println("[Connector Test]--> Evaluate Transaction: ReadCodeInfo, function returns all the current assets on the ledger")
	result, err := contract.EvaluateTransaction("ReadCodeInfo", "TestCID")
	if err != nil {
		log.Fatalf("Failed to evaluate transaction: %v", err)
	}
	log.Println("[Connector Test] : ", string(result))

	return contract
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
		log.Fatalf("Failed to evaluate transaction: %v", err)
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
