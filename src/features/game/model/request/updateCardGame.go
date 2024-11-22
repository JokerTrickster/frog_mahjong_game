package request

type ReqUpdateCard struct {
	Card UpdateCard `json:"card"`
}

type UpdateCard struct {
	CardID        int    `json:"cardID"`
	Name          string `json:"name"`
	Size          int    `json:"size"`
	Habitat       string `json:"habitat"`
	Nest          string `json:"nest"`
	BeakDirection string `json:"beakDirection"`
}
