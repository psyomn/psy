package book

import (
	"fmt"
	"log"

	"github.com/psyomn/psy/oauth1"
)

const baseURL = "https://www.goodreads.com"

// const baseURL = "http://127.0.0.1:9090"

var grEndpoint = oauth1.Endpoint{
	RequestTokenURL: baseURL + "/oauth/request_token",
	AuthorizeURL:    baseURL + "/oauth/authorize",
	AccessTokenURL:  baseURL + "/oauth/access_token",
}

var config oauth1.Config

func authenticate() error {
	config = oauth1.Config{
		CallbackURL:    "oob",
		ConsumerKey:    grKey,
		ConsumerSecret: grSecret,
		Endpoint:       grEndpoint,
	}

	fmt.Printf(".%s.\n", grKey)
	fmt.Printf(".%s.\n", grSecret)

	requestToken, requestSecret, err := login()
	if err != nil {
		log.Fatalf("Request Token Phase: %s", err.Error())
	}
	accessToken, err := receivePIN(requestToken, requestSecret)
	if err != nil {
		log.Fatalf("Access Token Phase: %s", err.Error())
	}

	fmt.Println("Consumer was granted an access token to act on behalf of a user.")
	fmt.Printf("token: %s\nsecret: %s\n", accessToken.Token, accessToken.TokenSecret)

	return nil
}

func receivePIN(requestToken, requestSecret string) (*oauth1.Token, error) {
	fmt.Printf("Choose whether to grant the application access.\nPaste " +
		"the oauth_verifier parameter (excluding trailing #_=_) from the " +
		"address bar: ")
	var verifier string
	_, err := fmt.Scanf("%s", &verifier)
	accessToken, accessSecret, err := config.AccessToken(requestToken, requestSecret, verifier)
	if err != nil {
		return nil, err
	}
	return oauth1.NewToken(accessToken, accessSecret), err
}

func login() (requestToken, requestSecret string, err error) {
	fmt.Println(config)
	requestToken, requestSecret, err = config.RequestToken()
	if err != nil {
		return "", "", err
	}
	authorizationURL, err := config.AuthorizationURL(requestToken)
	if err != nil {
		return "", "", err
	}
	fmt.Printf("Open this URL in your browser:\n%s\n", authorizationURL.String())
	return requestToken, requestSecret, err
}
