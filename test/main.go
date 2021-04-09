package main

import (
	"log"
	"os"

	// "github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	

	"ttest/bcconnector"
	"ttest/oauth2/errors"
)

func main() {
	log.Println("============ application-golang starts ============")
	bcconnector.NewConnector()
	
	err := os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	if err != nil {
		log.Fatalf("Error setting DISCOVERY_AS_LOCALHOST environemnt variable: %v", err)
	}

	// log.Println(oauth2.Code)
	log.Println(errors.ErrExpiredAccessToken)

	// wallet, err := gateway.NewFileSystemWallet("wallet")
	// if err != nil {
	// 	log.Fatalf("Failed to create wallet: %v", err)
	// }

	// if !wallet.Exists("appUser") {
	// 	log.Println("!!!!!!!!!!!")
	// }

	

}


