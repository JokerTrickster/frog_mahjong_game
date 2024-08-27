package request

// 닉네임, 비밀번호, 이메일 정도만 정보를 받는다.
type ReqSignup struct {
	Name     string `json:"name" validate:"min=2,max=6"`
	Password string `json:"password" validate:"min=6,max=10"`
	Email    string `json:"email" validate:"email"`
}
