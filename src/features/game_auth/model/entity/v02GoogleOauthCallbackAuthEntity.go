package entity


type V02GoogleOauthCallbackSQLQuery struct {
	Email string `json:"email"`
}

// entity
type V02GoogleUser struct {
    ID            string `json:"id"`
    Email         string `json:"email"`
    VerifiedEmail bool   `json:"verified_email"`
    Picture       string `json:"picture"`
    HD            string `json:"hd"`
}