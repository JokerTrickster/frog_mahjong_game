package response

type ResV02GoogleOauthCallback struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	UserID       uint   `json:"userID"`
}
