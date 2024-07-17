package request

type ReqRequestPassword struct {
	Email string `json:"email" validate:"required,email"`
}
