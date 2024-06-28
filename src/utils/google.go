package utils

import (
	"main/utils/aws"

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
		return Env.GoogleClientID, nil

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
		return Env.GoogleClientSecret, nil

	} else {
		ClientID, err := aws.AwsSsmGetParam("google_client_secret")
		if err != nil {
			return "", err
		}
		return ClientID, nil
	}
}
