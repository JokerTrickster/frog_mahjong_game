package response

type ResListCardGame struct {
	TotalCount int        `json:"totalCount"`
	Cards      []FrogCard `json:"cards"`
}

type FrogCard struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}
