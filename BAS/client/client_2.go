package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

const (
	authServerURL = "http://localhost:9096"
)

var (
	config = oauth2.Config{
		ClientID:     "did:BAS:987654321abcdefghi",
		ClientSecret: "22222222",
		Scopes:       []string{"all"},
		RedirectURL:  "http://localhost:9094/oauth2",
		Endpoint: oauth2.Endpoint{
			AuthURL:  authServerURL + "/oauth/authorize",
			TokenURL: authServerURL + "/oauth/token",
		},
	}
	globalToken *oauth2.Token // Non-concurrent security
)

// DY mod START
// To measure execution time
// var mStartTime Time
// var mElapsedTime Time

// DY mod END
var mStartTime = time.Now()
var mElapsedTime = time.Since(mStartTime)
var numTokenCreation int64 = 0
var numResourceAccess int64 = 0

func main() {

	// DY mod START
	if len(os.Args) < 3 {
		log.Printf("Argument Error\n")
		return
	}
	numTokenCreation, _ := strconv.Atoi(os.Args[1])
	numResourceAccess, _ := strconv.Atoi(os.Args[2])
	log.Printf("numTokenCreation : %d \n", numTokenCreation)
	log.Printf("numResourceAccess :%d \n", numResourceAccess)
	// DY mod END

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// DY mod START
		if numTokenCreation == 0 {
			return
		}
		mStartTime = time.Now()
		// DY mod END
		u := config.AuthCodeURL("xyz",
			oauth2.SetAuthURLParam("code_challenge", genCodeChallengeS256("s256example")),
			oauth2.SetAuthURLParam("code_challenge_method", "S256"))
		http.Redirect(w, r, u, http.StatusFound)
	})

	http.HandleFunc("/oauth2", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		state := r.Form.Get("state")
		if state != "xyz" {
			http.Error(w, "State invalid", http.StatusBadRequest)
			return
		}
		code := r.Form.Get("code")
		if code == "" {
			http.Error(w, "Code not found", http.StatusBadRequest)
			return
		}
		// DY mod START
		mElapsedTime = time.Since(mStartTime)
		if mElapsedTime < time.Millisecond*15 {
			log.Printf("CTtime : %d: %s", numTokenCreation, mElapsedTime)
			numTokenCreation -= 1
		}

		// DY mod END

		time.Sleep(time.Second * 1)

		mStartTime = time.Now()

		token, err := config.Exchange(context.Background(), code, oauth2.SetAuthURLParam("code_verifier", "s256example"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		globalToken = token

		// DY mod START
		mElapsedTime = time.Since(mStartTime)
		if mElapsedTime < time.Millisecond*20 {
			// log.Printf("Ttime : %d: %s", numTokenCreation, mElapsedTime)
			// numTokenCreation -= 1
		}

		// DY mod END

		e := json.NewEncoder(w)
		e.SetIndent("", "  ")
		e.Encode(token)

		// DY mod START
		// time.Sleep(time.Millisecond)

		http.Get("http://localhost:9094/")
		// DY mod END
	})

	http.HandleFunc("/refresh", func(w http.ResponseWriter, r *http.Request) {
		if globalToken == nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		globalToken.Expiry = time.Now()
		token, err := config.TokenSource(context.Background(), globalToken).Token()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		globalToken = token
		e := json.NewEncoder(w)
		e.SetIndent("", "  ")
		e.Encode(token)
	})

	http.HandleFunc("/try", func(w http.ResponseWriter, r *http.Request) {
		// DY mod START
		if numResourceAccess == 0 {
			return
		}
		mStartTime = time.Now()
		// DY mod END

		if globalToken == nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		resp, err := http.Get(fmt.Sprintf("%s/test?access_token=%s", authServerURL, globalToken.AccessToken))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer resp.Body.Close()

		io.Copy(w, resp.Body)

		// DY mod START
		mElapsedTime = time.Since(mStartTime)
		if mElapsedTime < time.Millisecond*15 {
			log.Printf("Rtime : %d: %s", numResourceAccess, mElapsedTime)
			numResourceAccess -= 1
		}
		time.Sleep(time.Millisecond)
		http.Get("http://localhost:9094/try")
		// DY mod END

	})

	http.HandleFunc("/pwd", func(w http.ResponseWriter, r *http.Request) {
		token, err := config.PasswordCredentialsToken(context.Background(), "test", "test")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		globalToken = token
		e := json.NewEncoder(w)
		e.SetIndent("", "  ")
		e.Encode(token)
	})

	http.HandleFunc("/client", func(w http.ResponseWriter, r *http.Request) {
		cfg := clientcredentials.Config{
			ClientID:     config.ClientID,
			ClientSecret: config.ClientSecret,
			TokenURL:     config.Endpoint.TokenURL,
		}

		token, err := cfg.Token(context.Background())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		e := json.NewEncoder(w)
		e.SetIndent("", "  ")
		e.Encode(token)
	})

	log.Println("Client is running at 9094 port.Please open http://localhost:9094")
	log.Fatal(http.ListenAndServe(":9094", nil))
}

func genCodeChallengeS256(s string) string {
	s256 := sha256.Sum256([]byte(s))
	return base64.URLEncoding.EncodeToString(s256[:])
}
