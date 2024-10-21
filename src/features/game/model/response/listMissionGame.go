package response

type ResListMissionGame struct {
	Missions []Mission `json:"missions"`
}

type Mission struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}
