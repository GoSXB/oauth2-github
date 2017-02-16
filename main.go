package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var (
	UserToken   *oauth2.Token
	OAuthConfig = oauth2.Config{
		RedirectURL: "http://localhost:5000",
		Endpoint:    github.Endpoint,
		Scopes:      []string{"user", "user:email"},
	}
)

func main() {
	OAuthConfig.ClientID = os.Args[1]
	OAuthConfig.ClientSecret = os.Args[2]

	http.HandleFunc("/profile", ProfileHandler)
	http.HandleFunc("/", OAuthCallbackHandler)

	go http.ListenAndServe(":5000", nil)
	log.Println("Listening on :5000\n")

	fmt.Printf("Go on the following URL and login:\n\n%v\n", OAuthConfig.AuthCodeURL(""))

	// Wait for SIGINT to quit (CTRL +C)
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGINT)
	<-signalChan

	fmt.Println("Got SIGINT, stopping app")
}

func OAuthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	tempOAuthCode := r.URL.Query().Get("code")
	if tempOAuthCode != "" {
		// If there is a callback code query parameter
		// https://developer.github.com/v3/oauth/#2-github-redirects-back-to-your-site

		// Exchange it for a permanent token
		handleTempToken(tempOAuthCode, w, r)
		return
	}
	w.WriteHeader(200)
	fmt.Fprintf(w, "nothing to do here")
}

func handleTempToken(code string, w http.ResponseWriter, r *http.Request) {
	token, err := OAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		w.WriteHeader(500)
		log.Println("Fail to exchange temp token")
		return
	}
	UserToken = token
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(200)
	fmt.Fprintf(w, "app is connected, go on <a href=\"/profile\">Profile Page</a>")
	return
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	// If app is not authenticated -> 401
	if UserToken == nil {
		w.WriteHeader(401)
		fmt.Fprintf(w, "unauthorized")
		return
	}

	// Otherwise display authenticated user profile in JSON
	client := OAuthConfig.Client(context.Background(), UserToken)
	res, err := client.Get("https://api.github.com/user")
	if err != nil {
		w.WriteHeader(500)
		log.Println("fail to make github request:", err)
		return
	}
	defer res.Body.Close()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	io.Copy(w, res.Body)
}
