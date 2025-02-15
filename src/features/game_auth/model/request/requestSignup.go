package request

type ReqGameRequestSignup struct {
	Email string `json:"email" validate:"required,email"`
}
