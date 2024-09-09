package request

type ReqCreate struct {
	Name     string `json:"name" validate:"required"`
	Password string `json:"password,omitempty"`
	MaxCount int    `json:"maxCount" validate:"required"`
	MinCount int    `json:"minCount" validate:"required"`
	Timer    int    `json:"timer" validate:"required"`
}
