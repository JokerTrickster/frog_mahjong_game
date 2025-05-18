package response

type ResListCardGame struct {
	TotalCount int        `json:"totalCount"`
	Cards      []FrogCard `json:"cards"`
}

type FrogCard struct {
	ID    int    `json:"id"`
	Count int    `json:"count"`
	Color string `json:"color"`
}
