package response

type ResV2ListCardGame struct {
	Cards      []BirdCard `json:"cards"`
	TotalCount int        `json:"totalCount"`
}

type BirdCard struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Image         string `json:"image"`
	Size          int    `json:"size"`
	Habitat       string `json:"habitat"`
	BeakDirection string `json:"beakDirection"`
	Nest          string `json:"nest"`
}
