package request

type ReqFindItPasswordCheckBoardGame struct {
	Password string `json:"password" validate:"required"`
}
