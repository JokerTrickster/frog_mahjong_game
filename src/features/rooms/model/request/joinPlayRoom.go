package request

type ReqJoinPlay struct {
	Password string `json:"password" query:"password" validate:"required"`
}
