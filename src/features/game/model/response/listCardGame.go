package response

type ResListCardGame struct {
	Cards []BirdCard `json:"cards"`
}

type BirdCard struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Image         string `json:"image"`
	Size          int    `json:"size"`
	Habitat       string `json:"habitat"`
	BeakDirection string `json:"beakDirection"`
}
