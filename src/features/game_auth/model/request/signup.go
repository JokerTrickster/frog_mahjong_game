package request

// 닉네임, 비밀번호, 이메일 정도만 정보를 받는다.
type ReqGameSignup struct {
	Name     string `json:"name" validate:"min=2,max=5"`
	AuthCode string `json:"authCode"`
	Password string `json:"password" validate:"min=6"`
	Email    string `json:"email" validate:"email"`
}
