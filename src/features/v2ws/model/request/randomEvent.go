package request

type ReqWSRandom struct {
	Count int `query:"count" validate:"required"`
}
