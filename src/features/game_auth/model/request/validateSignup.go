package request

type ReqGameValidateSignup struct {
	Code  string `json:"code" form:"code" query:"code"`
	Email string `json:"email" form:"email" query:"email"`
}
