package request

type ReqGameRequestPassword struct {
	Email string `json:"email" validate:"required,email"`
}
