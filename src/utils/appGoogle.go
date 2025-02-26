package utils

import (
	"fmt"
	"main/utils/aws"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var AppGoogleConfig = oauth2.Config{}

func InitAppGoogleOauth() error {
	clientID, err := getAppClientID()
	if err != nil {
		return err
	}
	clientSecret, err := getAppClientSecret()
	if err != nil {
		return err
	}
	AppGoogleConfig = oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
	if Env.IsLocal {
		AppGoogleConfig.RedirectURL = "http://localhost:8080/v0.1/game/auth/google/callback"
	} else {
		// GoogleConfig.RedirectURL = fmt.Sprintf("https://%s-frog-api.jokertrickster.com/v0.1/auth/google/callback", Env.Env)
		AppGoogleConfig.RedirectURL = "https://dev-frog-api.jokertrickster.com/v0.1/game/auth/google/callback"
		///v0.1/game/auth/google/callback

	}
	return nil
}

func getAppClientID() (string, error) {
	if Env.IsLocal {
		clientID, ok := os.LookupEnv("APP_GOOGLE_CLIENT_ID")
		if !ok {
			return "", fmt.Errorf("APP_GOOGLE_CLIENT_ID not found")
		}
		return clientID, nil
	} else {
		ClientID, err := aws.AwsSsmGetParam("app_google_client_id")
		if err != nil {
			return "", err
		}
		return ClientID, nil
	}
}

func getAppClientSecret() (string, error) {
	if Env.IsLocal {

		clientSecret, ok := os.LookupEnv("APP_GOOGLE_CLIENT_SECRET")
		if !ok {
			return "", fmt.Errorf("APP_GOOGLE_CLIENT_ID not found")
		}
		return clientSecret, nil

	} else {
		ClientID, err := aws.AwsSsmGetParam("app_google_client_secret")
		if err != nil {
			return "", err
		}
		return ClientID, nil
	}
}
