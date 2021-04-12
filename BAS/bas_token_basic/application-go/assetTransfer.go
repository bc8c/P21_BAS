/*
Copyright 2020 IBM All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

func main() {

	log.Println("============ application-golang starts ============")

	err := os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	if err != nil {
		log.Fatalf("Error setting DISCOVERY_AS_LOCALHOST environemnt variable: %v", err)
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

	log.Println("--> Submit Transaction: InitLedger, function creates the initial set of assets on the ledger")
	result, err := contract.SubmitTransaction("InitLedger")
	if err != nil {
		log.Fatalf("Failed to Submit transaction: %v", err)
	}
	log.Println(string(result))

	log.Println("--> Evaluate Transaction: ReadCodeInfo, function returns all the current assets on the ledger")
	result, err = contract.EvaluateTransaction("ReadCodeInfo", "TestCID")
	if err != nil {
		log.Fatalf("Failed to evaluate transaction: %v", err)
	}
	log.Println(string(result))

	log.Println("--> Evaluate Transaction: ReadTokenInfo, function returns all the current assets on the ledger")
	result, err = contract.EvaluateTransaction("ReadTokenInfo", "TestTID")
	if err != nil {
		log.Fatalf("Failed to evaluate transaction: %v", err)
	}
	log.Println(string(result))

	log.Println("--> Submit Transaction: CreateCodeInfo, creates new CodeInfo")
	// mHashV := [32]byte{0}
	// result, err = contract.SubmitTransaction("CreateCodeInfo", "TestCID5", "TestDIDRO2", "TESTDIDClient2", "all", string(mHashV[:]), "TestTime2", "TESTURI2", "Available", "TestTID2")
	result, err = contract.SubmitTransaction("CreateCodeInfo", "TestCID8", "TestDIDRO2", "TESTDIDClient2", "all", "HashTest", "TestTime2", "TESTURI2", "Available", "TestTID2")
	if err != nil {
		log.Fatalf("Failed to Submit transaction: %v", err)
	}
	log.Println(string(result))

	log.Println("--> Evaluate Transaction: ReadCodeInfo, function returns all the current assets on the ledger")
	result, err = contract.EvaluateTransaction("ReadCodeInfo", "TestCID8")
	if err != nil {
		log.Fatalf("Failed to evaluate transaction: %v", err)
	}
	log.Println(string(result))

	log.Println("--> Submit Transaction: CreateTokenInfo, creates new CodeInfo")
	// mHashV = [32]byte{0}
	// result, err = contract.SubmitTransaction("CreateTokenInfo", "TestTID5", "TestDIDRO2", "TESTDIDClient2", "all", string(mHashV[:]), "TestTime2", "TestDuration", "TESTURI2", "Available")
	result, err = contract.SubmitTransaction("CreateTokenInfo", "TestTID8", "TestDIDRO2", "TESTDIDClient2", "all", "HashTest", "TestTime2", "TESTURI2", "Available", "TestTID2")
	if err != nil {
		log.Fatalf("Failed to Submit transaction: %v", err)
	}
	log.Println(string(result))

	log.Println("--> Evaluate Transaction: ReadTokenInfo, function returns all the current assets on the ledger")
	result, err = contract.EvaluateTransaction("ReadTokenInfo", "TestTID8")
	if err != nil {
		log.Fatalf("Failed to evaluate transaction: %v", err)
	}
	log.Println(string(result))

	log.Println("--> Evaluate Transaction: GetCodeInfoByDID, function returns all the current assets on the ledger")
	result, err = contract.EvaluateTransaction("GetCodeInfoByDID", "TestDIDRO2")
	if err != nil {
		log.Fatalf("Failed to evaluate transaction: %v", err)
	}
	log.Println(string(result))

	log.Println("--> Evaluate Transaction: GetTokenInfoByDID, function returns all the current assets on the ledger")
	result, err = contract.EvaluateTransaction("GetTokenInfoByDID", "TestDIDRO2")
	if err != nil {
		log.Fatalf("Failed to evaluate transaction: %v", err)
	}
	log.Println(string(result))

	log.Println("============ application-golang ends ============")
}

func populateWallet(wallet *gateway.Wallet) error {
	log.Println("============ Populating wallet ============")
	credPath := filepath.Join(
		"..",
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
