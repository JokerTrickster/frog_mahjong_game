package request

type ReqV02GoogleOauthCallback struct {
	Code string `query:"code"`
}
