package entity

type ValidatePasswordAuthEntity struct {
	Email    string `json:"email"`
	Code     string `json:"code"`
	Password string `json:"password"`
}
