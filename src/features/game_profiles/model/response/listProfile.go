package response

type ResListGameProfile struct {
	Profiles []Profile `json:"profiles"`
}

type Profile struct {
	ProfileID  int    `json:"profileID"`
	Name       string `json:"name"`
	Image      string `json:"image"`
}
