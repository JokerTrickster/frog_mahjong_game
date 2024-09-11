package request

type ReqWSPlayTogetherEvent struct {
	Timer int `json:"timer" validate:"required"`
	Count int `json:"count" validate:"required"`
}
