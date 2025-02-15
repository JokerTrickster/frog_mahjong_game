package request

type ReqGameValidatePassword struct {
	Code     string `json:"code" form:"code" query:"code"`
	Password string `json:"password" form:"password" query:"password"`
	Email    string `json:"email" form:"email" query:"email"`
}
