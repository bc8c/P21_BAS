package main

import (
	"log"
	"time"
	"crypto/sha256"
	"encoding/base64"
	// "fmt"
	"bytes"
)

func main(){

	
	var mTime = time.Now().String()

	log.Printf("%s\n", mTime)

	var mHash = genHashS256(mTime)

	// log.Printf("%s\n", mID)

	// a := "CI_"+mID

	// log.Printf("%s\n", fmt.Sprintf("%s%s", "CI_", mID))

	var mID bytes.Buffer
	mID.WriteString("CI_")
	mID.WriteString(mHash)

	log.Printf("%s\n", mID.String())


	
}

// generate Hash value
func genHashS256(s string) string {
	s256 := sha256.Sum256([]byte(s))
	return base64.StdEncoding.EncodeToString(s256[:])
	//base64.StdEncoding.DecodeString(s256[:])
}

