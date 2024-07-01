package request

type ReqMessageChat struct {
	Secret string `query:"secret" validate:"required"`
}
