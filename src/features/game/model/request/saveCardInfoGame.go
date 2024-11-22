package request

type ReqSaveCardInfo struct {
	Cards []CardInfo `json:"cards"`
}

type CardInfo struct {
	Name          string `json:"name"`
	Size          int    `json:"size"`
	Habitat       string `json:"habitat"`
	Nest          string `json:"nest"`
	BeakDirection string `json:"beakDirection"`
	Image         string `json:"image"`
}
