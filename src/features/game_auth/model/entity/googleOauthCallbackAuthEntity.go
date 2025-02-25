package entity

type GoogleOauthCallbackSQLQuery struct {
	Email string `json:"email"`
}

// entity
type GoogleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
	HD            string `json:"hd"`
}
