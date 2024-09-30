package request

type ReqRequestSignup struct {
	Email string `json:"email" validate:"required,email"`
}
