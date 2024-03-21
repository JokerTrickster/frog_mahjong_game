package request

type ReqCreate struct {
	Name     string `json:"name" validate:"required"`
	Password string `json:"password,omitempty"`
	MaxCount int    `json:"max_count" validate:"required"`
	MinCount int    `json:"min_count" validate:"required"`
}
