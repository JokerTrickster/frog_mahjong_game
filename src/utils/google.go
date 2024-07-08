package utils

import (
	"fmt"
	"main/utils/aws"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var GoogleConfig = oauth2.Config{}

func InitGoogleOauth() error {
	clientID, err := getClientID()
	if err != nil {
		return err
	}
	clientSecret, err := getClientSecret()
	if err != nil {
		return err
	}
	GoogleConfig = oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  "http://localhost:8080/v0.1/auth/google/callback",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
	return nil
}

func getClientID() (string, error) {
	if Env.IsLocal {
		clientID, ok := os.LookupEnv("GOOGLE_CLIENT_ID")
		if !ok {
			return "", fmt.Errorf("GOOGLE_CLIENT_ID not found")
		}
		return clientID, nil
	} else {
		ClientID, err := aws.AwsSsmGetParam("google_client_id")
		if err != nil {
			return "", err
		}
		return ClientID, nil
	}
}

func getClientSecret() (string, error) {
	if Env.IsLocal {

		clientSecret, ok := os.LookupEnv("GOOGLE_CLIENT_SECRET")
		if !ok {
			return "", fmt.Errorf("GOOGLE_CLIENT_ID not found")
		}
		return clientSecret, nil

	} else {
		ClientID, err := aws.AwsSsmGetParam("google_client_secret")
		if err != nil {
			return "", err
		}
		return ClientID, nil
	}
}
