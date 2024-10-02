package response

type ResListProfile struct {
	Profiles []Profile `json:"profiles"`
}

type Profile struct {
	ProfileID  int    `json:"profileID"`
	Name       string `json:"name"`
	TotalCount int    `json:"totalCount"`
}
