package request

type ReqList struct {
	Page     int `query:"page" validate:"omitempty,gte=0"`
	PageSize int `query:"pageSize" validate:"omitempty,gt=0"`
}
