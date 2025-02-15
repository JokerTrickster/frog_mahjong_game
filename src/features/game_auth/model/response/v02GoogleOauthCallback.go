package response

type ResGameV02GoogleOauthCallback struct {
	AccessToken      string `json:"accessToken"`
	RefreshToken     string `json:"refreshToken"`
	UserID           uint   `json:"userID"`
	IsDuplicateLogin bool   `json:"isDuplicateLogin"`
}
